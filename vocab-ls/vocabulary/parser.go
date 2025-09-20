package vocabulary

import (
	"context"
	"time"
	lsproto "vocab/lsp"
	"vocab/syntax"
)

type ParsingError string

const (
	MalformedDate            ParsingError = "Malformed date"
	ExpectVocabulary         ParsingError = "Expect Vocabulary"
	ExpectLanguageExpression ParsingError = "The language of this section is not specified. Specified either (it) or (de)"
	UnrecognizedLanguage     ParsingError = "Unrecognized language identifier. Specify either (it) or (de)"
)

type Parser struct {
	ctx     context.Context
	uri     string
	scanner *Scanner
	ast     *VocabAst

	token      Token
	text       string
	tokenStart int // start pos on line
	tokenEnd   int // end pos on line
	line       int // line, 0-indexed

	writeCallback func(any)
}

func NewParser(ctx context.Context, uri string, scanner *Scanner, writeCallback func(any)) *Parser {
	return &Parser{
		ctx:     ctx,
		uri:     uri,
		scanner: scanner,
		ast:     &VocabAst{},

		token:      TokenUnknown,
		text:       "",
		tokenStart: -1,
		tokenEnd:   -1,

		writeCallback: writeCallback,
	}
}

func (p *Parser) Parse() {
	p.ast.Sections = []*VocabularySection{}

	for {
		p.nextToken()

		if p.token == TokenEOF {
			return
		}

		switch p.token {
		case TokenDateExpression:
			p.parseDateExpression()
		case TokenGreaterThan:
			fallthrough
		case TokenDoubleGreaterThan:
			p.parseVocabSection()
		default:
			p.parseSentenceSection()
		}
	}
}

func (p *Parser) parseDateExpression() {
	p.ast.Sections = append(p.ast.Sections, &VocabularySection{})

	text := p.text

	parsed, err := time.Parse(syntax.DateLayout, text)
	if err != nil {
		p.errorHere(&err, "Invalid date format")
	}

	date := &DateSection{Text: text, Time: parsed, Start: p.tokenStart, End: p.tokenEnd}
	p.currentVocabSection().Date = date
}

func (p *Parser) parseVocabSection() {
	currentSection := p.currentVocabSection()
	words := &WordsSection{
		Line: p.line,
	}
	if p.token == TokenGreaterThan {
		currentSection.NewWords = append(currentSection.NewWords, words)
	} else {
		currentSection.ReviewedWords = append(currentSection.NewWords, words)
	}

	p.nextToken()
	if p.token != TokenLanguageExpression {
		p.errorHere(nil, ExpectLanguageExpression)
		return
	}
	if p.text == "(it)" {
		words.Language = Italiano
	} else if p.text == "(de)" {
		words.Language = Deutsch
	} else {
		p.errorHere(nil, UnrecognizedLanguage)
		words.Language = Unrecognized
	}

	parsing := ""

	for {
		switch p.token {
		case TokenText:
			parsing += p.text
		case TokenComma:
			parsing = ""
		case TokenEOF:
			return
		case TokenLineBreak:
			return
		default:
			p.errorHere(nil, ExpectVocabulary)
		}
	}
}

func (p *Parser) parseSentenceSection() {

}

// Add a diagnostics error to this line.
func (p *Parser) errorHere(original *error, message ParsingError) {
	if original != nil {
		errorMessage := (*original).Error()
		p.writeCallback(errorMessage)
	}
	diag := p.currentVocabSection().Diagnostics
	newError := &lsproto.Diagnostic{
		Severity: lsproto.DiagnosticsSeverityError,
		Message:  string(message),
		Range: lsproto.Range{
			Start: lsproto.Position{
				Line:      p.line,
				Character: p.tokenStart,
			},
			End: lsproto.Position{
				Line:      p.line,
				Character: p.tokenEnd,
			},
		},
	}
	diag = append(diag, newError)
}

func (p *Parser) currentVocabSection() *VocabularySection {
	return p.ast.Sections[len(p.ast.Sections)-1]
}

func (p *Parser) nextToken() {
	p.line = p.scanner.line
	token, text := p.scanner.Scan()
	p.text = text
	p.token = token
	p.tokenEnd = p.scanner.lineOffset
	p.tokenStart = p.tokenEnd - len(text)
}
