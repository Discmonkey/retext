package parser

import (
	"aqwari.net/xml/xmltree"
	"archive/zip"
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"strings"
)

// DocxParser implements the converter interface
type DocXParser struct {
}

// Convert unzips the docx file, finds the content (document.xml), parses the xml tree, and finally walks the xml
// tree in a top down fashion
func (t DocXParser) Convert(unprocessed []byte) (Document, error) {
	d := Document{
		Paragraphs: make([]Paragraph, 0, 0)}

	reader, err := zip.NewReader(bytes.NewReader(unprocessed), int64(len(unprocessed)))
	if err != nil {
		return d, err
	}

	var main *zip.File = nil

	for _, file := range reader.File {

		if strings.HasSuffix(file.Name, "document.xml") {
			main = file
			break
		}
	}

	if main == nil {
		return d, errors.New("could not find document.xml in docx file")
	}

	contents, err := readFile(main)
	if err != nil {
		return d, err
	}

	root, err := xmltree.Parse(contents)

	if root == nil {
		return d, errors.New("could not parse document.xml xml tree")
	}

	topDownWalk(root, &documentBuilder{
		doc: &d,
	})

	return d, nil
}

func readFile(file *zip.File) ([]byte, error) {
	reader, err := file.Open()
	if err != nil {
		return nil, err
	}

	defer func() {
		err := reader.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	return ioutil.ReadAll(reader)
}

type visitor interface {
	Visit(element *xmltree.Element)
}

type documentBuilder struct {
	doc *Document
}

func (d *documentBuilder) Visit(element *xmltree.Element) {

	paragraphIndex := len(d.doc.Paragraphs) - 1

	// the second condition checks to confirm that we have added at least one sentence to the previous paragraph
	// before constructing a new paragraph
	if element.StartElement.Name.Local == "p" &&
		(paragraphIndex == -1 || len(d.doc.Paragraphs[paragraphIndex].Sentences) > 0) {

		d.doc.Paragraphs = append(d.doc.Paragraphs, Paragraph{
			Sentences: make([]Sentence, 0, 0)},
		)

		paragraphIndex += 1
	}

	if len(element.Children) == 0 {

		currentIndex := 0
		sentence := Sentence{}

		for currentIndex < len(element.Content) {
			sentence, currentIndex = readSentence(element.Content, currentIndex)
			d.doc.Paragraphs[paragraphIndex].Sentences = append(d.doc.Paragraphs[paragraphIndex].Sentences, sentence)
		}

	}
}

var _ visitor = &documentBuilder{}

func topDownWalk(element *xmltree.Element, v visitor) {
	v.Visit(element)

	for _, child := range element.Children {
		topDownWalk(&child, v)
	}
}
