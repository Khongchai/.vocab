package vocabulary

type Token int

const (
	TokenGreaterThan Token = iota
	TokenDoubleGreaterThan
	TokenBacktick
	TokenSlash
	TokenComma
)

var TextToToken = map[string]Token{
	">":  TokenGreaterThan,
	">>": TokenDoubleGreaterThan,
	"`":  TokenBacktick,
	"/":  TokenSlash,
	",":  TokenComma,
}
