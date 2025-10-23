package forest

import (
	"context"
	"fmt"
	"maps"
	"math"
	"slices"
	"time"
	lsproto "vocab/lsp"
	"vocab/vocabulary/parser"
)

// The compiler
type Forest struct {
	ctx context.Context
	// Map of document uri and the associated diagnostics from parser
	parsingDiagnostics map[string][]*lsproto.Diagnostic
	// Map of document uri and the associated trees
	trees map[string]*WordTree
	log   func(any)
}

func NewForest(ctx context.Context, log func(any)) *Forest {
	return &Forest{
		parsingDiagnostics: make(map[string][]*lsproto.Diagnostic),
		trees:              make(map[string]*WordTree),
		ctx:                ctx,
		log:                log,
	}
}

// Create or replace tree associated with documentUri and merge it back to the global tree.
//
// This also clears the diagnostics of the current documentUri
func (c *Forest) Plant(documentUri string, text string, changeRange *lsproto.Range) *Forest {
	scanner := parser.NewScanner(text)
	parser := parser.NewParser(c.ctx, documentUri, scanner, c.log)
	parser.Parse()

	c.parsingDiagnostics[documentUri] = []*lsproto.Diagnostic{}
	for _, section := range parser.Ast.Sections {
		c.parsingDiagnostics[documentUri] = append(c.parsingDiagnostics[documentUri], section.Diagnostics...)
	}

	c.trees[documentUri] = AstToWordTree(parser.Ast)

	return c
}

// Based on the built tree, compile tree into diagnostics.
func (c *Forest) Harvest() map[string][]lsproto.Diagnostic {
	mergedTree := NewWordTree()
	for _, tree := range c.trees {
		mergedTree.Graft(tree)
	}
	fruits := mergedTree.Harvest()

	diags := make(map[string][]lsproto.Diagnostic)
	for uri := range c.trees {
		diags[uri] = []lsproto.Diagnostic{}
	}

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

			diags[word.Uri()] = append(diags[word.Uri()], *err)
		}
	}

	for _, fruit := range fruits {
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

	for uri := range c.parsingDiagnostics {
		for _, diag := range c.parsingDiagnostics[uri] {
			diags[uri] = append(diags[uri], *diag)
		}
	}

	return diags
}

func (f *Forest) GetTreesLocations() []string {
	return slices.Collect(maps.Keys(f.trees))
}
