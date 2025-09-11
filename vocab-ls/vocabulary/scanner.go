package vocabulary

type Scanner struct {
	CurrentToken *Token
}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (*Scanner) NextToken() *Token {
	return nil
}
