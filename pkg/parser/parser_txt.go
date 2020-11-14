package parser

import (
	"errors"
	"strings"
)

type TxtParser struct {
}

func (t TxtParser) ConvertDemo(bytes []byte) (Demo, error) {
	return Demo{}, errors.New("text par")
}

// Convert for the TxtParser does not fail since it simply reads from the start of the bytes to the end
// It makes fairly arbitrary choices for how to handle certain datatypes.
func (t TxtParser) ConvertSource(unprocessed []byte) (Source, error) {
	d := Source{
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
		Words: []string{},
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
			w, newPosition, last := readWord(text, position)

			s.Words = append(s.Words, w)

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

func readWord(text []byte, position int) (string, int, isLast) {

	builder := strings.Builder{}

	for position < len(text) {
		switch text[position] {

		// non character
		case '\n', '\f', '\r', '\t', ' ':
			return builder.String(), position, true

		case '!', '?':
			builder.WriteByte(text[position])
			position++

			return builder.String(), position, true
		case '.':
			builder.WriteByte(text[position])
			position++

			if isTerminating(text, position-1) {
				return builder.String(), position, true
			}

		default:
			builder.WriteByte(text[position])
			position++
		}
	}
	return builder.String(), position, true
}

var _ Parser = TxtParser{}
