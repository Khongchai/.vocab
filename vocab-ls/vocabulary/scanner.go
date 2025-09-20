package vocabulary

import (
	"unicode/utf8"
	"vocab/lib"
	"vocab/syntax"
)

// Diagnostics error from scanners are added when multi-tokens identifier fail to match something.
// Most of vocab tokens are pretty simple, so we'll omit the for now.
type Scanner struct {
	text string
	// This always report the next position to be read from.
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
		c, cSize := s.charAt(0)
		if !lib.IsWhiteSpaceSingleLine(c) {
			break
		}
		s.forwardPos(cSize)
	}

	scanned, scannedSize := s.charAt(0)

	if s.atEnd() {
		return TokenEOF, ""
	}

	if lib.IsRecognizedLetter(scanned) {
		c, cSize := s.charAt(0)
		collected := string(c)
		s.forwardPos(cSize)
		for {
			c, cSize = s.charAt(0)
			if !lib.IsRecognizedLetter(c) {
				break
			}
			collected += string(c)
			s.forwardPos(cSize)
		}
		return TokenText, collected
	}

	if lib.IsLineBreak(scanned) {
		s.forwardLine()
		return TokenLineBreak, string(scanned)
	}

	if lib.IsDigit(scanned) {
		collected := string(scanned)
		s.forwardPos(1)
		for i := 1; i < syntax.DateLength; i++ {
			c, _ := s.charAt(0)
			if !lib.IsDigit(c) && i != syntax.DateSlashFirstPosition && i != syntax.DateSlashSecondPostion {
				return TokenText, collected
			}
			collected += string(c)
			s.forwardPos(1)
		}
		return TokenDateExpression, collected
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
		c2, _ := s.charAt(1)
		c3, _ := s.charAt(2)
		if c2 == lib.Minus && c3 == lib.GreaterThan {
			s.forwardPos(3)
			return TokenMarkdownCommentEnd, "-->"
		}
		s.forwardPos(1)
		return TokenMinus, "-"

	case lib.GreaterThan:
		next, _ := s.charAt(1)
		if lib.GreaterThan == next {
			s.forwardPos(2)
			return TokenDoubleGreaterThan, ">>"
		}
		s.forwardPos(1)
		return TokenGreaterThan, ">"

	case lib.LessThan:
		c1, _ := s.charAt(1)
		c2, _ := s.charAt(2)
		c3, _ := s.charAt(3)
		if c1 == lib.ExclamationMark && c2 == lib.Minus && c3 == lib.Minus {
			s.forwardPos(4)
			return TokenMarkdownCommentStart, "<!--"
		}
		s.forwardPos(1)
		return TokenLessThan, "<"

	case lib.Backtick:
		collected := string(scanned)
		s.forwardPos(1)

		for {
			cur, curLength := s.charAt(0)
			collected += string(cur)
			s.forwardPos(curLength)

			if cur == lib.Backtick {
				return TokenWordExpression, collected
			}

			if s.atEnd() || lib.IsLineBreak(cur) {
				return TokenText, collected
			}
		}

	case lib.LeftParen:
		possibleClosing, _ := s.charAt(3)
		if possibleClosing == lib.RightParen {
			expr := s.text[s.pos : 4+s.pos]
			s.forwardPos(4)
			return TokenLanguageExpression, expr
		}

		// To keep things consistent, we consume some number of text up until len("(xx)") before continuing.
		remainingSize := len(s.text[s.pos:])
		consumableSize := min(4, remainingSize)
		consumed := s.text[s.pos:consumableSize]
		s.forwardPos(len(consumed))
		return TokenText, consumed

	default:
		s.forwardPos(scannedSize)
		return TokenText, string(scanned)
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
func (s *Scanner) charAt(offset int) (rune, int) {
	if s.pos+offset > len(s.text)-1 {
		return -1, 0
	}
	r, size := utf8.DecodeRuneInString(s.text[s.pos+offset:])
	return r, size
}

func (s *Scanner) CurrentPosition() int {
	return s.pos - 1
}

func GetTokenStartPos(s *Scanner, tokenText string) int {
	return s.pos - len(tokenText) - 1
}
