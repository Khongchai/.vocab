package vocabulary

import (
	"vocab/lib"
	"vocab/syntax"
)

type Scanner struct {
	text string
	pos  int
	line int
}

func NewScanner(text string) *Scanner {
	return &Scanner{
		text: text,
		pos:  0,
		line: 0,
	}
}

func (s *Scanner) Scan() (Token, string) {
	if s.atEnd() {
		return TokenEOF, ""
	}

	if lib.IsLineBreak(s.char()) {
		return TokenLineBreak, ""
	}

	// Date literal
	if lib.IsDigit(s.char()) {
		collected := ""
		for i := range syntax.DateLength {
			collected += string(s.char())

			// invalid cases
			if (i == syntax.DateSlashFirstPosition || i == syntax.DateSlashSecondPostiion) && s.char() != lib.Slash {
				return TokenText, collected
			} else if !lib.IsDigit(s.char()) {
				return TokenText, collected
			}

			s.pos++
			if s.atEnd() {
				break
			}
		}

		s.pos += syntax.DateLength
		return TokenDateLiteral, collected
	}

	// try parse language ident
	if lib.LeftParen == s.char() {
		collected := ""

		for range syntax.LanguageIdentifierLength - 2 { // minus left and right parent
			collected += string(s.char())

			if !lib.IsASCIILetter(s.char()) {
				return TokenText, collected
			}

			s.pos++
			if s.atEnd() {
				break
			}
		}

		if s.char() == lib.RightParen {
			s.pos++ // If it's not right parent then we don't care / don't read and therefore not forwarding it.
			return TokenLanguageIdent, collected
		} else {
			return TokenText, collected
		}

	}

	if lib.LessThan == s.char() {
		// try parse <!-- end section -->
		// comment
		return TokenText, ""
	}

	if lib.GreaterThan == s.char() {
		// try parse new vocab and reviewed vocab section
		return TokenText, ""
	}

	return TokenText, string(s.char())

}

func (s *Scanner) atEnd() bool {
	return s.pos >= len(s.text)
}

func (s *Scanner) char() rune {
	return rune(s.text[s.pos])
}
