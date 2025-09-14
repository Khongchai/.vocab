package vocabulary

type Token int

const (
	TokenIgnored Token = iota
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

	TokenLanguageIdent

	TokenNumericLiteral // same as markdown text but a special case when it's a number for easier detection in parser
	TokenTextLiteral    // all valid markdown text

	TokenMarkdownCommentStart // <!--
	TokenMarkdownCommentEnd   // -->
	TokenMarkdownCodefence    // ```

	TokenDateExpression // xx/xx/xxxx
	TokenWordExpression // `literally` `a` `word`
)
