package parser

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strings"
)

type XlsxParser struct {
}

func parseSheet(file *zip.File) Attributes {
	return Attributes{}
}

func parseStrings(file *zip.File) []string {
	return nil
}

func mergeAttributes(source, destination Attributes) {
	for _, column := range source.Columns {
		destination.Columns = append(destination.Columns, column)
	}

	for str, attributeSet := range source.Values {
		var combined []string

		if value, ok := destination.Values[str]; ok {
			combined = value

			for _, val := range attributeSet {
				combined = append(combined, val)
			}
		} else {
			length := len(source.Columns) + len(destination.Columns)
			combined = make([]string, length, length)

			for i, val := range attributeSet {
				combined[len(source.Columns)+i] = val
			}
		}
	}
}

func (x XlsxParser) Convert(unprocessed []byte) (Document, error) {
	d := Document{
		Attributes: Attributes{
			Columns: make([]string, 0, 0),
			Values:  make(map[string][]string),
		},
	}

	reader, err := zip.NewReader(bytes.NewReader(unprocessed), int64(len(unprocessed)))
	if err != nil {
		return d, err
	}

	var sheet *zip.File = nil
	var sharedStrings *zip.File = nil

	for _, file := range reader.File {

		if strings.HasSuffix(file.Name, "sheet1.xml") {
			sheet = file
		} else if strings.HasSuffix(file.Name, "sharedStrings.xml") {
			sharedStrings = file
		}
	}

	contentsSheet, err := readZipFile(sheet)
	if err != nil {
		return d, err
	}

	contentsStrings, err := readZipFile(sharedStrings)
	if err != nil {
		return d, err
	}

	fmt.Println(string(contentsSheet))
	fmt.Println(string(contentsStrings))

	return d, nil
}

var _ Parser = XlsxParser{}
