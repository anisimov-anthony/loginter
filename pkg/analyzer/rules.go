package analyzer

import (
	"go/token"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

// checkLowercase checks that the log message starts with a lowercase letter.
// Returns a diagnostic with a suggested fix if the message starts with an uppercase letter.
func checkLowercase(msg string, pos, end token.Pos) *analysis.Diagnostic {
	if len(msg) == 0 {
		return nil
	}
	r, size := utf8.DecodeRuneInString(msg)
	if r == utf8.RuneError || !unicode.IsUpper(r) {
		return nil
	}

	fixed := string(unicode.ToLower(r)) + msg[size:]

	return &analysis.Diagnostic{
		Pos:     pos,
		End:     end,
		Message: "log message must start with a lowercase letter",
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: "Convert first letter to lowercase",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     pos,
						End:     end,
						NewText: []byte(`"` + fixed + `"`),
					},
				},
			},
		},
	}
}

// checkEnglish checks that the log message contains only English (ASCII) characters.
// Letters outside the ASCII range (e.g. Cyrillic, Chinese, etc.) are flagged.
// Digits, spaces, and common punctuation are allowed.
func checkEnglish(msg string, pos, end token.Pos) *analysis.Diagnostic {
	for _, r := range msg {
		if unicode.IsLetter(r) && r > 127 {
			return &analysis.Diagnostic{
				Pos:     pos,
				End:     end,
				Message: "log message must be in English only",
			}
		}
	}
	return nil
}

// specialChars is the set of characters considered "special" in log messages.
const specialChars = "!@#$^&*~;`"

// checkSpecialChars checks that the log message does not contain special characters or emoji.
// Common punctuation used in sentences (comma, period, hyphen, apostrophe, colon, parentheses,
// slash, question mark) is allowed.
func checkSpecialChars(msg string, pos, end token.Pos) *analysis.Diagnostic {
	for _, r := range msg {
		if strings.ContainsRune(specialChars, r) || isEmoji(r) || isEllipsis(r) {
			fixed := stripSpecialChars(msg)
			return &analysis.Diagnostic{
				Pos:     pos,
				End:     end,
				Message: "log message must not contain special characters or emoji",
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: "Remove special characters and emoji",
						TextEdits: []analysis.TextEdit{
							{
								Pos:     pos,
								End:     end,
								NewText: []byte(`"` + fixed + `"`),
							},
						},
					},
				},
			}
		}
	}
	return nil
}

// checkSensitiveData checks that the log message does not contain sensitive data keywords.
// It checks both the message string literal and any string concatenation with variables.
// Patterns are checked longest-first to prefer more specific matches (e.g. "api_secret" over "secret").
func checkSensitiveData(msg string, pos, end token.Pos, patterns []string) *analysis.Diagnostic {
	lower := strings.ToLower(msg)

	sorted := make([]string, len(patterns))
	copy(sorted, patterns)
	sort.Slice(sorted, func(i, j int) bool {
		return len(sorted[i]) > len(sorted[j])
	})

	for _, pattern := range sorted {
		lp := strings.ToLower(pattern)
		idx := strings.Index(lower, lp)
		if idx == -1 {
			continue
		}
		before := idx - 1
		after := idx + len(lp)
		if before >= 0 && isWordChar(rune(lower[before])) {
			continue
		}
		if after < len(lower) && isWordChar(rune(lower[after])) {
			continue
		}
		return &analysis.Diagnostic{
			Pos:     pos,
			End:     end,
			Message: "log message may contain sensitive data (\"" + pattern + "\")",
		}
	}
	return nil
}

// isEmoji returns true if the rune is an emoji character.
func isEmoji(r rune) bool {
	return (r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols and Pictographs
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map
		(r >= 0x1F1E0 && r <= 0x1F1FF) || // Flags
		(r >= 0x2600 && r <= 0x26FF) || // Misc symbols
		(r >= 0x2700 && r <= 0x27BF) || // Dingbats
		(r >= 0xFE00 && r <= 0xFE0F) || // Variation Selectors
		(r >= 0x1F900 && r <= 0x1F9FF) || // Supplemental Symbols
		(r >= 0x1FA00 && r <= 0x1FA6F) || // Chess Symbols
		(r >= 0x1FA70 && r <= 0x1FAFF) || // Symbols Extended-A
		(r >= 0x200D && r <= 0x200D) || // Zero Width Joiner
		(r >= 0x231A && r <= 0x231B) || // Watch, Hourglass
		(r >= 0x23E9 && r <= 0x23F3) || // Various symbols
		(r >= 0x23F8 && r <= 0x23FA) || // Various symbols
		(r >= 0x2B50 && r <= 0x2B55) // Stars and circles
}

// isWordChar returns true if the rune is a letter or digit (part of a "word").
func isWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isEllipsis checks if the rune is the Unicode ellipsis character.
func isEllipsis(r rune) bool {
	return r == '\u2026' // …
}

// stripSpecialChars removes special characters and emoji from a string.
func stripSpecialChars(s string) string {
	var b strings.Builder
	dotCount := 0
	for _, r := range s {
		if r == '.' {
			dotCount++
			if dotCount > 1 {
				continue
			}
			b.WriteRune(r)
			continue
		}
		dotCount = 0

		if strings.ContainsRune(specialChars, r) || isEmoji(r) || isEllipsis(r) {
			continue
		}
		b.WriteRune(r)
	}
	return strings.TrimSpace(b.String())
}
