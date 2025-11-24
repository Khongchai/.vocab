# vocab

Extension that adds a support for a custom .vocab file extension. 

.vocab is a custom language for vocabulary note taking format. It uses the spaced-repetition sm2 algorithm behind the scene to help remind you which words need review.

Right now it supports 3 languages identifiers: de, fr, and it. Even though these don't mean much other than just helping to remind you which language branch the words or phrases fall in, in the future, I might extend it with NLP capabilities to perform lemmatization. We'll see.

# Syntax

## Structure

The syntax is composed of 4 main sections, the date, new words, reviewed words, and utterance.

```
04/09/2025
> (de) schön
>> (it) `inoltre`
Inoltre, questo plugin sarà fantastico. Sono sicuro.
```

### Date

The date section has `dd/mm/yyyy` format. This marks the start of a section and plays into when the word is going to circle back.

### New and Reviewed Sections

`>` marks new and `>>` marks review. This is for clarity only and does not have any language server functionality.

The words in these sections are demarcated by a comma

```
> (de|fr|it) word1, word2
```

### Utterance

This section is for example sentences. Can be anything.

The utterance section can be any string of text for however many lines. The text are considered utterance sections up until the next date section.

## Grading

A word can be graded according to the [sm2](https://en.wikipedia.org/wiki/SuperMemo#Description_of_SM-2_algorithm) algorithm.

```
> (de|fr|it) word1(5), word2(2)
```

This affect the interval between the word's last appearance and when it needs to appear (be reviewed) again.

## Exact Match

Capture exact match by wrapping a word with backticks. 

`das Haus`, `word2`

The result is that the entire `das Haus` must reappear again later -- normally both indefinite and definite articles are stripped out.

## Comment

Comments are prepended with the pipe symbol `|`.

```
01/01/2026
> (de|fr|it) word1(5), word2(2) | word2 is so difficult...damn!
```

## Commands

`Review All From This File` will create a new section with words in the current file that needs review.

`Review All` will create a new section with words in all files in the current workspace that ends with .vocab that needs review.

# Example
```
13/10/2025
> (it) qualcuno(1), migliaia, decimi(1)
Qualcuno di voi ha chiesto:
Finalmente siamo arrivati, dopo migliaia di chilometri e decimi di ore di macchini.
Nach Tausenden von Kilometern und Zehntelstunden Fahrt sind wir endlich angekommen.
14/10/2025
> (de) der Nebensatz, der Relativsatz(3), `der Einschub`, die Schnodderigkeit(4)
Der Sprecher benutzt lange, zusammengesetzte Sätze mit Nebensätzen, Relativsätzen und erklärenden Einschüben.
Er ist voll mit Schnodderigkeit. Kann nicht mit ihm arbeiten...
```

