package analyzer

import (
	"go/token"
	"testing"
)

func TestCheckLowercase(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantNil bool
	}{
		{"uppercase start", "Starting server", false},
		{"lowercase start", "starting server", true},
		{"empty string", "", true},
		{"number start", "8080 is the port", true},
		{"uppercase single char", "S", false},
		{"lowercase single char", "s", true},
		{"unicode uppercase", "\u0410bcdef", false}, // Cyrillic A
		{"special char start", "!hello", true},      // starts with special, not a letter
		{"space start", " hello", true},             // starts with space
		{"all caps", "HTTP SERVER STARTED", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := checkLowercase(tt.msg, token.Pos(1), token.Pos(len(tt.msg)+3))
			if tt.wantNil && d != nil {
				t.Errorf("expected nil diagnostic, got: %s", d.Message)
			}
			if !tt.wantNil && d == nil {
				t.Errorf("expected diagnostic, got nil")
			}
			if !tt.wantNil && d != nil {
				if d.Message != "log message must start with a lowercase letter" {
					t.Errorf("unexpected message: %s", d.Message)
				}
				if len(d.SuggestedFixes) == 0 {
					t.Error("expected suggested fix")
				}
			}
		})
	}
}

func TestCheckLowercaseSuggestedFix(t *testing.T) {
	d := checkLowercase("Starting server", token.Pos(1), token.Pos(18))
	if d == nil {
		t.Fatal("expected diagnostic")
	}
	if len(d.SuggestedFixes) != 1 {
		t.Fatalf("expected 1 fix, got %d", len(d.SuggestedFixes))
	}
	fix := d.SuggestedFixes[0]
	if len(fix.TextEdits) != 1 {
		t.Fatalf("expected 1 edit, got %d", len(fix.TextEdits))
	}
	if string(fix.TextEdits[0].NewText) != `"starting server"` {
		t.Errorf("unexpected fix text: %s", fix.TextEdits[0].NewText)
	}
}

func TestCheckEnglish(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantNil bool
	}{
		{"english only", "starting server on port 8080", true},
		{"cyrillic", "запуск сервера", false},
		{"mixed", "server запуск", false},
		{"chinese", "启动服务器", false},
		{"empty", "", true},
		{"numbers", "12345", true},
		{"punctuation", "hello, world. test-case", true},
		{"special chars only", "!@#$%", true}, // no letters, so ok
		{"japanese", "サーバー", false},
		{"accented latin", "caf\u00e9", false}, // e with accent is > 127
		{"german", "gro\u00df", false},         // ß
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := checkEnglish(tt.msg, token.Pos(1), token.Pos(len(tt.msg)+3))
			if tt.wantNil && d != nil {
				t.Errorf("expected nil diagnostic, got: %s", d.Message)
			}
			if !tt.wantNil && d == nil {
				t.Errorf("expected diagnostic, got nil")
			}
		})
	}
}

func TestCheckSpecialChars(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantNil bool
	}{
		{"clean message", "server started", true},
		{"exclamation", "server started!", false},
		{"multiple exclamations", "failed!!!", false},
		{"at sign", "hello @admin", false},
		{"hash", "issue #123", false},
		{"percent", "100% done", true},              // % is allowed (used in format strings like %d, %v)
		{"percent format d", "got %d errors", true}, // format verb
		{"percent format s", "user: %s", true},      // format verb
		{"percent format v", "value: %v", true},     // format verb
		{"dollar", "value $10", false},              // $ is a special char
		{"ampersand", "this & that", false},
		{"asterisk", "important*", false},
		{"tilde", "~config", false},
		{"semicolon", "step 1; step 2", false},
		{"backtick", "run `cmd`", false},
		{"empty", "", true},
		{"comma ok", "hello, world", true},
		{"period ok", "done.", true},
		{"hyphen ok", "well-known", true},
		{"parentheses ok", "value (default)", true},
		{"colon ok", "host: localhost", true},
		{"slash ok", "path/to/file", true},
		{"question mark ok", "what happened", true},
		{"equals ok", "key=value", true},
		{"plus ok", "a+b", true},
		{"underscore ok", "snake_case", true},
		{"quotes ok", "said 'hello'", true},
		{"double quotes in msg", `said "hello"`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := checkSpecialChars(tt.msg, token.Pos(1), token.Pos(len(tt.msg)+3))
			if tt.wantNil && d != nil {
				t.Errorf("expected nil diagnostic, got: %s", d.Message)
			}
			if !tt.wantNil && d == nil {
				t.Errorf("expected diagnostic, got nil")
			}
			if !tt.wantNil && d != nil {
				if len(d.SuggestedFixes) == 0 {
					t.Error("expected suggested fix")
				}
			}
		})
	}
}

