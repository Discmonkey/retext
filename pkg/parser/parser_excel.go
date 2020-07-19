package parser

import (
	"aqwari.net/xml/xmltree"
	"archive/zip"
	"bytes"
	"errors"
	"strconv"
	"strings"
)

type XlsxParser struct {
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

	if sheet == nil {
		return d, errors.New("could not find content in excel file")
	}

	contentsSheet, err := readZipFile(sheet)
	if err != nil {
		return d, err
	}

	contentsStrings, err := readZipFile(sharedStrings)
	if err != nil {
		return d, err
	}

	commonStrings, err := parseStrings(contentsStrings)
	if err != nil {
		return d, err
	}

	d.Attributes, err = parseSheet(contentsSheet, commonStrings)

	return d, err
}

type arrayVisitor struct {
	items []string
}

func (a *arrayVisitor) Visit(element *xmltree.Element) {
	if element.Children == nil {
		a.items = append(a.items, string(element.Content))
	}
}

var _ visitor = &arrayVisitor{}

func parseStrings(contents []byte) ([]string, error) {
	root, err := xmltree.Parse(contents)
	if err != nil {
		return nil, err
	}

	visitor := arrayVisitor{
		items: make([]string, 0, 0),
	}

	topDownWalk(root, &visitor)

	return visitor.items, nil
}

func isRowsAndCols(element *xmltree.Element) bool {
	// check that element isn't empty
	return element.StartElement.Name.Local == "sheetData"
}

func getValue(element xmltree.Element, values []string) string {
	strval := string(element.Children[0].Content)

	isString := false
	for _, attr := range element.StartElement.Attr {
		if attr.Name.Local == "t" {
			isString = attr.Value == "s"
		}
	}

	if !isString {
		return strval
	}

	tableEntry, err := strconv.ParseInt(strval, 10, 64)
	if err != nil {
		return strval
	}

	if int(tableEntry) >= len(values) {
		return strval
	}

	return values[tableEntry]
}

func parseSheet(contents []byte, values []string) (Attributes, error) {
	a := Attributes{}
	root, err := xmltree.Parse(contents)
	if err != nil {
		return a, err
	}

	if root.Children == nil {
		return a, errors.New("root excel sheet element does not have any children")
	}

	var rowsCols *xmltree.Element = nil
	for _, child := range root.Children {
		if isRowsAndCols(&child) {
			rowsCols = &child
			break
		}
	}

	if rowsCols == nil {
		return a, errors.New("could not find rows/cols child of root")
	}

	numCols := len(rowsCols.Children[0].Children)
	a.Columns = make([]string, numCols, numCols)
	a.Values = make(map[string][]string)

	for i := 0; i < numCols; i++ {
		columnName := getValue(rowsCols.Children[0].Children[i], values)
		a.Columns[i] = columnName
		a.Values[columnName] = make([]string, 0, 0)
	}

	for i := 1; i < len(rowsCols.Children); i++ {
		row := rowsCols.Children[i]

		if len(row.Children) != numCols {
			continue
		}

		for j := 0; j < numCols; j++ {
			a.Values[a.Columns[j]] = append(a.Values[a.Columns[j]], getValue(row.Children[j], values))
		}
	}

	return a, nil
}

var _ Parser = XlsxParser{}
