package vocabulary

import (
	"context"
	"strings"
	"time"
	lsproto "vocab/lsp"
	"vocab/syntax"
)

const (
	MalformedDate            string = "Malformed date"
	ExpectDate               string = "Expect Date"
	ExpectVocabulary         string = "Expect Vocabulary"
	ExpectLanguageExpression string = "The language of this section is not specified. Specified either (it) or (de)"
	UnrecognizedLanguage     string = "Unrecognized language identifier. Specify either (it) or (de)"
	ExpectVocabSection       string = "Expect Vocab Section"
	ExpectDateSection        string = "Expect Date Section"
	UnexpectedToken          string = "Unexpected Token"
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

	printCallback func(any)
}

func NewParser(ctx context.Context, uri string, scanner *Scanner, printCallback func(any)) *Parser {
	return &Parser{
		ctx:     ctx,
		uri:     uri,
		scanner: scanner,
		ast:     &VocabAst{},

		token:      TokenUnknown,
		text:       "",
		tokenStart: -1,
		tokenEnd:   -1,

		printCallback: printCallback,
	}
}

func (p *Parser) Parse() {
	p.ast.Sections = []*VocabularySection{}

	var lastSection *VocabularySection = nil

	for {
		p.nextToken()

		if p.token == TokenEOF {
			return
		}

		if len(p.ast.Sections) > 0 {
			lastSection = p.ast.Sections[len(p.ast.Sections)-1]
		}

		switch p.token {
		case TokenLineBreakTrivia, TokenWhitespace:
			continue
		case TokenDateExpression:
			p.ast.Sections = append(p.ast.Sections, &VocabularySection{})
			p.parseDateExpression()
		case TokenGreaterThan, TokenDoubleGreaterThan:
			if lastSection.Date == nil {
				p.errorHere(nil, ExpectDateSection)
			}
			p.parseVocabSection()
		default:
			if lastSection == nil || len(lastSection.NewWords) == 0 && len(lastSection.ReviewedWords) == 0 {
				p.ast.Sections = append(p.ast.Sections, &VocabularySection{})
				p.errorHere(nil, ExpectVocabSection)
				return
			}
			p.parseUtteranceSection()
		}
	}
}

func (p *Parser) parseDateExpression() {
	parsed, err := time.Parse(syntax.DateLayout, p.text)
	parsedAsLocalTime := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.Local)
	date := &DateSection{Text: p.text, Time: parsedAsLocalTime, Start: p.tokenStart, End: p.tokenEnd, Line: p.line}
	p.currentVocabSection().Date = date

	if err != nil {
		p.errorHere(&err, MalformedDate)
		return
	}

	p.nextToken()

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

	p.nextTokenNotWhitespace()

	if p.token != TokenLanguageExpression {
		p.errorHere(nil, ExpectLanguageExpression)
		return
	}
	switch p.text {
	case "(it)":
		words.Language = Italiano
	case "(de)":
		words.Language = Deutsch
	default:
		p.errorHere(nil, UnrecognizedLanguage)
		words.Language = Unrecognized
	}

	parsing := ""

	p.nextTokenNotWhitespace()

outer:
	for {
		switch p.token {
		case TokenEOF, TokenLineBreakTrivia, TokenComma:
			newWord := &Word{Text: parsing, Start: p.tokenStart, End: p.tokenEnd}
			words.Words = append(words.Words, newWord)

			if p.token == TokenLineBreakTrivia || p.token == TokenEOF {
				break outer
			}

			parsing = ""
			p.nextTokenNotWhitespace()
		default:
			parsing += p.text
			p.nextToken()
		}
	}
}

func (p *Parser) parseUtteranceSection() {
	var sb strings.Builder

	for {
		switch p.token {
		case TokenEOF:
			return
		case TokenLineBreakTrivia:
			newUtterance := &UtteranceSection{}
			newUtterance.Text = sb.String()
			p.currentVocabSection().Utterance = append(p.currentVocabSection().Utterance, newUtterance)
			p.nextToken()
			return
		default:
			sb.WriteString(p.text)
		}
	}
}

// Add a diagnostics error to this line and forward until new line.
func (p *Parser) errorHere(original *error, message string) {
	if original != nil {
		p.printCallback(&original)
	}
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
	diag := p.currentVocabSection().Diagnostics
	p.currentVocabSection().Diagnostics = append(diag, newError)

	for {
		p.nextToken()

		if p.token == TokenLineBreakTrivia || p.token == TokenEOF {
			break
		}
	}
}

func (p *Parser) nextTokenNotWhitespace() {
	p.nextToken()
	for {
		if p.token != TokenWhitespace {
			return
		}
		p.nextToken()
	}
}

func (p *Parser) currentVocabSection() *VocabularySection {
	sectionCount := len(p.ast.Sections)
	if sectionCount == 0 {
		return nil
	}
	last := p.ast.Sections[sectionCount-1]
	return last
}

func (p *Parser) nextToken() {
	p.line = p.scanner.line
	token, text := p.scanner.Scan()
	p.text = text
	p.token = token
	p.tokenEnd = p.scanner.lineOffset
	p.tokenStart = p.tokenEnd - len(text)
}
