package reader

import (
	"encoding/csv"
	"os"
)

type CSVReader struct {
}

func init() {
	Register("csv", &CSVReader{})
}

func (r *CSVReader) CheckSupport(_ string) bool {
	return true
}

func (r *CSVReader) Read(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}
