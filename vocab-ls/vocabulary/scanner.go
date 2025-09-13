package vocabulary

type Scanner struct {
	currentToken *Token
	text         string
	position     int
	line         int
}

func NewScanner(text string) *Scanner {
	return &Scanner{
		text:         text,
		position:     0,
		line:         0,
		currentToken: nil,
	}
}

func (*Scanner) NextToken() *Token {
	return nil
}

func (s *Scanner) CurrentToken() *Token {
	return s.currentToken
}
