package vocabulary

type Token int

const (
	TokenGreaterThan Token = iota
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

	TokenMarkdownCommentStart // <!--
	TokenMarkdownCommentEnd   // -->

	TokenDateExpression     // xx/xx/xxxx
	TokenWordExpression     // `literally` `a` `word`
	TokenLanguageExpression // (de) or (it)

	TokenText // all text that do not match anything above
)
