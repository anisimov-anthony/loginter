package analyzer_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/anisimov-anthony/loginter/pkg/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func testdataDir(t *testing.T) string {
	t.Helper()
	_, testFilename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to get current test filename")
	}
	return filepath.Join(filepath.Dir(testFilename), "testdata")
}

// cfgLowercaseOnly returns a config with only the lowercase check enabled.
func cfgLowercaseOnly() analyzer.Config {
	return analyzer.Config{CheckLowercase: true}
}

// cfgEnglishOnly returns a config with only the English check enabled.
func cfgEnglishOnly() analyzer.Config {
	return analyzer.Config{CheckEnglish: true}
}

// cfgSpecialOnly returns a config with only the special chars check enabled.
func cfgSpecialOnly() analyzer.Config {
	return analyzer.Config{CheckSpecial: true}
}

// cfgSensitiveOnly returns a config with only the sensitive data check enabled.
func cfgSensitiveOnly() analyzer.Config {
	return analyzer.Config{CheckSensitive: true}
}

// TestAnalyzerBasic runs all checks on the basic slog testdata package.
func TestAnalyzerBasic(t *testing.T) {
	a := analyzer.NewAnalyzer(analyzer.DefaultConfig())
	analysistest.Run(t, testdataDir(t), a, "basic")
}

// TestAnalyzerZap runs all checks on the zap testdata package.
func TestAnalyzerZap(t *testing.T) {
	a := analyzer.NewAnalyzer(analyzer.DefaultConfig())
	analysistest.Run(t, testdataDir(t), a, "zap")
}

// TestAnalyzerSensitive runs all checks on the sensitive testdata package.
func TestAnalyzerSensitive(t *testing.T) {
	a := analyzer.NewAnalyzer(analyzer.DefaultConfig())
	analysistest.Run(t, testdataDir(t), a, "sensitive")
}

// TestAnalyzerClean verifies no diagnostics are reported for clean code.
func TestAnalyzerClean(t *testing.T) {
	a := analyzer.NewAnalyzer(analyzer.DefaultConfig())
	analysistest.Run(t, testdataDir(t), a, "clean")
}

// TestAnalyzerSlogLogger tests slog.Logger instance method calls.
func TestAnalyzerSlogLogger(t *testing.T) {
	a := analyzer.NewAnalyzer(analyzer.DefaultConfig())
	analysistest.Run(t, testdataDir(t), a, "sloglogger")
}

// TestAnalyzerDisabledChecks verifies that disabled rules produce no diagnostics.
func TestAnalyzerDisabledChecks(t *testing.T) {
	cfg := analyzer.Config{
		CheckLowercase: true,
		CheckEnglish:   false,
		CheckSpecial:   false,
		CheckSensitive: false,
	}
	a := analyzer.NewAnalyzer(cfg)
	analysistest.Run(t, testdataDir(t), a, "config")
}

// TestAnalyzerCustomPatterns verifies that user-supplied sensitive patterns work.
func TestAnalyzerCustomPatterns(t *testing.T) {
	cfg := analyzer.Config{
		CheckLowercase:    false,
		CheckEnglish:      false,
		CheckSpecial:      false,
		CheckSensitive:    true,
		SensitivePatterns: []string{"ssn", "credit_card"},
	}
	a := analyzer.NewAnalyzer(cfg)
	analysistest.Run(t, testdataDir(t), a, "custompatterns")
}

// TestSuggestedFixesBasic verifies auto-fixes on the slog basic testdata.
func TestSuggestedFixesBasic(t *testing.T) {
	a := analyzer.NewAnalyzer(analyzer.DefaultConfig())
	analysistest.RunWithSuggestedFixes(t, testdataDir(t), a, "basic")
}

// TestSuggestedFixesSlogLowercase verifies auto-fixes for the lowercase rule on slog.
func TestSuggestedFixesSlogLowercase(t *testing.T) {
	a := analyzer.NewAnalyzer(cfgLowercaseOnly())
	analysistest.RunWithSuggestedFixes(t, testdataDir(t), a, "slog_lowercase")
}

