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

func (x XlsxParser) ConvertSource(i []byte) (Source, error) {
	return Source{}, errors.New("excel files currently not supported as source files")
}

func check() {
	var _ Parser = XlsxParser{}
}

func (x XlsxParser) ConvertDemo(unprocessed []byte) (Demo, error) {
	demo := Demo{
		Columns: nil,
		Values:  nil,
	}

	reader, err := zip.NewReader(bytes.NewReader(unprocessed), int64(len(unprocessed)))
	if err != nil {
		return demo, err
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
		return demo, errors.New("could not find content in excel file")
	}

	contentsSheet, err := readZipFile(sheet)
	if err != nil {
		return demo, err
	}

	contentsStrings, err := readZipFile(sharedStrings)
	if err != nil {
		return demo, err
	}

	commonStrings, err := parseStrings(contentsStrings)
	if err != nil {
		return demo, err
	}

	return parseSheet(contentsSheet, commonStrings)

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

func parseSheet(contents []byte, values []string) (Demo, error) {
	demo := Demo{}
	root, err := xmltree.Parse(contents)
	if err != nil {
		return demo, err
	}

	if root.Children == nil {
		return demo, errors.New("root excel sheet element does not have any children")
	}

	var rowsCols *xmltree.Element = nil
	for _, child := range root.Children {
		if isRowsAndCols(&child) {
			rowsCols = &child
			break
		}
	}

	if rowsCols == nil {
		return demo, errors.New("could not find rows/cols child of root")
	}

	numCols := len(rowsCols.Children[0].Children)
	demo.Columns = make([]string, numCols, numCols)
	demo.Values = make(map[string][]string)

	for i := 0; i < numCols; i++ {
		columnName := getValue(rowsCols.Children[0].Children[i], values)
		demo.Columns[i] = columnName
		demo.Values[columnName] = make([]string, 0, 0)
	}

	for i := 1; i < len(rowsCols.Children); i++ {
		row := rowsCols.Children[i]

		if len(row.Children) != numCols {
			continue
		}

		for j := 0; j < numCols; j++ {
			demo.Values[demo.Columns[j]] = append(demo.Values[demo.Columns[j]], getValue(row.Children[j], values))
		}
	}

	return demo, nil
}
