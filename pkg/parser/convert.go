package parser

import (
	"strings"
)

type Word struct {
	Text string

	// frontend stuff starting here
	Selected bool
}

type Sentence struct {
	Parts []Word
}

type Paragraph struct {
	Sentences []Sentence
}

type Document struct {
	Tags       []string
	Title      []string
	Paragraphs []Paragraph
}

// Convert takes the bytes from a file, and converts those bytes to an intermediate format by breaking the document down into
// paragraphs, sentences, and  words. For now punctuation and words are roughly equivalent. This is somewhat problematic in that stop != stop.
// however it makes the implementation significantly easier for now. Saving punctuation forces the representationt o also
// keep tracing spaces for example "some - text" versus "some-text" versus "12/30/2020".
func Convert(unprocessed []byte) Document {

	d := Document{
		Paragraphs: make([]Paragraph, 0, 0)}

	cursor := 0
	length := len(unprocessed)
	for cursor < length {
		switch unprocessed[cursor] {
		case ' ', '\t', '\n', '\f', '\r':
			cursor++
		default:
			paragraph, cursorAfterParagraph := readParagraph(unprocessed, cursor)

			d.Paragraphs = append(d.Paragraphs, paragraph)

			cursor = cursorAfterParagraph
		}
	}

	return d
}

func readParagraph(text []byte, position int) (Paragraph, int) {
	p := Paragraph{
		Sentences: []Sentence{},
	}
	for position < len(text) {
		switch text[position] {

		// we assume we hit another paragraph even without punctuation
		case '\n', '\f', '\r':
			return p, position

		// skip spaces
		case '\t', ' ':
			position += 1

		// otherwise read the sentence
		default:
			sentence, newPosition := readSentence(text, position)

			p.Sentences = append(p.Sentences, sentence)

			position = newPosition
		}
	}

	return p, position
}

func readSentence(text []byte, position int) (Sentence, int) {
	s := Sentence{
		Parts: []Word{},
	}

	for position < len(text) {
		switch text[position] {

		// basically, only the paragraph reader is allowed to consume newlines
		case '\n', '\f', '\r':
			return s, position

		case ' ', '\t':
			position += 1
		// on any space we go ahead and read a word
		default:
			p, newPosition, last := readWord(text, position)

			s.Parts = append(s.Parts, p)

			position = newPosition

			if last {
				return s, position
			}
		}
	}

	return s, position
}

type isLast = bool

func readWord(text []byte, position int) (Word, int, isLast) {

	word := Word{}

	builder := strings.Builder{}

	for position < len(text) {
		switch text[position] {

		// as good as an ending
		case '\n', '\f', '\r':
			word.Text = builder.String()

			return word, position, true

		// on any space we go ahead and read a word
		case '\t', ' ':
			word.Text = builder.String()

			return word, position, false

		// these are ending characters
		case '.', '!', '?':
			builder.WriteByte(text[position])
			position++
			word.Text = builder.String()

			return word, position, true

		default:
			builder.WriteByte(text[position])
			position++
		}
	}

	word.Text = builder.String()

	return word, position, true
}
