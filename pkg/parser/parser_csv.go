package parser

import (
	"bytes"
	"encoding/csv"
	"io"
)

type CsvRow struct {
	values []string
	err    error
}

type CsvParser struct {
}

func (x CsvParser) Convert(unprocessed []byte) (Document, error) {
	doc := Document{
		Attributes: Attributes{
			Columns: make([]string, 0, 0),
			Values:  make(map[string][]string),
		},
	}

	rowChannel := getCsvRows(bytes.NewReader(unprocessed))
	headerRow := <-rowChannel

	if headerRow.err != nil {
		return doc, headerRow.err
	}

	a := Attributes{
		Columns: headerRow.values,
		Values:  make(map[string][]string),
	}

	numCols := len(a.Columns)
	for _, columnName := range a.Columns {
		a.Values[columnName] = make([]string, 0)
	}

	for csvRow := range rowChannel {
		if csvRow.err != nil {
			return doc, csvRow.err
		}
		if len(csvRow.values) != numCols {
			continue
		}

		for j, val := range csvRow.values {
			columnName := a.Columns[j]
			a.Values[columnName] = append(a.Values[columnName], val)
		}
	}

	doc.Attributes = a

	return doc, nil
}

// mostly https://stackoverflow.com/questions/32027590/efficient-read-and-write-csv-in-go
func getCsvRows(rc io.Reader) chan CsvRow {
	ch := make(chan CsvRow, 10)
	go func() {
		r := csv.NewReader(rc)
		defer close(ch)
		for {
			rec, err := r.Read()
			if err == io.EOF {
				break
			}

			r := CsvRow{
				values: rec,
				err:    err,
			}

			ch <- r

			if err != nil {
				break
			}
		}
	}()
	return ch
}

var _ Parser = CsvParser{}
