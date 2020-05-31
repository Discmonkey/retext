package parser

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestConvertDocX(t *testing.T) {

	b, err := ioutil.ReadFile("test_documents/future.docx")
	if err != nil {
		t.Fatal(err)
	}

	parser := DocXParser{}

	document, err := parser.Convert(b)

	if err != nil {
		t.Fatal(err)
	}

	for _, paragraph := range document.Paragraphs {
		fmt.Println()

		for _, sentence := range paragraph.Sentences {

			for _, part := range sentence.Parts {
				fmt.Print(part.Text, " ")
			}
		}
	}

}
