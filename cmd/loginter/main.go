package main

import (
	"github.com/anisimov-anthony/loginter/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.NewAnalyzer(analyzer.DefaultConfig()))
}
