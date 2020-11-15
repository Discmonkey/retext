package parser

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
)

type CsvRow struct {
	values []string
	err    error
}

type CsvParser struct {
}

func (x CsvParser) ConvertSource(i []byte) (Source, error) {
	return Source{}, errors.New("csv parsing currently not supported for source files")
}

func init() {
	var _ Parser = CsvParser{}
}

func (x CsvParser) ConvertDemo(unprocessed []byte) (Demo, error) {
	demo := Demo{
		Columns: nil,
		Values:  make(map[string][]string),
	}

	rowChannel := getCsvRows(bytes.NewReader(unprocessed))
	headerRow := <-rowChannel

	if headerRow.err != nil {
		return demo, headerRow.err
	}

	demo.Columns = headerRow.values

	numCols := len(demo.Columns)
	for _, columnName := range demo.Columns {
		demo.Values[columnName] = make([]string, 0)
	}

	for csvRow := range rowChannel {
		if csvRow.err != nil {
			return demo, csvRow.err
		}
		if len(csvRow.values) != numCols {
			continue
		}

		for j, val := range csvRow.values {
			columnName := demo.Columns[j]
			demo.Values[columnName] = append(demo.Values[columnName], val)
		}
	}

	return demo, nil
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
