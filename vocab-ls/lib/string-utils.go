package lib

func IsWhiteSpaceSingleLine(ch rune) bool {
	switch ch {
	case
		' ',    // space
		'\t',   // tab
		'\v',   // verticalTab
		'\f',   // formFeed
		0x0085, // nextLine
		0x00A0, // nonBreakingSpace
		0x1680, // ogham
		0x2000, // enQuad
		0x2001, // emQuad
		0x2002, // enSpace
		0x2003, // emSpace
		0x2004, // threePerEmSpace
		0x2005, // fourPerEmSpace
		0x2006, // sixPerEmSpace
		0x2007, // figureSpace
		0x2008, // punctuationEmSpace
		0x2009, // thinSpace
		0x200A, // hairSpace
		0x200B, // zeroWidthSpace
		0x202F, // narrowNoBreakSpace
		0x205F, // mathematicalSpace
		0x3000, // ideographicSpace
		0xFEFF: // byteOrderMark
		return true
	}
	return false
}

const (
	// Whitespace and control characters
	LineFeed       rune = 0x0a   // \n
	CarriageReturn rune = 0x0d   // \r
	Space          rune = 0x0020 // " "
	FormFeed       rune = 0x0c   // \f
	Tab            rune = 0x09   // \t

	// Digits
	Digit0 rune = 0x30
	Digit1 rune = 0x31
	Digit2 rune = 0x32
	Digit3 rune = 0x33
	Digit4 rune = 0x34
	Digit5 rune = 0x35
	Digit6 rune = 0x36
	Digit7 rune = 0x37
	Digit8 rune = 0x38
	Digit9 rune = 0x39

	// Lowercase letters
	LowerA rune = 0x61
	LowerB rune = 0x62
	LowerC rune = 0x63
	LowerD rune = 0x64
	LowerE rune = 0x65
	LowerF rune = 0x66
	LowerG rune = 0x67
	LowerH rune = 0x68
	LowerI rune = 0x69
	LowerJ rune = 0x6a
	LowerK rune = 0x6b
	LowerL rune = 0x6c
	LowerM rune = 0x6d
	LowerN rune = 0x6e
	LowerO rune = 0x6f
	LowerP rune = 0x70
	LowerQ rune = 0x71
	LowerR rune = 0x72
	LowerS rune = 0x73
	LowerT rune = 0x74
	LowerU rune = 0x75
	LowerV rune = 0x76
	LowerW rune = 0x77
	LowerX rune = 0x78
	LowerY rune = 0x79
	LowerZ rune = 0x7a

	// Uppercase letters
	UpperA rune = 0x41
	UpperB rune = 0x42
	UpperC rune = 0x43
	UpperD rune = 0x44
	UpperE rune = 0x45
	UpperF rune = 0x46
	UpperG rune = 0x47
	UpperH rune = 0x48
	UpperI rune = 0x49
	UpperJ rune = 0x4a
	UpperK rune = 0x4b
	UpperL rune = 0x4c
	UpperM rune = 0x4d
	UpperN rune = 0x4e
	UpperO rune = 0x4f
	UpperP rune = 0x50
	UpperQ rune = 0x51
	UpperR rune = 0x52
	UpperS rune = 0x53
	UpperT rune = 0x54
	UpperU rune = 0x55
	UpperV rune = 0x56
	UpperW rune = 0x57
	UpperX rune = 0x58
	UpperY rune = 0x59
	UpperZ rune = 0x5a

	// Symbols
	Asterisk        rune = 0x2a // *
	Backslash       rune = 0x5c // \
	CloseBrace      rune = 0x7d // }
	CloseBracket    rune = 0x5d // ]
	Colon           rune = 0x3a // :
	Comma           rune = 0x2c // ,
	Dash            rune = 0x2d // -
	Dot             rune = 0x2e // .
	DoubleQuote     rune = 0x22 // "
	ExclamationMark rune = 0x21 // !
	GreaterThan     rune = 0x3e // >
	LessThan        rune = 0x3c // <
	Minus           rune = 0x2d // -
	OpenBrace       rune = 0x7b // {
	LeftParen       rune = 0x28 // (
	RightParen      rune = 0x29 // )
	Plus            rune = 0x2b // +
	Slash           rune = 0x2f // /
	Backtick        rune = 0x60 // `
)

func IsLineBreak(ch rune) bool {
	switch ch {
	case
		'\n',   // lineFeed
		'\r',   // carriageReturn
		0x2028, // lineSeparator
		0x2029: // paragraphSeparator
		return true
	}
	return false
}

func IsGermanOrItalianLetter(ch rune) bool {
	// isSpecialChar checks if a rune is a German or Italian special character
	switch ch {
	// German
	case 'Ä', 'ä', 'Ö', 'ö', 'Ü', 'ü', 'ß':
		fallthrough
	case 'À', 'à', 'È', 'è', 'É', 'é', 'Ì', 'ì', 'Ò', 'ò', 'Ù', 'ù':
		return true
	}
	return false
}

func IsDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func IsASCIILetter(ch rune) bool {
	return ch >= 'A' && ch <= 'Z' || ch >= 'a' && ch <= 'z'
}

func IsRecognizedLetter(ch rune) bool {
	return IsASCIILetter(ch) || IsGermanOrItalianLetter(ch)
}
