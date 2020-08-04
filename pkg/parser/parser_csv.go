package parser

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
)

type CsvParser struct {
}

func (x CsvParser) Convert(unprocessed []byte) (doc Document, err error) {
	d := Document{
		Attributes: Attributes{
			Columns: make([]string, 0, 0),
			Values:  make(map[string][]string),
		},
	}

	defer func() (Document, error) {
		if r := recover(); r != nil {
			return d, errors.New(fmt.Sprintf("%s", r))
		}
		return d, err
	}()

	rowChannel := getCSVRows(bytes.NewReader(unprocessed))
	headerRow := <-rowChannel

	a := Attributes{
		Columns: headerRow,
		Values:  make(map[string][]string),
	}
	numCols := len(a.Columns)
	for _, columnName := range a.Columns {
		a.Values[columnName] = make([]string, 0)
	}

	for row := range rowChannel {
		if len(row) != numCols {
			continue
		}

		for j, val := range row {
			columnName := a.Columns[j]
			a.Values[columnName] = append(a.Values[columnName], val)
		}
	}

	d.Attributes = a
	return d, err
}

// mostly https://stackoverflow.com/questions/32027590/efficient-read-and-write-csv-in-go
func getCSVRows(rc io.Reader) (ch chan []string) {
	ch = make(chan []string, 10)
	go func() {
		r := csv.NewReader(rc)
		defer close(ch)
		for {
			rec, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)

			}
			ch <- rec
		}
	}()
	return
}

var _ Parser = CsvParser{}
