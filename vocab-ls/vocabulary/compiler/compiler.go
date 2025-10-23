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
	ctx                context.Context
	parsingDiagnostics map[string][]*lsproto.Diagnostic
	tree               *WordTree
	log                func(any)
}

func NewCompiler(ctx context.Context, log func(any)) *Compiler {
	return &Compiler{
		parsingDiagnostics: make(map[string][]*lsproto.Diagnostic),
		ctx:                ctx,
		tree:               nil,
		log:                log,
	}
}

// This is called "accept" not compile because the compiler can incrementally build the inner language tree representation.
// Calling "Accept" multiple times will update the existing global tree state.
func (c *Compiler) Accept(documentUri string, text string, changeRange *lsproto.Range) *Compiler {

	scanner := parser.NewScanner(text)
	parser := parser.NewParser(c.ctx, documentUri, scanner, c.log)
	parser.Parse()

	c.parsingDiagnostics[documentUri] = []*lsproto.Diagnostic{}
	for _, section := range parser.Ast.Sections {
		c.parsingDiagnostics[documentUri] = append(c.parsingDiagnostics[documentUri], section.Diagnostics...)
	}

	newWordTree := AstToWordTree(parser.Ast)
	if c.tree == nil {
		c.tree = newWordTree
	} else {
		c.tree.Graft(newWordTree)
	}

	return c
}

// Based on the built tree, compile tree into diagnostics.
func (c *Compiler) Compile() []lsproto.Diagnostic {
	fruits := c.tree.Harvest()
	diags := []lsproto.Diagnostic{}

	addDiagToAllWordPositions := func(timeRemaining float64, severitiy lsproto.DiagnosticsSeverity, words []*parser.Word) {
		for _, word := range words {
			message := func() string {
				if timeRemaining == 0 {
					return "Review now!"
				}
				// can keep this for hover action
				if timeRemaining > 0 {
					return ""
				}
				return fmt.Sprintf("%d days past deadline", int(math.Ceil(timeRemaining*-1)))
			}()
			if message == "" {
				continue
			}
			err := lsproto.MakeDiagnostics(
				message,
				word.Line,
				word.Start,
				word.End,
				severitiy,
			)

			diags = append(diags, *err)
		}
	}

	for _, fruit := range fruits {
		for _, starting := range fruit.StartingDiagnostics {
			diags = append(diags, *starting)
		}

		severity, remainingDays := func() (lsproto.DiagnosticsSeverity, float64) {
			interval := math.Ceil(fruit.Interval)
			deadline := fruit.LastSeenDate.AddDate(0, 0, int(interval))
			remainingHours := time.Until(deadline).Hours()
			var remainingDays float64 = remainingHours / 24

			if remainingDays <= 1 {
				return lsproto.DiagnosticsSeverityError, remainingDays
			} else if remainingDays < 3 {
				return lsproto.DiagnosticsSeverityHint, remainingDays
			}

			return lsproto.DiagnosticsSeverityInformation, remainingDays
		}()

		addDiagToAllWordPositions(remainingDays, severity, fruit.Words)
	}

	for key := range c.parsingDiagnostics {
		for _, parsingDiag := range c.parsingDiagnostics[key] {
			diags = append(diags, *parsingDiag)
		}
	}

	return diags
}