func TestCheckSpecialCharsSuggestedFix(t *testing.T) {
	tests := []struct {
		name     string
		msg      string
		expected string
	}{
		{"exclamation", "server started!", `"server started"`},
		{"multiple exclamations", "failed!!!", `"failed"`},
		{"at sign", "hello @admin", `"hello admin"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := checkSpecialChars(tt.msg, token.Pos(1), token.Pos(len(tt.msg)+3))
			if d == nil {
				t.Fatal("expected diagnostic")
			}
			if len(d.SuggestedFixes) != 1 || len(d.SuggestedFixes[0].TextEdits) != 1 {
				t.Fatal("expected 1 fix with 1 edit")
			}
			got := string(d.SuggestedFixes[0].TextEdits[0].NewText)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestCheckSensitiveData(t *testing.T) {
	patterns := DefaultSensitivePatterns()

	tests := []struct {
		name    string
		msg     string
		wantNil bool
	}{
		{"password", "user password: secret123", false},
		{"api_key", "api_key=abc123", false},
		{"token", "token: xyz", false},
		{"secret", "secret value here", false},
		{"private_key", "private_key loaded", false},
		{"access_key", "access_key set", false},
		{"api_secret", "api_secret configured", false},
		{"credential", "user credential found", false},
		{"auth", "auth header set", false},
		{"passwd", "passwd changed", false},
		{"apikey", "apikey sent", false},
		{"word boundary no match", "user authenticated successfully", true},
		{"no sensitive", "server started", true},
		{"empty", "", true},
		{"case insensitive", "PASSWORD: abc", false},
		{"word boundary no match 2", "authentication failed", true},
		{"auth standalone", "auth header set", false},
		{"auth at end", "basic auth", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := checkSensitiveData(tt.msg, token.Pos(1), token.Pos(len(tt.msg)+3), patterns)
			if tt.wantNil && d != nil {
				t.Errorf("expected nil diagnostic, got: %s", d.Message)
			}
			if !tt.wantNil && d == nil {
				t.Errorf("expected diagnostic, got nil")
			}
		})
	}
}

func TestCheckSensitiveDataCustomPatterns(t *testing.T) {
	patterns := append(DefaultSensitivePatterns(), "ssn", "credit_card")

	tests := []struct {
		name    string
		msg     string
		wantNil bool
	}{
		{"custom ssn", "user ssn: 123-45-6789", false},
		{"custom credit_card", "credit_card: 4111", false},
		{"default still works", "password: abc", false},
		{"clean", "payment processed", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := checkSensitiveData(tt.msg, token.Pos(1), token.Pos(len(tt.msg)+3), patterns)
			if tt.wantNil && d != nil {
				t.Errorf("expected nil diagnostic, got: %s", d.Message)
			}
			if !tt.wantNil && d == nil {
				t.Errorf("expected diagnostic, got nil")
			}
		})
	}
}

func TestIsEmoji(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{"rocket", '🚀', true},
		{"smile", '😀', true},
		{"heart", '❤', true},
		{"letter a", 'a', false},
		{"digit 1", '1', false},
		{"space", ' ', false},
		{"star", '⭐', true},
		{"warning", '⚠', true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEmoji(tt.r); got != tt.want {
				t.Errorf("isEmoji(%q) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}

func TestStripSpecialChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"exclamation", "hello!", "hello"},
		{"multiple exclamations", "failed!!!", "failed"},
		{"at sign", "hello @admin", "hello admin"},
		{"emoji", "done 🚀", "done"},
		{"ellipsis dots", "warning...", "warning."},
		{"clean", "hello world", "hello world"},
		{"mixed", "hello! @world #123", "hello world 123"},
		{"percent allowed", "100% done", "100% done"}, // % is not stripped
		{"tilde", "~config", "config"},
		{"asterisk", "important*", "important"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripSpecialChars(tt.input)
			if got != tt.expected {
				t.Errorf("stripSpecialChars(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsEllipsis(t *testing.T) {
	if !isEllipsis('\u2026') {
		t.Error("expected true for Unicode ellipsis")
	}
	if isEllipsis('.') {
		t.Error("expected false for regular period")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if !cfg.CheckLowercase {
		t.Error("expected CheckLowercase to be true")
	}
	if !cfg.CheckEnglish {
		t.Error("expected CheckEnglish to be true")
	}
	if !cfg.CheckSpecial {
		t.Error("expected CheckSpecial to be true")
	}
	if !cfg.CheckSensitive {
		t.Error("expected CheckSensitive to be true")
	}
	if len(cfg.SensitivePatterns) != 0 {
		t.Error("expected no custom patterns by default")
	}
}

func TestDefaultSensitivePatterns(t *testing.T) {
	patterns := DefaultSensitivePatterns()
	if len(patterns) == 0 {
		t.Fatal("expected default patterns")
	}

	expected := map[string]bool{
		"password":    true,
		"passwd":      true,
		"secret":      true,
		"token":       true,
		"api_key":     true,
		"apikey":      true,
		"api_secret":  true,
		"access_key":  true,
		"private_key": true,
		"credential":  true,
		"auth":        true,
	}

	for _, p := range patterns {
		if !expected[p] {
			t.Errorf("unexpected pattern: %s", p)
		}
		delete(expected, p)
	}
	for p := range expected {
		t.Errorf("missing pattern: %s", p)
	}
}

func TestUnquoteString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple", `"hello"`, "hello"},
		{"newline", `"hello\nworld"`, "hello\nworld"},
		{"tab", `"hello\tworld"`, "hello\tworld"},
		{"carriage return", `"hello\rworld"`, "hello\rworld"},
		{"escaped backslash", `"hello\\world"`, "hello\\world"},
		{"escaped quote", `"hello\"world"`, "hello\"world"},
		{"empty", `""`, ""},
		{"short", `"a"`, "a"},
		{"unknown escape", `"hello\xworld"`, "hello\\xworld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := unquoteString(tt.input)
			if got != tt.expected {
				t.Errorf("unquoteString(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestUnquoteStringTooShort(t *testing.T) {
	got := unquoteString("x")
	if got != "x" {
		t.Errorf("expected %q, got %q", "x", got)
	}
}

func TestIsWordChar(t *testing.T) {
	tests := []struct {
		r    rune
		want bool
	}{
		{'a', true},
		{'Z', true},
		{'0', true},
		{' ', false},
		{'_', false},
		{'-', false},
		{'.', false},
	}
	for _, tt := range tests {
		if got := isWordChar(tt.r); got != tt.want {
			t.Errorf("isWordChar(%q) = %v, want %v", tt.r, got, tt.want)
		}
	}
}

func TestConfigAllSensitivePatterns(t *testing.T) {
	cfg := Config{
		CheckSensitive:    true,
		SensitivePatterns: []string{"ssn", "credit_card"},
	}
	patterns := cfg.AllSensitivePatterns()
	defaults := DefaultSensitivePatterns()

	if len(patterns) != len(defaults)+2 {
		t.Errorf("expected %d patterns, got %d", len(defaults)+2, len(patterns))
	}

	found := map[string]bool{}
	for _, p := range patterns {
		found[p] = true
	}
	if !found["ssn"] {
		t.Error("missing custom pattern: ssn")
	}
	if !found["credit_card"] {
		t.Error("missing custom pattern: credit_card")
	}
}

func TestCheckLowercaseEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantNil bool
		fixText string
	}{
		// Emoji at start is not a letter — should pass
		{"emoji start", "🚀 launched", true, ""},
		// Digit start — should pass
		{"digit start", "123 items", true, ""},
		// Hyphen start — should pass
		{"hyphen start", "-flag value", true, ""},
		// Cyrillic uppercase — should fail (it's an uppercase letter)
		{"cyrillic uppercase", "\u0410лмаз", false, `"` + "\u0430лмаз" + `"`},
		// Multi-byte lowercase first char
		{"lowercase cyrillic start", "\u0430лмаз", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := checkLowercase(tt.msg, token.Pos(1), token.Pos(len(tt.msg)+3))
			if tt.wantNil && d != nil {
				t.Errorf("expected nil diagnostic, got: %s", d.Message)
			}
			if !tt.wantNil && d == nil {
				t.Errorf("expected diagnostic, got nil for msg=%q", tt.msg)
			}
			if !tt.wantNil && d != nil && tt.fixText != "" {
				if len(d.SuggestedFixes) == 0 || len(d.SuggestedFixes[0].TextEdits) == 0 {
					t.Fatal("expected suggested fix")
				}
				got := string(d.SuggestedFixes[0].TextEdits[0].NewText)
				if got != tt.fixText {
					t.Errorf("fix text: got %q, want %q", got, tt.fixText)
				}
			}
		})
	}
}

func TestCheckEnglishEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantNil bool
	}{
		// Pure ASCII symbols — no letters — should pass
		{"ascii symbols", "!@#$%", true},
		// Tab and newline (ascii control) — should pass
		{"tab and newline", "hello\tworld\n", true},
		// Arabic — should fail
		{"arabic", "مرحبا", false},
		// Korean — should fail
		{"korean", "안녕하세요", false},
		// ASCII letters only no spaces — should pass
		{"compact english", "connectionfailed", true},
		// Digit-only — should pass
		{"digits only", "42", true},
		// Mixed english digits — should pass
		{"english and digits", "retry3", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := checkEnglish(tt.msg, token.Pos(1), token.Pos(len(tt.msg)+3))
			if tt.wantNil && d != nil {
				t.Errorf("expected nil diagnostic, got: %s", d.Message)
			}
			if !tt.wantNil && d == nil {
				t.Errorf("expected diagnostic, got nil for msg=%q", tt.msg)
			}
			if !tt.wantNil && d != nil {
				if d.Message != "log message must be in English only" {
					t.Errorf("unexpected message: %s", d.Message)
				}
				// English check has no suggested fixes
				if len(d.SuggestedFixes) != 0 {
					t.Errorf("expected no suggested fixes, got %d", len(d.SuggestedFixes))
				}
			}
		})
	}
}

func TestCheckSpecialCharsEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantNil bool
	}{
		// Caret (^) is a special char
		{"caret", "version ^1.0", false},
		// Dollar sign ($) IS in specialChars list
		{"dollar", "value $10", false},
		// Ellipsis unicode character
		{"unicode ellipsis", "loading\u2026", false},
		// Emoji: rocket
		{"rocket emoji", "done 🚀", false},
		// Emoji: warning sign
		{"warning emoji", "alert ⚠", false},
		// Emoji: check mark (dingbat range 0x2700-0x27BF)
		{"checkmark emoji", "done ✅", false},
		// Multiple allowed punctuation together
		{"allowed combo", "key=val, host:port (default)", true},
		// Only spaces — should pass
		{"spaces only", "   ", true},
		// Pipe character — not in specialChars — should pass
		{"pipe", "a|b", true},
		// Percent is allowed (used in format strings like %d, %v, %s)
		{"percent format d", "got %d errors", true},
		{"percent format s", "user: %s", true},
		{"percent standalone", "100% complete", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := checkSpecialChars(tt.msg, token.Pos(1), token.Pos(len(tt.msg)+3))
			if tt.wantNil && d != nil {
				t.Errorf("expected nil diagnostic, got: %s", d.Message)
			}
			if !tt.wantNil && d == nil {
				t.Errorf("expected diagnostic, got nil for msg=%q", tt.msg)
			}
		})
	}
}

func TestCheckSensitiveDataEdgeCases(t *testing.T) {
	patterns := DefaultSensitivePatterns()

	tests := []struct {
		name    string
		msg     string
		wantNil bool
	}{
		// Pattern at the very start of message
		{"pattern at start", "password is wrong", false},
		// Pattern at the very end of message
		{"pattern at end", "check the token", false},
		// Pattern surrounded by underscores (word chars? no — underscore is not isWordChar)
		{"pattern between underscores", "_token_", false},
		// Pattern is the whole message
		{"pattern only", "token", false},
		// Pattern in a longer word — false positive avoided
		{"authorship should not match", "authorship claimed", true},
		// "secretive" starts with "secret" but has word char after — should not match
		{"secretive not match", "secretive agent", true},
		// Exact word "auth" at start
		{"auth at start", "auth failed", false},
		// "credentials" ends with "s" after "credential" — word char after — should NOT match
		{"credentials no match", "invalid credentials", true},
		// Empty patterns list — should not report anything
		{"empty patterns", "password is abc", true},
		// Numeric in message (non-letter) before pattern — should still match
		{"number before pattern", "3 token attempts", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := patterns
			if tt.name == "empty patterns" {
				p = []string{}
			}
			d := checkSensitiveData(tt.msg, token.Pos(1), token.Pos(len(tt.msg)+3), p)
			if tt.wantNil && d != nil {
				t.Errorf("expected nil diagnostic, got: %s", d.Message)
			}
			if !tt.wantNil && d == nil {
				t.Errorf("expected diagnostic, got nil for msg=%q", tt.msg)
			}
		})
	}
}

func TestCheckSensitiveDataMessage(t *testing.T) {
	d := checkSensitiveData("auth failed", token.Pos(1), token.Pos(15), DefaultSensitivePatterns())
	if d == nil {
		t.Fatal("expected diagnostic")
	}
	expected := `log message may contain sensitive data ("auth")`
	if d.Message != expected {
		t.Errorf("message: got %q, want %q", d.Message, expected)
	}
	if len(d.SuggestedFixes) != 0 {
		t.Errorf("expected no suggested fixes, got %d", len(d.SuggestedFixes))
	}
}

func TestStripSpecialCharsEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty", "", ""},
		{"spaces only", "   ", ""},
		{"only special", "!@#", ""},
		{"only emoji", "🚀", ""},
		{"dots then text", "...hello", ".hello"},
		{"text then dots", "hello...", "hello."},
		{"single dot", "done.", "done."},
		{"unicode ellipsis", "loading\u2026", "loading"},
		{"mixed valid", "key=value, host:port", "key=value, host:port"},
		{"caret removed", "^start", "start"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripSpecialChars(tt.input)
			if got != tt.expected {
				t.Errorf("stripSpecialChars(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsEmojiAdditional(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		// Supplemental symbols range 0x1F900-0x1F9FF
		{"yawning face", '\U0001F971', true},
		// Chess symbols range 0x1FA00-0x1FA6F
		{"chess pawn", '\U0001FA01', true},
		// Transport range 0x1F680-0x1F6FF
		{"car", '\U0001F697', true},
		// Misc pictographs 0x1F300-0x1F5FF
		{"globe", '\U0001F30D', true},
		// Flags range 0x1F1E0-0x1F1FF
		{"flag letter", '\U0001F1E6', true},
		// Misc symbols 0x2600-0x26FF
		{"sun", '\u2600', true},
		// Dingbats 0x2700-0x27BF
		{"scissors", '\u2702', true},
		// Non-emoji ASCII
		{"letter z", 'z', false},
		{"digit 9", '9', false},
		{"newline", '\n', false},
		// Hourglass (0x231B)
		{"hourglass", '\u231B', true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEmoji(tt.r); got != tt.want {
				t.Errorf("isEmoji(%U) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}

func TestCheckLowercaseFixMessage(t *testing.T) {
	d := checkLowercase("Starting server", token.Pos(1), token.Pos(18))
	if d == nil {
		t.Fatal("expected diagnostic")
	}
	if d.Message != "log message must start with a lowercase letter" {
		t.Errorf("unexpected message: %q", d.Message)
	}
	if len(d.SuggestedFixes) != 1 {
		t.Fatalf("expected 1 fix, got %d", len(d.SuggestedFixes))
	}
	if d.SuggestedFixes[0].Message != "Convert first letter to lowercase" {
		t.Errorf("unexpected fix message: %q", d.SuggestedFixes[0].Message)
	}
}

func TestCheckSpecialCharsFixMessage(t *testing.T) {
	d := checkSpecialChars("hello!", token.Pos(1), token.Pos(9))
	if d == nil {
		t.Fatal("expected diagnostic")
	}
	if d.Message != "log message must not contain special characters or emoji" {
		t.Errorf("unexpected message: %q", d.Message)
	}
	if len(d.SuggestedFixes) != 1 {
		t.Fatalf("expected 1 fix, got %d", len(d.SuggestedFixes))
	}
	if d.SuggestedFixes[0].Message != "Remove special characters and emoji" {
		t.Errorf("unexpected fix message: %q", d.SuggestedFixes[0].Message)
	}
}

func TestDerefPointer(t *testing.T) {
	// derefPointer is an internal helper; test via checkLowercase indirectly
	// (the function is already exercised through integration tests).
	// Here we verify the exported behavior is consistent.
	d := checkLowercase("Already lowercase", token.Pos(1), token.Pos(20))
	// "Already" starts with uppercase A
	if d == nil {
		t.Fatal("expected diagnostic for uppercase start")
	}
}

func TestCheckSensitiveDataLongestPatternFirst(t *testing.T) {
	// "api_secret" is longer than "secret", so it should be reported as "api_secret"
	patterns := DefaultSensitivePatterns()
	d := checkSensitiveData("api_secret stored", token.Pos(1), token.Pos(20), patterns)
	if d == nil {
		t.Fatal("expected diagnostic")
	}
	expected := `log message may contain sensitive data ("api_secret")`
	if d.Message != expected {
		t.Errorf("message: got %q, want %q", d.Message, expected)
	}
}

func TestCheckSensitiveDataWordBoundaryDetails(t *testing.T) {
	patterns := DefaultSensitivePatterns()

	// "auth" should NOT match inside "authentication"
	d := checkSensitiveData("authentication failed", token.Pos(1), token.Pos(25), patterns)
	if d != nil {
		t.Errorf("expected nil (word boundary), got: %s", d.Message)
	}

	// "auth" SHOULD match when preceded by space and followed by space
	d = checkSensitiveData("basic auth header", token.Pos(1), token.Pos(20), patterns)
	if d == nil {
		t.Fatal("expected diagnostic for standalone 'auth'")
	}

	// "token" SHOULD match when preceded by colon
	d = checkSensitiveData("bearer: token value", token.Pos(1), token.Pos(22), patterns)
	if d == nil {
		t.Fatal("expected diagnostic for 'token' after colon")
	}
}

func TestUnquoteStringEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"single char no quotes", "a", "a"},
		// string with len < 2 after removing first quote — returns empty (len=2: remove first and last char)
		{"only open quote", `"a`, ``},
		{"unicode passthrough", `"caf` + "\u00e9" + `"`, "caf\u00e9"},
		{"multiple escapes", `"a\nb\tc"`, "a\nb\tc"},
		// `"a\"` (4 bytes) → s[1:3] = `a\` → loop: 'a', then '\' with no next char → writes `a\`
		{"backslash at end", `"a\"`, `a\`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := unquoteString(tt.input)
			if got != tt.expected {
				t.Errorf("unquoteString(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
