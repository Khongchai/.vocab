package vocabulary

import (
	"vocab/lib"
)

type Scanner struct {
	text       string
	pos        int
	lineOffset int
	line       int
}

func NewScanner(text string) *Scanner {
	return &Scanner{
		text:       text,
		pos:        0,
		lineOffset: 0,
		line:       0,
	}
}

func (s *Scanner) Scan() (Token, string) {
	if s.atEnd() {
		return TokenEOF, ""
	}

	// It should not be possible to overscan!
	switch scanned := s.charAt(0); scanned {
	case lib.Slash:
		s.forwardPos(1)
		return TokenSlash, "/"

	case lib.Comma:
		s.forwardPos(1)
		return TokenComma, ","

	case lib.GreaterThan:
		if lib.GreaterThan == s.charAt(1) {
			s.forwardPos(2)
			return TokenDoubleGreaterThan, ">>"
		}
		s.forwardPos(1)
		return TokenGreaterThan, ">"

	case lib.Backtick:
		if s.charAt(1) == lib.Backtick && s.charAt(2) == lib.Backtick {
			return TokenMarkdownCodefence, "```"
		}
		return TokenBacktick, "`"

	case lib.LeftParen:
		s.forwardPos(1)
		return TokenLeftParen, "("

	case lib.RightParen:
		s.forwardPos(1)
		return TokenRightParen, ")"

	default:
		if lib.IsASCIILetter(scanned) {
			collected := ""
			for current := s.charAt(0); !lib.IsASCIILetter(current); s.forwardPos(1) {
				collected += string(current)
			}
			return TokenTextLiteral, collected
		}

		if lib.IsLineBreak(scanned) {
			s.forwardLine()
			return TokenLineBreak, ""
		}

		if lib.IsDigit(scanned) {
			collected := ""
			for current := s.charAt(0); !lib.IsDigit(current); s.forwardPos(1) {
				collected += string(current)
			}
			return TokenNumericLiteral, collected
		}

		// all ignored characters here
		s.forwardPos(1)
		return TokenUnknown, string(scanned)
	}

}

func (s *Scanner) atEnd() bool {
	return s.pos >= len(s.text)
}

func (s *Scanner) forwardPos(by int) {
	s.pos += by
	s.lineOffset += by
}

func (s *Scanner) forwardLine() {
	s.pos++
	s.line++
	s.lineOffset = 0
}

// Does not throw error and return -1 if index out of range
func (s *Scanner) charAt(offset int) rune {
	if s.pos+offset > len(s.text)-1 {
		return -1
	}
	r := rune(s.text[s.pos+offset])
	return r
}
