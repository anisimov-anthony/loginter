package analyzer

// Config holds the configuration for the loginter analyzer.
type Config struct {
	// CheckLowercase enables checking that log messages start with a lowercase letter.
	CheckLowercase bool `json:"check_lowercase"`

	// CheckEnglish enables checking that log messages contain only English characters.
	CheckEnglish bool `json:"check_english"`

	// CheckSpecial enables checking that log messages do not contain special characters or emoji.
	CheckSpecial bool `json:"check_special"`

	// CheckSensitive enables checking that log messages do not contain sensitive data.
	CheckSensitive bool `json:"check_sensitive"`

	// SensitivePatterns is a list of additional patterns to check for sensitive data.
	// These are added to the default list of sensitive keywords.
	SensitivePatterns []string `json:"sensitive_patterns"`
}

// DefaultConfig returns a Config with all checks enabled and default sensitive patterns.
func DefaultConfig() Config {
	return Config{
		CheckLowercase: true,
		CheckEnglish:   true,
		CheckSpecial:   true,
		CheckSensitive: true,
	}
}

// DefaultSensitivePatterns returns the built-in list of sensitive keywords.
func DefaultSensitivePatterns() []string {
	return []string{
		"password",
		"passwd",
		"secret",
		"token",
		"api_key",
		"apikey",
		"api_secret",
		"access_key",
		"private_key",
		"credential",
		"auth",
	}
}

// AllSensitivePatterns returns the combined list of default and custom sensitive patterns.
func (c Config) AllSensitivePatterns() []string {
	patterns := DefaultSensitivePatterns()
	patterns = append(patterns, c.SensitivePatterns...)
	return patterns
}
