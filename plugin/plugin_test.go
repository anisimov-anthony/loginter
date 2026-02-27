package plugin

import (
	"testing"

	"github.com/golangci/plugin-module-register/register"
)

func TestNew(t *testing.T) {
	p, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil plugin")
	}
}

func TestNewWithSettings(t *testing.T) {
	settings := map[string]any{
		"check_lowercase":    true,
		"check_english":      false,
		"check_special":      true,
		"check_sensitive":    true,
		"sensitive_patterns": []any{"ssn", "credit_card"},
	}
	p, err := New(settings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil plugin")
	}
}

func TestNewDefaultsWhenAllDisabled(t *testing.T) {
	// When all checks are disabled (zero values), defaults should be applied
	settings := map[string]any{
		"check_lowercase": false,
		"check_english":   false,
		"check_special":   false,
		"check_sensitive": false,
	}
	p, err := New(settings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil plugin")
	}
}

func TestBuildAnalyzers(t *testing.T) {
	p, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	analyzers, err := p.BuildAnalyzers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(analyzers) != 1 {
		t.Fatalf("expected 1 analyzer, got %d", len(analyzers))
	}
	if analyzers[0].Name != "loginter" {
		t.Errorf("expected analyzer name 'loginter', got %q", analyzers[0].Name)
	}
}

func TestGetLoadMode(t *testing.T) {
	p, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mode := p.GetLoadMode()
	if mode != register.LoadModeTypesInfo {
		t.Errorf("expected %q, got %q", register.LoadModeTypesInfo, mode)
	}
}

func TestPluginRegistered(t *testing.T) {
	newPlugin, err := register.GetPlugin("loginter")
	if err != nil {
		t.Fatalf("plugin not registered: %v", err)
	}
	if newPlugin == nil {
		t.Fatal("expected non-nil plugin constructor")
	}
}
