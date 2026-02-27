package plugin

import (
	"github.com/anisimov-anthony/loginter/pkg/analyzer"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("loginter", New)
}

// New creates a new loginter plugin with the given settings.
func New(settings any) (register.LinterPlugin, error) {
	cfg, err := register.DecodeSettings[analyzer.Config](settings)
	if err != nil {
		return nil, err
	}

	if !cfg.CheckLowercase && !cfg.CheckEnglish && !cfg.CheckSpecial && !cfg.CheckSensitive {
		cfg = analyzer.DefaultConfig()
	}

	return &loginterPlugin{cfg: cfg}, nil
}

type loginterPlugin struct {
	cfg analyzer.Config
}

func (p *loginterPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.NewAnalyzer(p.cfg)}, nil
}

func (p *loginterPlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
