package vocabulary

type Token int

const (
	TokenUnknown Token = iota
	TokenGreaterThan
	TokenDoubleGreaterThan
	TokenBacktick
	TokenSlash
	TokenComma
	TokenEOF
	TokenLineBreak
	TokenOpenBracket
	TokenCloseBracket
	TokenSpace
	TokenMinus
	TokenLessThan

	TokenLanguageIdent

	TokenNumericLiteral
	TokenWordLiteral // `literally` `a` `word`

	TokenDateLiteral // special syntax for xx/xx/xxxx where x is any digit
	TokenMarkdownComment

	TokenText // all valid markdown text
)
