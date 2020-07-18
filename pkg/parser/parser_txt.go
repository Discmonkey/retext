package parser

import (
	"strings"
)

type TxtParser struct {
}

// Convert for the TxtParser does not fail since it simply reads from the start of the bytes to the end
// It makes fairly arbitrary choices for how to handle certain datatypes.
func (t TxtParser) Convert(unprocessed []byte) (Document, error) {
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

	return d, nil
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

// isTerminating checks whether or not the current index is the last word in a paragraph
func isTerminating(text []byte, index int) bool {
	if index == len(text)-1 {
		return true
	}

	switch text[index+1] {
	case '\n', '\f', '\r', '\t', ' ':
		return true
	default:
		return false
	}
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

		case '!', '?':
			builder.WriteByte(text[position])
			position++
			word.Text = builder.String()

			return word, position, true

		case '.':
			builder.WriteByte(text[position])
			position++

			if isTerminating(text, position-1) {
				word.Text = builder.String()

				return word, position, true
			}

		default:
			builder.WriteByte(text[position])
			position++
		}
	}

	word.Text = builder.String()

	return word, position, true
}

var _ Parser = TxtParser{}
