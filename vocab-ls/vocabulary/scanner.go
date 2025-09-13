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
		s.line++
		return TokenLineBreak, ""
	}

	// Date literal
	if lib.IsDigit(s.char()) {
		collected := string(s.consume())
		for i := 1; i < syntax.DateLength; i++ {
			if (i == syntax.DateSlashFirstPosition || i == syntax.DateSlashSecondPostiion) && s.char() != lib.Slash {
				return TokenText, collected
			} else if !lib.IsDigit(s.char()) {
				return TokenText, collected
			}
			collected += string(s.consume())
		}

		return TokenDateLiteral, collected
	}

	// try parse language ident
	if lib.LeftParen == s.char() {
		collected := string(s.consume())

		for i := 1; i < syntax.LanguageIdentifierLength; i++ {
			if i == syntax.LanguageIdentifierRightParenPos && s.char() != lib.RightParen {
				return TokenText, collected
			} else if !lib.IsASCIILetter(s.char()) {
				return TokenText, collected
			}
			collected += string(s.consume())
		}

		return TokenLanguageIdent, collected
	}

	if lib.LessThan == s.char() {
		// try parse <!-- end section -->
		// comment

	}

	if lib.GreaterThan == s.char() {
		// try parse new vocab and reviewed vocab section
	}

	// try parse markdown comment or word literal
	if lib.Backtick == s.char() {
		if s.peek(1) == lib.Backtick && s.peek(2) == lib.Backtick {
			tripleTicks := "```"
			s.pos += len(tripleTicks)
			return TokenMarkdownComment, tripleTicks
		}

		collected := string(s.consume())
		for {
			if s.atEnd() {
				return TokenText, collected
			}
			if lib.IsLineBreak(s.char()) {
				return TokenText, collected
			}
			if s.char() == lib.Backtick {
				return TokenWordLiteral, collected
			}
			collected = string(s.consume())
		}
	}

	return TokenText, string(s.char())

}

func (s *Scanner) atEnd() bool {
	return s.pos >= len(s.text)
}

func (s *Scanner) char() rune {
	return rune(s.text[s.pos])
}

// get rune at s.pos and forward s.pos by 1
func (s *Scanner) consume() rune {
	r := rune(s.text[s.pos])
	s.pos++
	return r
}

func (s *Scanner) peek(offset int) rune {
	if s.pos+offset > len(s.text)-1 {
		return -1
	}
	r := rune(s.text[s.pos+offset])
	return r
}
