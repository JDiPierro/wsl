package wsl

import (
	"go/token"

	"golang.org/x/tools/go/analysis"
)

func NewAnalyzerWithProcessor(p *Processor) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:             "wsl",
		Doc:              "add or remove empty lines",
		Run:              runFunc(p),
		RunDespiteErrors: true,
	}
}

func runFunc(p *Processor) func(pass *analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		p.ProcessAstFiles(pass.Fset, pass.Files)

		for _, err := range p.result {
			var (
				newText []byte
				pos     token.Pos
				end     token.Pos
			)

			switch err.Type {
			case ShouldAddWS:
				pos = err.PreviousEndPos
				end = err.PreviousEndPos
				newText = []byte("\n")
			case ShouldRemoveWS:
				// TODO: This is not yet implemented.
				continue

				// nolint: govet: unreachable
				// Leading whitespaces
				pos = err.ErrorNode.Pos()
				end = err.ErrorNode.Pos()

				// Trailing whitespaces
				// pos = err.End - 2
				// end = err.End - 2

				newText = []byte("")
			case InvalidFile:
				// Nothing to do
				continue
			}

			d := analysis.Diagnostic{
				Pos:      pos,
				Category: "",
				Message:  err.Reason,
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: err.Type.String(),
						TextEdits: []analysis.TextEdit{
							{
								Pos:     pos,
								End:     end,
								NewText: newText,
							},
						},
					},
				},
			}

			pass.Report(d)
		}

		return nil, nil
	}
}
