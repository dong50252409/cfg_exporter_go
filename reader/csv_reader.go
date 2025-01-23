package reader

import (
	"encoding/csv"
	"os"
)

type CSVReader struct {
	*Reader
}

func init() {
	Register("csv", newCSVReader)
}

func newCSVReader(r *Reader) IReader {
	return &CSVReader{r}
}

func (r *CSVReader) Read() ([][]string, error) {
	file, err := os.Open(r.Path)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) { _ = file.Close() }(file)

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}
