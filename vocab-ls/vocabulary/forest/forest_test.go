package forest

import (
	"fmt"
	"testing"
	"time"
	"vocab/syntax"
	test "vocab/vocab_testing"
)

func TestShouldCompile(t *testing.T) {
	text := test.TrimLines(`
		20/05/2025
		> (it) something
		slfjalfjd 
		lkjadsflajsdf
		lksjdf
	`)

	forest := NewForest(t.Context(), func(any) {})
	forest.Plant("xxx", text, nil)
	forest.Harvest()
}

func TestShouldActuallyEmitError(t *testing.T) {
	// parser error
	forst := NewForest(t.Context(), func(a any) {})
	parsingDiag := forst.Plant("xxx", "> (it) la magia, bene,scorprire", nil).Harvest()
	test.Expect(t, true, len(parsingDiag) > 0)

	// compiler error
	forst = NewForest(t.Context(), func(a any) {})
	compilationDiag := forst.Plant("xxx", test.TrimLines(fmt.Sprintf(`
		12/10/1000
		> (it) mostrare(0), %stante%s
	`, "`", "`")), nil).Harvest()
	test.Expect(t, true, len(compilationDiag["xxx"]) > 0)
	diag1 := compilationDiag["xxx"][0]
	test.Expect(t, 1, diag1.Range.Start.Line, diag1.Range.End.Line)
	test.Expect(t, 7, diag1.Range.Start.Character)
	test.Expect(t, 15, diag1.Range.End.Character)
	diag2 := compilationDiag["xxx"][1]
	test.Expect(t, 1, diag2.Range.Start.Line, diag2.Range.End.Line)
	test.Expect(t, 20, diag2.Range.Start.Character)
	test.Expect(t, 27, diag2.Range.End.Character)
}

// Very unlikely in real life, but just for programmatic correctness.
func TestShouldEmitErrorIfReviewedWordIsOfDifferentLanguageBranch(t *testing.T) {
	forest := NewForest(t.Context(), func(a any) {})
	okText := fmt.Sprintf(`
	    20/01/2025
		> (it) sport
		%s
		>> (de) Sport
	`, time.Now().Format(syntax.DateLayout))

	// act
	finalDiag := forest.Plant("doc-1", test.TrimLines(okText), nil).Harvest()

	test.Expect(t, 1, len(finalDiag["doc-1"]))
	test.Expect(t, 1, finalDiag["doc-1"][0].Range.Start.Line, finalDiag["doc-1"][0].Range.End.Line)
}

func TestShouldClearOldParsingDiagnosticsOfCorrectDocument_OnceErrorIsFixed(t *testing.T) {
	// parser error
	forest := NewForest(t.Context(), func(a any) {})
	parsingDiag1 := forest.Plant("doc-1", "> (it) la magia, bene,scorprire", nil).Harvest()
	test.Expect(t, true, len(parsingDiag1["doc-1"]) > 0)

	parsingDiag2 := forest.Plant("doc-2", "> (it) la magia, bene,scorprire", nil).Harvest()
	test.Expect(t, true, len(parsingDiag2["doc-1"]) > 0)

	// act: clear errors from doc-1
	today := time.Now().Format(syntax.DateLayout)
	okText := fmt.Sprintf(`
		%s
		> (it) mostrare
		Mostrare
	`, today)
	finalDiag := forest.Plant("doc-1", test.TrimLines(okText), nil).Harvest()
	test.Expect(t, true, len(finalDiag["doc-1"]) == len(parsingDiag2["doc-2"])-len(parsingDiag1["doc-1"]))
}

func TestShouldReAddParsingDiagnosticsOfCorrectDocument_OnceErrorIsBack(t *testing.T) {
	forest := NewForest(t.Context(), func(a any) {})
	input1 := test.TrimLines(`
		01/01/2025
		> (it) la magia
	`)
	diags := forest.Plant("doc-1", input1, nil).Harvest()
	test.Expect(t, true, len(diags["doc-1"]) > 0)

	// clear errors
	input2 := test.TrimLines(fmt.Sprintf(`
		%s
		> (it) la magia
	`, time.Now().Format(syntax.DateLayout)))
	diags = forest.Plant("doc-1", input2, nil).Harvest()
	test.Expect(t, true, len(diags["doc-1"]) == 0)

	// errors should be back here
	diags = forest.Plant("doc-1", input1, nil).Harvest()
	test.Expect(t, true, len(diags["doc-1"]) > 0)
}

