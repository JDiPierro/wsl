package main

import (
	"github.com/bombsimon/wsl/v2"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	p := wsl.NewProcessor()

	singlechecker.Main(
		wsl.NewAnalyzerWithProcessor(p),
	)
}
