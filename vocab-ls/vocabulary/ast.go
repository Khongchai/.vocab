package vocabulary

import "time"

type Word struct {
	text string
}

type Sentence struct {
	text string
	// future positions for reviewed and such will be here.
}

type VocabularySection struct {
	date          time.Time
	newWords      []Word
	reviewedWords []Word
	sentences     []Sentence
}

type Document struct {
	sections []VocabularySection
}

type Ast struct {
	documents []Document
}

func (ast *Ast) parse() {

}
