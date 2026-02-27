//go:build e2e

// Package e2e runs end-to-end tests by building the standalone loginter binary
// and executing it against real Go source code, verifying the exit code and
// the diagnostic messages in its output.
package e2e_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// binaryPath holds the path to the compiled loginter binary for the test session.
var binaryPath string

// TestMain builds the loginter binary once before running all e2e tests.
func TestMain(m *testing.M) {
	bin, err := buildBinary()
	if err != nil {
		panic("failed to build loginter binary: " + err.Error())
	}
	binaryPath = bin
	defer os.Remove(binaryPath)

	os.Exit(m.Run())
}

// buildBinary compiles cmd/loginter into a temp file and returns its path.
func buildBinary() (string, error) {
	_, thisFile, _, _ := runtime.Caller(0)
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..")

	bin, err := os.CreateTemp("", "loginter-e2e-*")
	if err != nil {
		return "", err
	}
	bin.Close()

	cmd := exec.Command("go", "build", "-o", bin.Name(), "./cmd/loginter")
	cmd.Dir = repoRoot
	if out, buildErr := cmd.CombinedOutput(); buildErr != nil {
		return "", fmt.Errorf("build failed: %w\n%s", buildErr, out)
	}
	return bin.Name(), nil
}

// testdataDir returns an absolute path to a named testcase directory.
func testdataDir(name string) string {
	_, thisFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(thisFile), "testcases", name)
}

// runLinter executes the binary on the given package directory.
// It returns (stdout+stderr combined, exit code).
func runLinter(t *testing.T, dir string) (string, int) {
	t.Helper()
	cmd := exec.Command(binaryPath, "./...")
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	code := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code = exitErr.ExitCode()
		} else {
			t.Fatalf("unexpected error running linter: %v", err)
		}
	}
	return buf.String(), code
}

// TestExitCodeNonZeroOnViolations verifies that the linter exits with a non-zero
// code when violations are present.
func TestExitCodeNonZeroOnViolations(t *testing.T) {
	_, code := runLinter(t, testdataDir("violations"))
	if code == 0 {
		t.Error("expected non-zero exit code when violations are present, got 0")
	}
}

// TestRule1Lowercase verifies that uppercase-starting messages are reported.
func TestRule1Lowercase(t *testing.T) {
	out, _ := runLinter(t, testdataDir("violations"))
	want := []string{
		"log message must start with a lowercase letter",
	}
	assertContains(t, out, want)
	// 5 uppercase messages: lines 14-18
	assertLineNumbers(t, out, "log message must start with a lowercase letter", []int{14, 15, 16, 17, 18})
}

// TestRule2English verifies that non-English messages are reported.
func TestRule2English(t *testing.T) {
	out, _ := runLinter(t, testdataDir("violations"))
	want := []string{
		"log message must be in English only",
	}
	assertContains(t, out, want)
	// 3 non-English messages: lines 21-23
	assertLineNumbers(t, out, "log message must be in English only", []int{21, 22, 23})
}

// TestRule3SpecialChars verifies that special characters and emoji are reported.
func TestRule3SpecialChars(t *testing.T) {
	out, _ := runLinter(t, testdataDir("violations"))
	want := []string{
		"log message must not contain special characters or emoji",
	}
	assertContains(t, out, want)
}

// TestRule3SpecificChars verifies individual special characters are caught (by line).
func TestRule3SpecificChars(t *testing.T) {
	out, _ := runLinter(t, testdataDir("violations"))
	// Lines 26-34 each contain a different special character / emoji
	assertLineNumbers(t, out, "log message must not contain special characters or emoji", []int{26, 27, 28, 29, 30, 31, 32, 33, 34})
}

// TestRule4SensitiveData verifies that sensitive keyword patterns are reported.
func TestRule4SensitiveData(t *testing.T) {
	out, _ := runLinter(t, testdataDir("violations"))
	want := []string{
		`log message may contain sensitive data ("password")`,
		`log message may contain sensitive data ("token")`,
		`log message may contain sensitive data ("secret")`,
		`log message may contain sensitive data ("api_key")`,
		`log message may contain sensitive data ("api_secret")`,
		`log message may contain sensitive data ("private_key")`,
		`log message may contain sensitive data ("access_key")`,
		`log message may contain sensitive data ("credential")`,
		`log message may contain sensitive data ("auth")`,
		`log message may contain sensitive data ("passwd")`,
		`log message may contain sensitive data ("apikey")`,
	}
	assertContains(t, out, want)
}

// TestExitCodeZeroOnCleanCode verifies that the linter exits 0 for compliant code.
func TestExitCodeZeroOnCleanCode(t *testing.T) {
	out, code := runLinter(t, testdataDir("clean"))
	if code != 0 {
		t.Errorf("expected exit code 0 for clean code, got %d\noutput:\n%s", code, out)
	}
}

// TestNoFalsePositivesOnCleanCode verifies there are no diagnostics for clean code.
func TestNoFalsePositivesOnCleanCode(t *testing.T) {
	out, _ := runLinter(t, testdataDir("clean"))
	forbidden := []string{
		"log message must start with a lowercase letter",
		"log message must be in English only",
		"log message must not contain special characters or emoji",
		"log message may contain sensitive data",
	}
	assertNotContains(t, out, forbidden)
}

// TestWordBoundaryNoFalsePositives verifies word-boundary logic on the sensitive rule.
// "authenticated" must not trigger "auth", "authorization" must not trigger "auth".
func TestWordBoundaryNoFalsePositives(t *testing.T) {
	out, _ := runLinter(t, testdataDir("clean"))
	if strings.Contains(out, "sensitive data") {
		t.Errorf("word boundary false positive in clean code:\n%s", out)
	}
}

// TestPercentAllowedInMessages verifies that %d / %v / %s are not flagged.
func TestPercentAllowedInMessages(t *testing.T) {
	out, _ := runLinter(t, testdataDir("clean"))
	if strings.Contains(out, "special characters") {
		t.Errorf("percent sign falsely flagged as special char:\n%s", out)
	}
}

func assertContains(t *testing.T, output string, patterns []string) {
	t.Helper()
	for _, p := range patterns {
		if !strings.Contains(output, p) {
			t.Errorf("expected linter output to contain %q\nfull output:\n%s", p, output)
		}
	}
}

func assertNotContains(t *testing.T, output string, patterns []string) {
	t.Helper()
	for _, p := range patterns {
		if strings.Contains(output, p) {
			t.Errorf("linter output must NOT contain %q\nfull output:\n%s", p, output)
		}
	}
}

// assertLineNumbers checks that each line number appears in the output paired with the given message.
// Output format: "file.go:LINE:COL: message"
func assertLineNumbers(t *testing.T, output, message string, lines []int) {
	t.Helper()
	for _, ln := range lines {
		needle := fmt.Sprintf(":%d:", ln)
		found := false
		for _, line := range strings.Split(output, "\n") {
			if strings.Contains(line, needle) && strings.Contains(line, message) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected diagnostic %q at line %d\nfull output:\n%s", message, ln, output)
		}
	}
}
