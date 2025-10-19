package parser

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
	pos int
	// `pos` at which the current token ends
	tokenLineOffsetEnd int
	// `pos` at which the current token begins
	tokenLineOffsetStart int
	line                 int
}

func NewScanner(text string) *Scanner {
	return &Scanner{
		text:                 text,
		pos:                  0,
		tokenLineOffsetEnd:   0,
		tokenLineOffsetStart: 0,
		line:                 0,
	}
}

func (s *Scanner) Scan() (Token, string) {
	scanned, scannedSize := s.charAt(0)

	s.tokenLineOffsetStart = s.tokenLineOffsetEnd

	if lib.IsWhiteSpaceSingleLine(scanned) {
		s.forwardPos(scannedSize)
		return TokenWhitespace, " "
	}

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
			if lib.IsLineBreak(c) || s.atEnd() {
				return TokenText, collected
			}

			if !lib.IsDigit(c) {
				if lib.Slash != c && (i != syntax.DateSlashFirstPosition && i != syntax.DateSlashSecondPostion) {
					return TokenText, collected
				}
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

	case lib.VerticalLine:
		for {
			s.forwardPos(1)
			cur, _ := s.charAt(0)
			if lib.IsLineBreak(cur) || s.atEnd() {
				return TokenCommentTrivia, ""
			}
		}

	case lib.GreaterThan:
		next, _ := s.charAt(1)
		if lib.GreaterThan == next {
			s.forwardPos(2)
			return TokenDoubleGreaterThan, ">>"
		}
		s.forwardPos(1)
		return TokenGreaterThan, ">"

	case lib.Backtick:
		nextChar, nextCharLength := s.charAt(1)
		collected := string(nextChar)
		s.forwardPos(nextCharLength)

		for {
			next, nextLength := s.charAt(1)

			if next == lib.Backtick {
				// if it's a backtick we need to also chomp it too, so forward by 2, the
				// last character in the literal and the backtick itself
				s.forwardPos(2)
				return TokenWordLiteral, collected
			}
			if next == -1 || lib.IsLineBreak(next) {
				s.forwardPos(1)
				return TokenWordLiteral, collected
			}

			collected += string(next)
			s.forwardPos(nextLength)
		}

	case lib.LeftParen:
		collected := "("
		s.forwardPos(1)

		for {
			thisChar, thisCharLen := s.charAt(0)
			switch thisChar {
			case lib.RightParen:
				s.forwardPos(1)
				return TokenSemanticSpecifierLiteral, collected[1:]
			default:
				if s.atEnd() || lib.IsLineBreak(thisChar) {
					return TokenText, collected
				}
				collected += string(thisChar)
				s.forwardPos(thisCharLen)
			}
		}

	default:
		s.forwardPos(scannedSize)
		return TokenText, string(scanned)
	}

}

func (s *Scanner) atEnd() bool {
	return s.pos >= len(s.text)
}
func (s *Scanner) PeekNext() (rune, int) {
	return s.charAt(1)
}

func (s *Scanner) forwardPos(by int) {
	s.pos += by
	s.tokenLineOffsetEnd += by
}

func (s *Scanner) forwardLine() {
	s.pos++
	s.line++
	s.tokenLineOffsetEnd = 0
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
