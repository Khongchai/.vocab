package compiler

import (
	"fmt"
	lsproto "vocab/lsp"
	"vocab/vocabulary/parser"
)

func AstToWordTree(ast *parser.VocabAst) *WordTree {
	tree := NewWordTree()

	for _, section := range ast.Sections {
		allWordSections := append(section.NewWords, section.ReviewedWords...)
		for _, wordSection := range allWordSections {
			lang := wordSection.Language
			reviewed := wordSection.Reviewed
			for _, word := range wordSection.Words {
				var diag []*lsproto.Diagnostic
				existingTwigs := tree.GetTwigs(lang, word.Text)

				warnText := func() string {
					if !reviewed && len(existingTwigs) != 0 {
						return fmt.Sprintf("This has been seen before. First occurence in %s", existingTwigs[0].GetLocation())
					} else if reviewed && len(existingTwigs) == 0 {
						return "This has never been seen before"
					}
					return ""
				}()

				if warnText != "" {
					warn := lsproto.MakeDiagnostics(
						warnText,
						wordSection.Line,
						word.Start,
						word.End,
						lsproto.DiagnosticsSeverityWarning,
					)
					diag = append(diag, warn)
				}

				tree.AddTwig(lang, word, ast.Uri, section, diag)
			}
		}
	}
	return tree
}
