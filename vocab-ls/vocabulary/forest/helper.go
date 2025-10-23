package forest

import (
	lsproto "vocab/lsp"
	"vocab/vocabulary/parser"
)

func AstToWordTree(ast *parser.VocabAst) *WordTree {
	tree := NewWordTree()

	for _, section := range ast.Sections {
		allWordSections := append(section.NewWords, section.ReviewedWords...)
		for _, wordSection := range allWordSections {
			lang := wordSection.Language
			for _, word := range wordSection.Words {
				tree.AddTwig(lang, word, ast.Uri, section, []*lsproto.Diagnostic{})
			}
		}
	}
	return tree
}
