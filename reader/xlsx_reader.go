package reader

import (
	"cfg_exporter/config"
	"cfg_exporter/util"
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

	defer func() { _ = file.Close() }()

	var records [][]string
	filename := filepath.Base(path)
	tableName, err := util.SubTableName(filename)
	if err != nil {
		return nil, err
	}
	for _, sheetName := range file.GetSheetList() {
		if name, _ := util.SubTableName(sheetName); tableName == name {
			if tempRecords, err := file.GetRows(sheetName); err != nil {
				return nil, err
			} else {
				if len(records) == 0 {
					records = tempRecords
				} else {
					combineRecords(&records, &tempRecords)
				}

			}

		}
	}

	return records, err
}

func combineRecords(records *[][]string, newRecords *[][]string) {
	fieldNameRow := config.Config.Schema[config.SchemaName].FieldNameRow
	bodyStartRow := config.Config.BodyStartRow

	r1 := (*records)[fieldNameRow-1]
	r2 := (*newRecords)[fieldNameRow-1]

	set1 := make(map[string]int)
	for index, fieldName := range r1 {
		set1[fieldName] = index
	}

	set2 := make(map[string]int)
	for index, fieldName := range r2 {
		set2[fieldName] = index
	}
	lastIndex := len(*records)
	for index, val1 := range r1 {
		val2 := r2[index]
		if val1 == val2 {
			// 总sheet和新sheet里都有相同字段，将新sheet里的数据追加到总sheet里
			for rowIndex, row := range (*newRecords)[bodyStartRow-1:] {
				if len(*records) <= lastIndex+rowIndex {
					*records = append(*records, []string{})
				}
				if len(row) > index {
					(*records)[lastIndex+rowIndex] = append((*records)[lastIndex+rowIndex], row[index])
				} else {
					(*records)[lastIndex+rowIndex] = append((*records)[lastIndex+rowIndex], "")
				}
			}
		} else {
			if colIndex, ok := set2[val1]; ok {
				// 新sheet里有总sheet的字段，只不过不在同一列，将新sheet里的数据追加到总sheet里
				for rowIndex, row := range (*newRecords)[bodyStartRow-1:] {
					if len(*records) <= lastIndex+rowIndex {
						*records = append(*records, []string{})
					}
					if len(row) > colIndex {
						(*records)[lastIndex+rowIndex] = append((*records)[lastIndex+rowIndex], row[colIndex])
					} else {
						(*records)[lastIndex+rowIndex] = append((*records)[lastIndex+rowIndex], "")
					}
				}
			} else {
				// 新sheet里没有总sheet的字段，创建空数据追加到总sheet里
				for rowIndex := 0; rowIndex < len((*newRecords)[bodyStartRow-1:]); rowIndex++ {
					if len(*records) <= lastIndex+rowIndex {
						*records = append(*records, []string{})
					}
					(*records)[lastIndex+rowIndex] = append((*records)[lastIndex+rowIndex], "")
				}
			}
		}
	}
	for index, val2 := range r2 {
		if _, ok := set1[val2]; !ok {
			// 总sheet里没有新sheet有的字段，将新sheet里的数据加到总sheet的最后面
			for rowIndex, row := range (*newRecords)[0:bodyStartRow] {
				if len(row) > index {
					(*records)[rowIndex] = append((*records)[rowIndex], row[index])
				} else {
					(*records)[rowIndex] = append((*records)[rowIndex], "")
				}
			}

			for rowIndex, row := range (*newRecords)[bodyStartRow-1:] {
				if len(row) > index {
					(*records)[lastIndex+rowIndex] = append((*records)[lastIndex+rowIndex], row[index])
				} else {
					(*records)[lastIndex+rowIndex] = append((*records)[lastIndex+rowIndex], "")
				}
			}
		}
	}
}
