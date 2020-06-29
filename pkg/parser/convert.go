package parser

import "errors"

type DocumentType int

const (
	Text DocumentType = iota
	DocX
	Xlsx
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

type Attributes struct {
	Columns []string
	Values  map[string][]string
}

type Document struct {
	Attributes Attributes
	Title      string
	Paragraphs []Paragraph
}

type Parser interface {
	Convert([]byte) (Document, error)
}

// GetParser returns the correct parser for the filetype or error if there is no correct parser
// TODO may be worth adding an extra method to the Parser interface that returns a list of parseable documents
// TODO then parsers can automatically declare themselves. On the other hand we would then need to keep around
// TODO instances of the different parsers. The current approach uses a bit less state, with a lot of short lived
// TODO parsers. This feels appropriate since parsers should never have state that extends past the lifetime of
// TODO the document that they have just finished parsing.
func GetParser(t DocumentType) (Parser, error) {
	switch t {
	case Text:
		return TxtParser{}, nil
	case DocX:
		return DocXParser{}, nil
	case Xlsx:
		return XlsxParser{}, nil
	default:
		return nil, errors.New("could not find registered parser for doc type")
	}
}

// Convert takes the bytes from a file and the file type and generates a simple intermediate representation
// which can be manipulated.
func Convert(unprocessed []byte, t DocumentType) (Document, error) {
	p, err := GetParser(t)
	if err != nil {
		return Document{}, err
	}

	return p.Convert(unprocessed)
}
