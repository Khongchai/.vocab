package compiler

import (
	"testing"
	test "vocab/vocab_testing"
	"vocab/vocabulary/parser"
)

func TestAstToWordTree(t *testing.T) {
	text := test.TrimLines(`
		20/08/2025	
		> (it) sfatare, il mito
		Non per sfatare il mito della città eterna, della Roma magica che vedete anche nei film, soprattutto nei film di Fellini, ma per mostrarvi la verità o comunque la vera realtà di Roma e dei romani. 
	`)
	parser := parser.NewParser(t.Context(), "xxx", parser.NewScanner(text), func(a any) {})
	parser.Parse()

}