// TestSuggestedFixesSlogSpecial verifies auto-fixes for the special-chars rule on slog.
func TestSuggestedFixesSlogSpecial(t *testing.T) {
	a := analyzer.NewAnalyzer(cfgSpecialOnly())
	analysistest.RunWithSuggestedFixes(t, testdataDir(t), a, "slog_special")
}

// TestSuggestedFixesZapLowercase verifies auto-fixes for the lowercase rule on zap.
func TestSuggestedFixesZapLowercase(t *testing.T) {
	a := analyzer.NewAnalyzer(cfgLowercaseOnly())
	analysistest.RunWithSuggestedFixes(t, testdataDir(t), a, "zap_lowercase")
}

// TestSuggestedFixesZapSpecial verifies auto-fixes for the special-chars rule on zap.
func TestSuggestedFixesZapSpecial(t *testing.T) {
	a := analyzer.NewAnalyzer(cfgSpecialOnly())
	analysistest.RunWithSuggestedFixes(t, testdataDir(t), a, "zap_special")
}

// TestSlogLowercase runs only the lowercase check on dedicated slog testdata.
func TestSlogLowercase(t *testing.T) {
	a := analyzer.NewAnalyzer(cfgLowercaseOnly())
	analysistest.Run(t, testdataDir(t), a, "slog_lowercase")
}

// TestSlogEnglish runs only the English check on dedicated slog testdata.
func TestSlogEnglish(t *testing.T) {
	a := analyzer.NewAnalyzer(cfgEnglishOnly())
	analysistest.Run(t, testdataDir(t), a, "slog_english")
}

// TestSlogSpecialChars runs only the special-chars check on dedicated slog testdata.
func TestSlogSpecialChars(t *testing.T) {
	a := analyzer.NewAnalyzer(cfgSpecialOnly())
	analysistest.Run(t, testdataDir(t), a, "slog_special")
}

// TestSlogSensitive runs only the sensitive-data check on dedicated slog testdata.
func TestSlogSensitive(t *testing.T) {
	a := analyzer.NewAnalyzer(cfgSensitiveOnly())
	analysistest.Run(t, testdataDir(t), a, "slog_sensitive")
}

// TestZapLowercase runs only the lowercase check on dedicated zap testdata.
func TestZapLowercase(t *testing.T) {
	a := analyzer.NewAnalyzer(cfgLowercaseOnly())
	analysistest.Run(t, testdataDir(t), a, "zap_lowercase")
}

// TestZapEnglish runs only the English check on dedicated zap testdata.
func TestZapEnglish(t *testing.T) {
	a := analyzer.NewAnalyzer(cfgEnglishOnly())
	analysistest.Run(t, testdataDir(t), a, "zap_english")
}

// TestZapSpecialChars runs only the special-chars check on dedicated zap testdata.
func TestZapSpecialChars(t *testing.T) {
	a := analyzer.NewAnalyzer(cfgSpecialOnly())
	analysistest.Run(t, testdataDir(t), a, "zap_special")
}

// TestZapSensitive runs only the sensitive-data check on dedicated zap testdata.
func TestZapSensitive(t *testing.T) {
	a := analyzer.NewAnalyzer(cfgSensitiveOnly())
	analysistest.Run(t, testdataDir(t), a, "zap_sensitive")
}

// TestAnalyzerAllDisabledNoReports verifies that disabling all checks produces no output.
func TestAnalyzerAllDisabledNoReports(t *testing.T) {
	cfg := analyzer.Config{
		CheckLowercase: false,
		CheckEnglish:   false,
		CheckSpecial:   false,
		CheckSensitive: false,
	}
	a := analyzer.NewAnalyzer(cfg)
	analysistest.Run(t, testdataDir(t), a, "clean")
}
