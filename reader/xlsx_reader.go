package reader

import (
	"github.com/xuri/excelize/v2"
	"path/filepath"
	"strings"
)

type XLSXReader struct{}

func init() {
	Register("xlsx", &XLSXReader{})
}

func (r *XLSXReader) CheckSupport(path string) bool {
	return strings.Index(filepath.Base(path), "~$") == -1
}

func (r *XLSXReader) Read(path string) ([][]string, error) {
	file, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = file.Close()
	}()

	records, err := file.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}
	return records, nil
}
