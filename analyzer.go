package wsl

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

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
				pos = err.ErrorNode.Pos()
				end = err.ErrorNode.Pos()
				// newText = newlineAndIndent(err.ErrorNode, err.PreviousNode)
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

// We must calculate how much the text is indented to get the number of tabs we
// must indent the line after applying the newline.
func newlineAndIndent(n1, n2 ast.Node) []byte {
	var (
		// TODO: Not reliable with comments
		posDiff = n1.Pos() - n2.End() - 1 // -1 to remove the newline
		tabs    = make([]string, posDiff)
	)

	for i := range tabs {
		tabs[i] = "\t"
	}

	return []byte(fmt.Sprintf(
		"\n%s", strings.Join(tabs, ""),
	))
}
