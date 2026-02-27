package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// NewAnalyzer creates a new loginter analyzer with the given configuration.
func NewAnalyzer(cfg Config) *analysis.Analyzer {
	a := &loginterAnalyzer{cfg: cfg}
	return &analysis.Analyzer{
		Name:     "loginter",
		Doc:      "checks log messages for common issues\n\nloginter verifies that log messages start with a lowercase letter, are in English, do not contain special characters or emoji, and do not contain sensitive data.",
		Run:      a.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

type loginterAnalyzer struct {
	cfg Config
}

// slog package-level function names that accept a message as the first argument.
var slogFunctions = map[string]bool{
	"Info":         true,
	"Warn":         true,
	"Error":        true,
	"Debug":        true,
	"InfoContext":  true,
	"WarnContext":  true,
	"ErrorContext": true,
	"DebugContext": true,
	"Log":          true,
}

// zap Logger method names that accept a message as the first argument.
var zapLoggerMethods = map[string]bool{
	"Info":   true,
	"Warn":   true,
	"Error":  true,
	"Debug":  true,
	"Fatal":  true,
	"Panic":  true,
	"DPanic": true,
	"Log":    true,
}

// zap SugaredLogger method names that accept a message as the first argument.
var zapSugarMethods = map[string]bool{
	"Infof":   true,
	"Warnf":   true,
	"Errorf":  true,
	"Debugf":  true,
	"Fatalf":  true,
	"Panicf":  true,
	"DPanicf": true,
	"Infow":   true,
	"Warnw":   true,
	"Errorw":  true,
	"Debugw":  true,
	"Fatalw":  true,
	"Panicw":  true,
	"DPanicw": true,
	"Infoln":  true,
	"Warnln":  true,
	"Errorln": true,
	"Debugln": true,
	"Fatalln": true,
	"Panicln": true,
}

func (a *loginterAnalyzer) run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		a.checkCall(pass, call)
	})

	return nil, nil
}

func (a *loginterAnalyzer) checkCall(pass *analysis.Pass, call *ast.CallExpr) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	methodName := sel.Sel.Name

	if a.isSlogCall(pass, sel, methodName) || a.isZapCall(pass, sel, methodName) {
		a.checkLogArgs(pass, call, sel, methodName)
	}
}

// isSlogCall checks if the call is to the log/slog package.
func (a *loginterAnalyzer) isSlogCall(pass *analysis.Pass, sel *ast.SelectorExpr, methodName string) bool {
	if ident, ok := sel.X.(*ast.Ident); ok {
		if obj := pass.TypesInfo.Uses[ident]; obj != nil {
			if pkgName, ok := obj.(*types.PkgName); ok {
				if pkgName.Imported().Path() == "log/slog" && slogFunctions[methodName] {
					return true
				}
			}
		}
	}

	typ := pass.TypesInfo.TypeOf(sel.X)
	if typ == nil {
		return false
	}
	typ = derefPointer(typ)
	if named, ok := typ.(*types.Named); ok {
		if named.Obj() != nil && named.Obj().Pkg() != nil {
			if named.Obj().Pkg().Path() == "log/slog" && named.Obj().Name() == "Logger" && slogFunctions[methodName] {
				return true
			}
		}
	}

	return false
}

// isZapCall checks if the call is to the go.uber.org/zap package.
func (a *loginterAnalyzer) isZapCall(pass *analysis.Pass, sel *ast.SelectorExpr, methodName string) bool {
	typ := pass.TypesInfo.TypeOf(sel.X)
	if typ == nil {
		return false
	}
	typ = derefPointer(typ)
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	if named.Obj() == nil || named.Obj().Pkg() == nil {
		return false
	}
	pkgPath := named.Obj().Pkg().Path()
	if pkgPath != "go.uber.org/zap" {
		return false
	}
	typeName := named.Obj().Name()

	if typeName == "Logger" && zapLoggerMethods[methodName] {
		return true
	}
	if typeName == "SugaredLogger" && zapSugarMethods[methodName] {
		return true
	}

	return false
}

