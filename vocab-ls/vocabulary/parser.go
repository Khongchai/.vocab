package vocabulary

import (
	"context"
	"time"
	lsproto "vocab/lsp"
)

type ParsingError string

const (
	MalformedDate ParsingError = "Malformed date"
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
			p.parseVocabSection()
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

	parsed, err := time.Parse("02/01/2006", text)
	if err != nil {
		p.errorHere(err, "Invalid date format")
	}

	date := &DateSection{Text: text, Time: parsed, Start: p.tokenStart, End: p.tokenEnd}
	p.currentVocabSection().Date = date
}

// Add a diagnostics error to this line.
func (p *Parser) errorHere(original error, message ParsingError) {
	p.writeCallback(original)
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
