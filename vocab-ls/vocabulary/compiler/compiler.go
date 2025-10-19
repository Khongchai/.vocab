package compiler

import (
	"context"
	"fmt"
	"math"
	"time"
	lsproto "vocab/lsp"
	"vocab/vocabulary/parser"
)

type Compiler struct {
	ctx  context.Context
	tree *WordTree
	log  func(any)
}

func NewCompiler(ctx context.Context, log func(any)) *Compiler {
	return &Compiler{
		ctx:  ctx,
		tree: nil,
		log:  log,
	}
}

// This is called "accept" not compile because the compiler can incrementally build the inner language tree representation.
// Calling "Accept" multiple times will update the existing global tree state.
func (c *Compiler) Accept(documentUri string, text string, changeRange *lsproto.Range) {
	scanner := parser.NewScanner(text)
	parser := parser.NewParser(c.ctx, documentUri, scanner, c.log)
	parser.Parse()

	newWordTree := AstToWordTree(parser.Ast)
	if c.tree == nil {
		c.tree = newWordTree
	} else {
		c.tree.Graft(newWordTree)
	}
}

// Based on the built tree, compile tree into diagnostics.
func (c *Compiler) Compile() []lsproto.Diagnostic {
	details := c.tree.Harvest()
	diags := []lsproto.Diagnostic{}

	addDiagToWords := func(timeRemaining float64, severitiy lsproto.DiagnosticsSeverity, words []*parser.Word) {
		for _, word := range words {
			err := lsproto.MakeDiagnostics(
				fmt.Sprintf("Needs review within: %f day(s)", timeRemaining),
				word.Line,
				word.Start,
				word.End,
				severitiy,
			)

			diags = append(diags, *err)
		}
	}

	for _, detail := range details {
		for _, starting := range detail.StartingDiagnostics {
			diags = append(diags, *starting)
		}

		severity, remainingDays := func() (lsproto.DiagnosticsSeverity, float64) {
			interval := math.Ceil(detail.Interval)
			deadline := detail.LastSeenDate.AddDate(0, 0, int(interval))
			remainingHours := time.Until(deadline).Hours()
			var remainingDays float64 = remainingHours / 24

			if remainingDays <= 1 {
				return lsproto.DiagnosticsSeverityError, remainingDays
			} else if remainingDays < 3 {
				return lsproto.DiagnosticsSeverityHint, remainingDays
			}

			return lsproto.DiagnosticsSeverityInformation, remainingDays
		}()

		addDiagToWords(remainingDays, severity, detail.Words)
	}

	return diags
}