// checkLogArgs extracts the message argument from a log call and runs checks.
func (a *loginterAnalyzer) checkLogArgs(pass *analysis.Pass, call *ast.CallExpr, sel *ast.SelectorExpr, methodName string) {
	if len(call.Args) == 0 {
		return
	}

	msgArgIndex := 0
	if methodName == "Log" && a.isSlogCall(pass, sel, methodName) {
		msgArgIndex = 2
	} else if strings.HasSuffix(methodName, "Context") {
		msgArgIndex = 1
	}

	if msgArgIndex >= len(call.Args) {
		return
	}

	msgArg := call.Args[msgArgIndex]

	if a.cfg.CheckSensitive {
		a.checkSensitiveConcat(pass, msgArg)
	}

	msg, pos, end, ok := extractStringLiteral(pass, msgArg)
	if !ok {
		return
	}

	if a.cfg.CheckLowercase {
		if d := checkLowercase(msg, pos, end); d != nil {
			pass.Report(*d)
		}
	}

	if a.cfg.CheckEnglish {
		if d := checkEnglish(msg, pos, end); d != nil {
			pass.Report(*d)
		}
	}

	if a.cfg.CheckSpecial {
		if d := checkSpecialChars(msg, pos, end); d != nil {
			pass.Report(*d)
		}
	}

	if a.cfg.CheckSensitive {
		if d := checkSensitiveData(msg, pos, end, a.cfg.AllSensitivePatterns()); d != nil {
			pass.Report(*d)
		}
	}
}

// checkSensitiveConcat checks for string concatenation that may contain sensitive data.
// For example: log.Info("user password: " + password)
func (a *loginterAnalyzer) checkSensitiveConcat(pass *analysis.Pass, expr ast.Expr) {
	binExpr, ok := expr.(*ast.BinaryExpr)
	if !ok {
		return
	}
	if binExpr.Op != token.ADD {
		return
	}

	checkSide := func(side ast.Expr) {
		lit, ok := side.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return
		}
		val := strings.Trim(lit.Value, "`\"")
		lower := strings.ToLower(val)
		for _, pattern := range a.cfg.AllSensitivePatterns() {
			if strings.Contains(lower, strings.ToLower(pattern)) {
				pass.Report(analysis.Diagnostic{
					Pos:     binExpr.Pos(),
					End:     binExpr.End(),
					Message: "log message may contain sensitive data (\"" + pattern + "\")",
				})
				return
			}
		}
	}

	checkSide(binExpr.X)
	checkSide(binExpr.Y)

	if leftBin, ok := binExpr.X.(*ast.BinaryExpr); ok && leftBin.Op == token.ADD {
		a.checkSensitiveConcat(pass, binExpr.X)
	}
}

// extractStringLiteral extracts the string value from a string literal expression.
// Returns the unquoted string, the position of the literal (including quotes),
// and whether extraction was successful.
func extractStringLiteral(pass *analysis.Pass, expr ast.Expr) (string, token.Pos, token.Pos, bool) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", 0, 0, false
	}

	val := lit.Value
	if len(val) < 2 {
		return "", 0, 0, false
	}

	var unquoted string
	if val[0] == '`' {
		unquoted = val[1 : len(val)-1]
	} else {
		unquoted = unquoteString(val)
	}

	return unquoted, lit.Pos(), lit.End(), true
}

// unquoteString removes surrounding double quotes and handles basic escape sequences.
func unquoteString(s string) string {
	if len(s) < 2 {
		return s
	}
	s = s[1 : len(s)-1]

	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case 'r':
				b.WriteByte('\r')
			case '\\':
				b.WriteByte('\\')
			case '"':
				b.WriteByte('"')
			default:
				b.WriteByte(s[i])
				b.WriteByte(s[i+1])
			}
			i++
		} else {
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

// derefPointer returns the underlying type if t is a pointer type.
func derefPointer(t types.Type) types.Type {
	if ptr, ok := t.(*types.Pointer); ok {
		return ptr.Elem()
	}
	return t
}
