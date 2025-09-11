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
)

var TextToToken = map[string]Token{
	"<":  TokenLessThan,
	">":  TokenGreaterThan,
	">>": TokenDoubleGreaterThan,
	"`":  TokenBacktick,
	"/":  TokenSlash,
	",":  TokenComma,
	")":  TokenOpenBracket,
	"(":  TokenCloseBracket,
	" ":  TokenSpace,
	"-":  TokenMinus,
}
