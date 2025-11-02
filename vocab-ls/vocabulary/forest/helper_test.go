package forest

import (
	"fmt"
	"testing"
	test "vocab/vocab_testing"
	parser "vocab/vocabulary/parser"
)

func TestVocabAstToWordTree(t *testing.T) {
	text := test.TrimLines(fmt.Sprintf(`
		12/10/2025
		> (it) mostrare
		>> (it) spiegare
		E oggi voglio spiegarvi un po' com'è essere cittadino a Roma. Che cosa fanno i cittadini a Roma perchè vivono qui. Dov'è vivono la com'è la città e soprattutto com'è la parte che non mostriamo ai turisti.
		13/10/2025
		> (it) qualcuno(1), migliaia, decimi(1)
		Qualcuno di voi ha chiesto:
		Finalmente siamo arrivati, dopo migliaia di chilometri e decimi di ore di macchini.
		Nach Tausenden von Kilometern und Zehntelstunden Fahrt sind wir endlich angekommen.
		14/10/2025
		> (de) der Nebensatz, der Relativsatz(3), der Einschub, die Schnodderigkeit(4)
		Der Sprecher benutzt lange, zusammengesetzte Sätze mit Nebensätzen, Relativsätzen und erklärenden Einschüben.
		Er ist voll mit Schnodderigkeit. Kann nicht mit ihm arbeiten...
		15/10/2025
		> (de) ausnahmsweise(1), der Rechner(5), nachhaltig(4)
		>> (de) gewöhnlich
		>> (it) qualcuno(3)
		Ich werd' ausnahmsweise mal nichts machen.
		Das ist so gewöhnlich dass keiner davon im Geringsten die Augenbrauen hochziehen würde.
		Deutschland hatte zu Beginn des Krieges (1939) eine beeindruckend effiziente Kriegswirtschaft, aber sie war nicht nachhaltig.
		16/10/2025
		> (it) %sspiegarmi%s(2), %scom'è%s, eterno, il mito, sfatare, la verità, vera
		> (de) zertreuern, entlarven
		>> (de) gewöhnlich, ewig
		Puoi spiegarmi questa frase pezzo per pezzo?
		Nicht um den Mythos der Ewigen Stadt zu zerstreuen(oder entlarven), das magische Rom, das Sie in Filmen sehen, insbesondere in denen von Fellini, sondern um Ihnen die Wahrheit oder zumindest die wahre Realität von Rom und den Römern zu zeigen.
		Non per sfatare il mito della città eterna, della Roma magica che vedete anche nei film, soprattutto nei film di Fellini, ma per mostrarvi la verità o comunque la vera realtà di Roma e dei romani. 
	`, "`", "`", "`", "`"))

	p := parser.NewParser(t.Context(), "xxx", parser.NewScanner(text), func(a any) {})
	p.Parse()

	for _, section := range p.Ast.Sections {
		test.Expect(t, 0, len(section.Diagnostics)) // hopefully no syntax error!
	}

	tree := AstToWordTree(p.Ast)
	test.Expect(t, 2, len(tree.GetBranches()))

	italianBranch := tree.branches[string(parser.Italiano)].twigs
	test.Expect(t, 12, len(italianBranch))

	germanBranch := tree.branches[string(parser.Deutsch)].twigs
	test.Expect(t, 11, len(germanBranch))

	test.Expect(t, 1, len(italianBranch["mostrare"]))
	test.Expect(t, 1, len(italianBranch["spiegare"]))
	test.Expect(t, 2, len(italianBranch["qualcuno"]))
	test.Expect(t, 1, len(italianBranch["migliaia"]))
	test.Expect(t, 1, len(italianBranch["decimi"]))
	test.Expect(t, 1, len(italianBranch["spiegarmi"]))
	test.Expect(t, 1, len(italianBranch["com'è"]))
	test.Expect(t, 1, len(italianBranch["eterno"]))
	test.Expect(t, 1, len(italianBranch["mito"]))
	test.Expect(t, 1, len(italianBranch["sfatare"]))
	test.Expect(t, 1, len(italianBranch["verità"]))
	test.Expect(t, 1, len(italianBranch["vera"]))

	test.Expect(t, 1, len(germanBranch["Nebensatz"]))
	test.Expect(t, 1, len(germanBranch["Relativsatz"]))
	test.Expect(t, 1, len(germanBranch["Einschub"]))
	test.Expect(t, 1, len(germanBranch["Schnodderigkeit"]))
	test.Expect(t, 1, len(germanBranch["ausnahmsweise"]))
	test.Expect(t, 1, len(germanBranch["Rechner"]))
	test.Expect(t, 1, len(germanBranch["nachhaltig"]))
	test.Expect(t, 1, len(germanBranch["zertreuern"]))
	test.Expect(t, 1, len(germanBranch["entlarven"]))
	test.Expect(t, 2, len(germanBranch["gewöhnlich"]))
	test.Expect(t, 1, len(germanBranch["ewig"]))
}