func TestHarvestShouldReturnKeyOfDocumentEvenIfThereAreNoErrors(t *testing.T) {
	forest := NewForest(t.Context(), func(any) {})
	okText := test.TrimLines(fmt.Sprintf(`
		%s
		> (it) la magia
	`, time.Now().Format(syntax.DateLayout)))
	harvested := forest.Plant("xxx", okText, nil).Harvest()
	test.Expect(t, true, harvested["xxx"] != nil)
	test.Expect(t, 0, len(harvested["xxx"]))
}

func TestShouldAllowIncrementalCompilation(t *testing.T) {
	forest := NewForest(t.Context(), func(any) {})

	forest.Plant("xxx", test.TrimLines(`
		12/10/2025
		> (it) mostrare
		>> (it) spiegare
		E oggi voglio spiegarvi un po' com'è essere cittadino a Roma. Che cosa fanno i cittadini a Roma perchè vivono qui. Dov'è vivono la com'è la città e soprattutto com'è la parte che non mostriamo ai turisti.
	`), nil)
	forest.Plant("xxx", test.TrimLines(`
		13/10/2025
		> (it) qualcuno(1), migliaia, decimi(1)
		Qualcuno di voi ha chiesto:
		Finalmente siamo arrivati, dopo migliaia di chilometri e decimi di ore di macchini.
		Nach Tausenden von Kilometern und Zehntelstunden Fahrt sind wir endlich angekommen.
		14/10/2025
		> (de) der Nebensatz, der Relativsatz(3), der Einschub, die Schnodderigkeit(4)
		Der Sprecher benutzt lange, zusammengesetzte Sätze mit Nebensätzen, Relativsätzen und erklärenden Einschüben.
		Er ist voll mit Schnodderigkeit. Kann nicht mit ihm arbeiten...
	`), nil)
	forest.Harvest()

	forest.Plant("xxx", test.TrimLines(`
		16/10/2025
		> (it) %sspiegarmi%s(2), %scom'è%s, eterno, il mito, sfatare, la verità, vera
		> (de) zertreuern, entlarven
		>> (de) gewöhnlich, ewig
		Puoi spiegarmi questa frase pezzo per pezzo?
		Nicht um den Mythos der Ewigen Stadt zu zerstreuen(oder entlarven), das magische Rom, das Sie in Filmen sehen, insbesondere in denen von Fellini, sondern um Ihnen die Wahrheit oder zumindest die wahre Realität von Rom und den Römern zu zeigen.
		Non per sfatare il mito della città eterna, della Roma magica che vedete anche nei film, soprattutto nei film di Fellini, ma per mostrarvi la verità o comunque la vera realtà di Roma e dei romani. 
	`), nil)
	forest.Harvest()
}

func TestEmitErrorAtWordWithSpecialCharacter(t *testing.T) {
	forest := NewForest(t.Context(), func(a any) {})

	// act
	forest.Plant("xxx", "16/10/2025 \n> (it) `com'è`, risolvere", nil)

	harvested := forest.Harvest()
	test.Expect(t, 2, len(harvested["xxx"]))
	diag1 := harvested["xxx"][0]
	test.Expect(t, 1, diag1.Range.Start.Line, diag1.Range.End.Line)
	test.Expect(t, 7, diag1.Range.Start.Character)
	test.Expect(t, 14, diag1.Range.End.Character)

	diag2 := harvested["xxx"][1]
	test.Expect(t, 1, diag2.Range.Start.Line, diag2.Range.End.Line)
	test.Expect(t, 16, diag2.Range.Start.Character)
	test.Expect(t, 25, diag2.Range.End.Character)
}
