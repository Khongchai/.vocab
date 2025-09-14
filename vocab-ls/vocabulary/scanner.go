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
	for {
		if !lib.IsWhiteSpaceSingleLine(s.charAt(0)) {
			break
		}
		s.forwardPos(len(string(s.charAt(0))))
	}

	scanned := s.charAt(0)

	if s.atEnd() {
		return TokenEOF, ""
	}

	if lib.IsASCIILetter(scanned) {
		collected := string(s.charAt(0))
		s.forwardPos(1)
		for {
			cur := s.charAt(0)
			if !lib.IsASCIILetter(cur) {
				break
			}
			collected += string(cur)
			s.forwardPos(1)
		}
		return TokenTextLiteral, collected
	}

	if lib.IsLineBreak(scanned) {
		s.forwardLine()
		return TokenLineBreak, string(scanned)
	}

	if lib.IsDigit(scanned) {
		collected := string(s.charAt(0))
		s.forwardPos(1)
		for {
			cur := s.charAt(0)
			if !lib.IsDigit(cur) {
				break
			}
			collected += string(cur)
			s.forwardPos(1)
		}
		return TokenNumericLiteral, collected
	}

	// It should not be possible to overscan!
	switch scanned {
	case lib.Slash:
		s.forwardPos(1)
		return TokenSlash, "/"

	case lib.Comma:
		s.forwardPos(1)
		return TokenComma, ","

	case lib.Minus:
		if s.charAt(0) == lib.Minus && s.charAt(1) == lib.Minus && s.charAt(2) == lib.GreaterThan {
			s.forwardPos(3)
			return TokenMarkdownCommentEnd, "-->"
		}
		s.forwardPos(1)
		return TokenMinus, "-"

	case lib.GreaterThan:
		if lib.GreaterThan == s.charAt(1) {
			s.forwardPos(2)
			return TokenDoubleGreaterThan, ">>"
		}
		s.forwardPos(1)
		return TokenGreaterThan, ">"

	case lib.LessThan:
		if s.charAt(1) == lib.ExclamationMark && s.charAt(2) == lib.Minus && s.charAt(3) == lib.Minus {
			s.forwardPos(4)
			return TokenMarkdownCommentStart, "<!--"
		}
		s.forwardPos(1)
		return TokenLessThan, "<"

	case lib.Backtick:
		if s.charAt(1) == lib.Backtick && s.charAt(2) == lib.Backtick {
			s.forwardPos(3)
			return TokenMarkdownCodefence, "```"
		}
		s.forwardPos(1)
		return TokenBacktick, "`"

	case lib.LeftParen:
		s.forwardPos(1)
		return TokenLeftParen, "("

	case lib.RightParen:
		s.forwardPos(1)
		return TokenRightParen, ")"

	default:

		// all ignored characters here
		s.forwardPos(1)
		return TokenIgnored, string(scanned)
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

func GetTokenStartPos(s *Scanner, tokenText string) int {
	return s.pos - len(tokenText)
}
