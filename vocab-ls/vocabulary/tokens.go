package vocabulary

type Token int

const (
	TokenGreaterThan Token = iota
	TokenDoubleGreaterThan
	TokenBacktick
	TokenSlash
	TokenComma
	TokenOpenBracket
	TokenCloseBracket
	TokenSpace
	TokenMinus
	TokenLessThan

	TokenItalianKeyword
	TokenGermanKeyword

	TokenNumericLiteral
	TokenWordLiteral // `literally` `a` `word`
)
