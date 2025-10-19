package parser

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
	TokenSpace
	TokenMinus
	TokenLessThan
	TokenLeftParen
	TokenRightParen

	TokenCommentTrivia

	TokenDateExpression           // xx/xx/xxxx
	TokenWordLiteral              // `literally` `a` `word`
	TokenSemanticSpecifierLiteral // (de), (it), (4), (5)

	TokenText // all text that do not match anything above
	// We need to emit whitespace here. We can't be 100% sure yet in the scanner whether we're in the
	// example section or not. The parser knows, so all spaces need to be forwarded.
	TokenWhitespace
)
