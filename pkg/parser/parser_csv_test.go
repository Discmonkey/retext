package parser

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestConvertCsv(t *testing.T) {
	b, err := ioutil.ReadFile("test_documents/comma.csv")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("all good", func(t *testing.T) {
		parser := CsvParser{}

		document, err := parser.Convert(b)

		if err != nil {
			t.Fatal(err)
		}

		for col, values := range document.Attributes.Values {
			fmt.Println(col, values)
		}
	})
}
