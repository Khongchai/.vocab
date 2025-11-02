package languages

import "testing"

func TestStripGermanArticleFromWord(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"der Hund", "Hund"},
		{"die Katze", "Katze"},
		{"das Auto", "Auto"},
		{"ein Mann", "Mann"},
		{"einem Kind", "Kind"},
		{"eine Frau", "Frau"},
		{"ohneArtikel", "ohneArtikel"}, // no article
		{"den Baum", "Baum"},
		{"dem Haus", "Haus"},
		{"des Mannes", "Mannes"},
		// Edge case: double spaces
		{"der  Hund", " Hund"},
	}

	for _, tt := range tests {
		got := StripGermanArticleFromWord(tt.in)
		if got != tt.want {
			t.Errorf("StripGermanArticleFromWord(%q) = %q; want %q", tt.in, got, tt.want)
		}
	}
}

func TestStripItalianArticleFromWord(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"il gatto", "gatto"},
		{"la casa", "casa"},
		{"lo studente", "studente"},
		{"i libri", "libri"},
		{"gli amici", "amici"},
		{"le case", "case"},
		{"una ragazza", "ragazza"},
		{"uno studente", "studente"},
		{"un'amica", "amica"},
		{"l'acqua", "acqua"},
		{"pizza", "pizza"}, // no article
		// Edge case: double space
		{"la  casa", " casa"},
	}

	for _, tt := range tests {
		got := StripItalianArticleFromWord(tt.in)
		if got != tt.want {
			t.Errorf("StripItalianArticleFromWord(%q) = %q; want %q", tt.in, got, tt.want)
		}
	}
}
