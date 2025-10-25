package parser

import (
	"context"
	"strconv"
	"strings"
	"time"
	lsproto "vocab/lsp"
	"vocab/syntax"
)

const (
	MalformedDate            string = "Malformed date -- expected dd/mm/yyyy"
	ExpectDateSection        string = "Expect a date section here."
	ExpectVocabulary         string = "Expect Vocabulary"
	ExpectLanguageExpression string = "The language of this section is not specified. Specified either (it) or (de)"
	UnrecognizedLanguage     string = "Unrecognized language identifier. Specify either (it) or (de)"
	ExpectVocabSection       string = "Expect Vocab Section"
	UnexpectedToken          string = "Unexpected Token"
	InvalidScore             string = "Score must be a number"
	DuplicateToken           string = "Duplicate token in same section"
)

type Parser struct {
	ctx     context.Context
	uri     string
	scanner *Scanner
	Ast     *VocabAst

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
		Ast:     &VocabAst{},

		token:      TokenUnknown,
		text:       "",
		tokenStart: -1,
		tokenEnd:   -1,

		printCallback: printCallback,
	}
}

func (p *Parser) Parse() *Parser {
	p.Ast.Sections = []*VocabularySection{}

	var lastSection *VocabularySection = nil

	startNewSection := func() {
		p.Ast.Sections = append(p.Ast.Sections, NewVocabularySection(p.uri))
	}

	for {
		p.nextToken()

		if p.token == TokenEOF {
			return p
		}

		if len(p.Ast.Sections) > 0 {
			lastSection = p.Ast.Sections[len(p.Ast.Sections)-1]
		}

		switch p.token {
		case TokenCommentTrivia:
			continue
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
	section := p.currentVocabSection()
	date := &DateSection{
		Parent: section,
		Text:   p.text,
		Time:   parsedAsLocalTime,
		Start:  p.tokenStart,
		End:    p.tokenEnd, Line: p.line,
	}
	section.Date = date

	if err != nil {
		p.errorHere(&err, MalformedDate)
		return
	}

	p.nextToken()

}

func (p *Parser) parseVocabSection() {
	currentSection := p.currentVocabSection()
	words := &WordsSection{
		Line:   p.line,
		Parent: currentSection,
	}
	if p.token == TokenGreaterThan {
		words.Reviewed = false
		currentSection.NewWords = append(currentSection.NewWords, words)
	} else {
		words.Reviewed = true
		currentSection.ReviewedWords = append(currentSection.ReviewedWords, words)
	}

	p.nextTokenNotWhitespace()

	if p.token != TokenSemanticSpecifierLiteral {
		p.errorHere(nil, ExpectLanguageExpression)
		return
	}
	switch p.text {
	case "it":
		words.Language = Italiano
	case "de":
		words.Language = Deutsch
	default:
		p.errorHere(nil, UnrecognizedLanguage)
		words.Language = Unrecognized
	}

	p.nextTokenNotWhitespace()

	parsing := ""
	currentGrade := 0
	parsingStart := -1

	assignGradeAndClear := func() {
		words.Words[len(words.Words)-1].Grade = currentGrade
		currentGrade = 0
	}

	newWordFromText := func(t string) {
		isWordLiteral := p.token == TokenWordLiteral
		text := func() string {
			if isWordLiteral {
				return p.text
			}
			return t
		}()
		trailingSpaceCount := func() int {
			i := len(text) - 1
			trailingSpaceCount := 0
			for {
				if text[i] != ' ' {
					break
				}
				i--
				trailingSpaceCount++
			}
			return trailingSpaceCount
		}()
		start := parsingStart
		end := p.tokenStart - trailingSpaceCount

		newWord := &Word{Parent: words, Text: text[:len(text)-trailingSpaceCount], Start: start, End: end, Literally: isWordLiteral, Line: p.line}

		for _, word := range words.Words {
			if word.Text == newWord.Text {
				p.diagnosticsAt(nil, DuplicateToken, start, end, lsproto.DiagnosticsSeverityWarning)
				return
			}
		}

		words.Words = append(words.Words, newWord)

		parsingStart = -1
	}

	for {
		switch p.token {

		case TokenLineBreak, TokenEOF, TokenCommentTrivia:
			if parsing != "" {
				newWordFromText(parsing)
				assignGradeAndClear()
			}
			return
		case TokenSemanticSpecifierLiteral:
			number, err := strconv.Atoi(p.text)
			if err != nil {
				p.errorHere(nil, InvalidScore)
				for {
					if p.token == TokenLineBreak || p.token == TokenEOF {
						return
					}
					p.nextToken()
				}
			}
			currentGrade = number
			p.nextTokenNotWhitespace()
		case TokenWordLiteral:
			newWordFromText(parsing)
			p.nextTokenNotWhitespace()
		case TokenComma:
			if parsing == "" {
				assignGradeAndClear()
				p.nextTokenNotWhitespace()
				continue
			}

			newWordFromText(parsing)
			assignGradeAndClear()
			p.nextTokenNotWhitespace()
			parsing = ""
		default:
			if parsingStart == -1 {
				parsingStart = p.tokenStart
			}
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
			newUtterance := &UtteranceSection{
				Parent: p.currentVocabSection(),
				Line:   p.line,
				Start:  start,
				End:    start + len(text),
				Text:   sb.String(),
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
	p.diagnosticsAt(original, message, p.tokenStart, p.tokenEnd, lsproto.DiagnosticsSeverityError)

	// Only 1 error per line for simplicity.
	for {
		p.nextToken()

		if p.token == TokenLineBreak || p.token == TokenEOF {
			break
		}
	}
}

func (p *Parser) diagnosticsAt(original *error, message string, start int, end int, severity lsproto.DiagnosticsSeverity) {
	if original != nil {
		p.printCallback(&original)
	}
	newError := &lsproto.Diagnostic{
		Severity: severity,
		Message:  string(message),
		Range: lsproto.Range{
			Start: lsproto.Position{
				Line:      p.line,
				Character: start,
			},
			End: lsproto.Position{
				Line:      p.line,
				Character: end,
			},
		},
	}
	diag := p.currentVocabSection().Diagnostics
	p.currentVocabSection().Diagnostics = append(diag, newError)

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
	sectionCount := len(p.Ast.Sections)
	if sectionCount == 0 {
		return nil
	}
	last := p.Ast.Sections[sectionCount-1]
	return last
}

func (p *Parser) nextToken() {
	p.line = p.scanner.line
	token, text := p.scanner.Scan()
	p.text = text
	p.token = token
	p.tokenEnd = p.scanner.tokenLineOffsetEnd
	p.tokenStart = p.scanner.tokenLineOffsetStart
}
