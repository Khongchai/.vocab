package vocabulary

import (
	"context"
	"strings"
	"time"
	lsproto "vocab/lsp"
	"vocab/syntax"
	entity "vocab/vocabulary/entity"
)

const (
	MalformedDate            string = "Malformed date"
	ExpectDateSection        string = "Expect a date section here."
	ExpectVocabulary         string = "Expect Vocabulary"
	ExpectLanguageExpression string = "The language of this section is not specified. Specified either (it) or (de)"
	UnrecognizedLanguage     string = "Unrecognized language identifier. Specify either (it) or (de)"
	ExpectVocabSection       string = "Expect Vocab Section"
	UnexpectedToken          string = "Unexpected Token"
)

type Parser struct {
	ctx     context.Context
	uri     string
	scanner *Scanner
	ast     *entity.VocabAst

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
		ast:     &entity.VocabAst{},

		token:      TokenUnknown,
		text:       "",
		tokenStart: -1,
		tokenEnd:   -1,

		printCallback: printCallback,
	}
}

func (p *Parser) Parse() {
	p.ast.Sections = []*entity.VocabularySection{}

	var lastSection *entity.VocabularySection = nil

	startNewSection := func() {
		p.ast.Sections = append(p.ast.Sections, entity.NewVocabularySection(p.uri))
	}

	for {
		p.nextToken()

		if p.token == TokenEOF {
			return
		}

		if len(p.ast.Sections) > 0 {
			lastSection = p.ast.Sections[len(p.ast.Sections)-1]
		}

		switch p.token {
		case TokenLineBreak, TokenWhitespace:
			continue
		case TokenDateExpression:
			startNewSection()
			p.parseDateExpression()
		case TokenGreaterThan, TokenDoubleGreaterThan:
			if lastSection == nil || lastSection.Date == nil {
				startNewSection()
				p.errorHere(nil, ExpectDateSection)
			}
			p.parseVocabSection()
		default:
			if lastSection == nil {
				startNewSection()
				p.errorHere(nil, ExpectDateSection)
				continue
			}
			if len(lastSection.NewWords) == 0 && len(lastSection.ReviewedWords) == 0 {
				startNewSection()
				p.errorHere(nil, ExpectVocabSection)
				continue
			}
			p.parseUtteranceSection()
		}
	}
}

func (p *Parser) parseDateExpression() {
	parsed, err := time.Parse(syntax.DateLayout, p.text)
	parsedAsLocalTime := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.Local)
	date := &entity.DateSection{Text: p.text, Time: parsedAsLocalTime, Start: p.tokenStart, End: p.tokenEnd, Line: p.line}
	p.currentVocabSection().Date = date

	if err != nil {
		p.errorHere(&err, MalformedDate)
		return
	}

	p.nextToken()

}

func (p *Parser) parseVocabSection() {
	currentSection := p.currentVocabSection()
	words := &entity.WordsSection{
		Line: p.line,
	}
	if p.token == TokenGreaterThan {
		words.Reviewed = false
		currentSection.NewWords = append(currentSection.NewWords, words)
	} else {
		words.Reviewed = true
		currentSection.ReviewedWords = append(currentSection.ReviewedWords, words)
	}

	p.nextTokenNotWhitespace()

	if p.token != TokenLanguageLiteral {
		p.errorHere(nil, ExpectLanguageExpression)
		return
	}
	switch p.text {
	case "it":
		words.Language = entity.Italiano
	case "de":
		words.Language = entity.Deutsch
	default:
		p.errorHere(nil, UnrecognizedLanguage)
		words.Language = entity.Unrecognized
	}

	parsing := ""

	p.nextTokenNotWhitespace()

	newWordFromText := func(t string) {
		isWordLiteral := p.token == TokenWordLiteral
		text := func() string {
			if isWordLiteral {
				return p.text
			}
			return t
		}()

		newWord := &entity.Word{Text: text, Start: p.tokenStart, End: p.tokenEnd, Literally: isWordLiteral}
		words.Words = append(words.Words, newWord)

	}

	for {
		switch p.token {
		case TokenLineBreak, TokenEOF:
			if parsing != "" {
				newWordFromText(parsing)
			}
			return
		case TokenWordLiteral:
			newWordFromText(parsing)
			p.nextTokenNotWhitespace()
		case TokenComma:
			if parsing == "" {
				p.nextTokenNotWhitespace()
				continue
			}

			newWordFromText(parsing)
			p.nextTokenNotWhitespace()
			parsing = ""
		default:
			parsing += p.text
			p.nextToken()
		}
	}
}

func (p *Parser) parseUtteranceSection() {
	var sb strings.Builder

	start := p.tokenStart

	for {
		switch p.token {
		case TokenLineBreak, TokenEOF:
			text := sb.String()
			newUtterance := &entity.UtteranceSection{
				Line:  p.line,
				Start: start,
				End:   start + len(text),
				Text:  sb.String(),
			}
			p.currentVocabSection().Utterance = append(p.currentVocabSection().Utterance, newUtterance)
			return
		default:
			sb.WriteString(p.text)
			p.nextToken()
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

	// Only 1 error per line for simplicity.
	for {
		p.nextToken()

		if p.token == TokenLineBreak || p.token == TokenEOF {
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

func (p *Parser) currentVocabSection() *entity.VocabularySection {
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
	p.tokenStart = max(p.tokenEnd-len(text), 0)
}
