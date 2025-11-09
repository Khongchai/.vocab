package forest

import (
	"context"
	"fmt"
	"maps"
	"math"
	"slices"
	"sync"
	"time"
	lib "vocab/lib"
	lsproto "vocab/lsp"
	"vocab/vocabulary/parser"
)

// The compiler
type Forest struct {
	ctx context.Context
	// Map of document uri and the associated diagnostics from parser
	parsingDiagnostics map[string][]*lsproto.Diagnostic
	// Map of document uri and the associated trees
	trees        map[string]*WordTree
	log          func(any)
	pool         *lib.GoWorkerPool
	harvestMutex sync.Mutex
}

func NewForest(ctx context.Context, log func(any)) *Forest {
	return &Forest{
		parsingDiagnostics: make(map[string][]*lsproto.Diagnostic),
		trees:              make(map[string]*WordTree),
		ctx:                ctx,
		log:                log,
		pool:               lib.NewGoWorkerPool(ctx),
	}
}

// Create or replace tree associated with documentUri and merge it back to the global tree.
//
// # This also clears the diagnostics of the current documentUri
//
// This method spawns a new thread if available and parse the given file.
func (c *Forest) Plant(documentUri string, text string, changeRange *lsproto.Range) *Forest {
	c.pool.Run(documentUri, func() {
		scanner := parser.NewScanner(text)
		parser := parser.NewParser(c.ctx, documentUri, scanner, c.log)
		parser.Parse()

		c.parsingDiagnostics[documentUri] = []*lsproto.Diagnostic{}
		for _, section := range parser.Ast.Sections {
			c.parsingDiagnostics[documentUri] = append(c.parsingDiagnostics[documentUri], section.Diagnostics...)
		}

		c.trees[documentUri] = AstToWordTree(parser.Ast)
	})
	return c
}

func (c *Forest) Remove(documentUri string) {
	c.pool.Run(documentUri, func() {
		c.trees[documentUri] = nil
	})
}

type HarvestedDiagnostic struct {
	Diagnostic lsproto.Diagnostic
	Word       string
	Lang       parser.Language
}

// Based on the built tree, compile tree into diagnostics.
func (c *Forest) Harvest() map[string][]HarvestedDiagnostic {
	c.pool.WaitAll()
	c.harvestMutex.Lock()
	defer c.harvestMutex.Unlock()

	mergedTree := NewWordTree()
	for _, tree := range c.trees {
		mergedTree.Graft(tree)
	}
	fruits := mergedTree.Harvest()

	diags := make(map[string][]HarvestedDiagnostic)
	for uri := range c.trees {
		diags[uri] = []HarvestedDiagnostic{}
	}

	addDiagToAllWordPositions := func(timeRemaining float64, severitiy lsproto.DiagnosticsSeverity, fruit *WordFruit) {
		for _, word := range fruit.Words {
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

			diags[word.Uri()] = append(diags[word.Uri()], HarvestedDiagnostic{
				Lang:       fruit.Lang,
				Diagnostic: *err,
				Word:       fruit.Text,
			})
		}
	}

	for _, fruit := range fruits {
		severity, remainingDays := func() (lsproto.DiagnosticsSeverity, float64) {
			remainingDays := fruitToRemainingDays(fruit)

			if remainingDays <= 1 {
				return lsproto.DiagnosticsSeverityError, remainingDays
			} else if remainingDays < 3 {
				return lsproto.DiagnosticsSeverityHint, remainingDays
			}

			return lsproto.DiagnosticsSeverityInformation, remainingDays
		}()

		addDiagToAllWordPositions(remainingDays, severity, fruit)
	}

	for uri := range c.parsingDiagnostics {
		for _, diag := range c.parsingDiagnostics[uri] {
			diags[uri] = append(diags[uri], HarvestedDiagnostic{
				Diagnostic: *diag,
				Word:       "",
			})
		}
	}

	return diags
}

func (f *Forest) GetTreesLocations() []string {
	return slices.Collect(maps.Keys(f.trees))
}

// Pick a fruit based on its location in the tree and return its remaining days description
func (f *Forest) Pick(textDocument string, line int, character int) (string, bool) {
	if tree, exists := f.trees[textDocument]; exists {
		picked := tree.Pick(line, character)
		if picked == nil {
			return "", false
		}
		remaining := fruitToRemainingDays(picked)
		return fmt.Sprintf("Remaining days: %f", remaining), true
	}
	return "", false
}

func fruitToRemainingDays(fruit *WordFruit) float64 {
	if fruit == nil {
		panic("Fruit is null here...what?!")
	}
	interval := math.Ceil(fruit.Interval)
	deadline := fruit.LastSeenDate.AddDate(0, 0, int(interval))
	remainingHours := time.Until(deadline).Hours()
	var remainingDays float64 = remainingHours / 24
	return remainingDays
}
