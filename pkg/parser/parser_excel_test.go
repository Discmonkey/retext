package parser

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestConvertXlsx(t *testing.T) {

	b, err := ioutil.ReadFile("test_documents/excel.xlsx")
	if err != nil {
		t.Fatal(err)
	}

	parser := XlsxParser{}

	document, err := parser.Convert(b)

	if err != nil {
		t.Fatal(err)
	}

	for column := range document.Attributes.Columns {
		fmt.Println(column)
	}

}
