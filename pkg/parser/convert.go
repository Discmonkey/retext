package parser

import (
	"errors"
	"github.com/discmonkey/retext/pkg/swagger"
)

type DocumentType int

const (
	Text DocumentType = iota
	DocX
	Xlsx
	Csv
)

type Sentence = swagger.Sentence
type Paragraph = swagger.Paragraph
type Source = swagger.Source
type Demo = swagger.Demo

type Parser interface {
	ConvertSource([]byte) (Source, error)
	ConvertDemo([]byte) (Demo, error)
}

// GetParser returns the matching parser for the filetype or error if there is no correct parser
func NewParser(t DocumentType) (Parser, error) {
	switch t {
	case Text:
		return TxtParser{}, nil
	case DocX:
		return DocXParser{}, nil
	case Xlsx:
		return XlsxParser{}, nil
	case Csv:
		return CsvParser{}, nil
	default:
		return nil, errors.New("could not find registered parser for doc type")
	}
}

// ConvertSource takes the bytes from a file and the file type and generates a simple intermediate representation
// which can be manipulated.
func ConvertSource(unprocessed []byte, t DocumentType) (Source, error) {
	p, err := NewParser(t)
	if err != nil {
		return Source{}, err
	}

	return p.ConvertSource(unprocessed)
}

/// ConvertDemo
func ConvertDemo(unprocessed []byte, t DocumentType) (Demo, error) {
	p, err := NewParser(t)
	if err != nil {
		return Demo{}, err
	}

	return p.ConvertDemo(unprocessed)
}
