package compiler

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

	compiler := NewCompiler(t.Context(), func(any) {})
	compiler.Accept("xxx", text, nil)
	compiler.Compile()
}

func TestShouldActuallyEmitError(t *testing.T) {
	// parser error
	compiler := NewCompiler(t.Context(), func(a any) {})
	parsingDiag := compiler.Accept("xxx", "> (it) la magia, bene,scorprire", nil).Compile()
	test.Expect(t, true, len(parsingDiag) > 0)

	// compiler error
	compiler = NewCompiler(t.Context(), func(a any) {})
	compilationDiag := compiler.Accept("xxx", test.TrimLines(`
		12/10/1000
		> (it) mostrare
		Mostrare
	`), nil).Compile()
	test.Expect(t, true, len(compilationDiag) > 0)
}

func TestShouldClearOldParsingDiagnosticsOfCorrectDocument_OnceErrorIsFixed(t *testing.T) {
	// parser error
	compiler := NewCompiler(t.Context(), func(a any) {})
	parsingDiag1 := compiler.Accept("doc-1", "> (it) la magia, bene,scorprire", nil).Compile()
	test.Expect(t, true, len(parsingDiag1) > 0)

	parsingDiag2 := compiler.Accept("doc-2", "> (it) la magia, bene,scorprire", nil).Compile()
	test.Expect(t, true, len(parsingDiag2) > 0)

	// act: clear errors from doc-1
	today := time.Now().Format(syntax.DateLayout)
	okText := fmt.Sprintf(`
		%s
		> (it) mostrare
		Mostrare
	`, today)
	finalDiag := compiler.Accept("doc-1", test.TrimLines(okText), nil).Compile()
	test.Expect(t, true, len(finalDiag) == len(parsingDiag2)-len(parsingDiag1))
}

func TestShouldAllowIncrementalCompilation(t *testing.T) {
	compiler := NewCompiler(t.Context(), func(any) {})

	compiler.Accept("xxx", test.TrimLines(`
		12/10/2025
		> (it) mostrare
		>> (it) spiegare
		E oggi voglio spiegarvi un po' com'è essere cittadino a Roma. Che cosa fanno i cittadini a Roma perchè vivono qui. Dov'è vivono la com'è la città e soprattutto com'è la parte che non mostriamo ai turisti.
	`), nil)
	compiler.Accept("xxx", test.TrimLines(`
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
	compiler.Compile()

	compiler.Accept("xxx", test.TrimLines(`
		16/10/2025
		> (it) %sspiegarmi%s(2), %scom'è%s, eterno, il mito, sfatare, la verità, vera
		> (de) zertreuern, entlarven
		>> (de) gewöhnlich, ewig
		Puoi spiegarmi questa frase pezzo per pezzo?
		Nicht um den Mythos der Ewigen Stadt zu zerstreuen(oder entlarven), das magische Rom, das Sie in Filmen sehen, insbesondere in denen von Fellini, sondern um Ihnen die Wahrheit oder zumindest die wahre Realität von Rom und den Römern zu zeigen.
		Non per sfatare il mito della città eterna, della Roma magica che vedete anche nei film, soprattutto nei film di Fellini, ma per mostrarvi la verità o comunque la vera realtà di Roma e dei romani. 
	`), nil)
	compiler.Compile()
}
