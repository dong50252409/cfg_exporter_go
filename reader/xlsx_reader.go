package reader

import (
	"cfg_exporter/config"
	"cfg_exporter/util"
	"github.com/xuri/excelize/v2"
	"hash"
	"hash/fnv"
	"path/filepath"
)

var h hash.Hash32

type XLSXReader struct {
	*Reader
}

func init() {
	Register("xlsx", newXLSXReader)
	h = fnv.New32a()
}

func newXLSXReader(reader *Reader) IReader {
	return &XLSXReader{reader}
}

func (r *XLSXReader) Read() ([][]string, error) {
	file, err := excelize.OpenFile(r.Path)
	if err != nil {
		return nil, err
	}

	defer func() { _ = file.Close() }()

	var records [][]string
	filename := filepath.Base(r.Path)
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

// TODO 读取多个字段+装饰器来合并数据
func combineRecords(records *[][]string, newRecords *[][]string) {
	fieldCommentRow := config.Config.FieldCommentRow
	fieldNameRow := config.Config.Schema[config.SchemaName].FieldNameRow
	bodyStartRow := config.Config.BodyStartRow

	fnr1 := (*records)[fieldNameRow-1]
	fnr2 := (*newRecords)[fieldNameRow-1]

	fnSet1 := make(map[string]int)
	fnSet2 := make(map[string]int)

	// 字段名set
	for index, fieldName := range fnr1 {
		if fieldName != "" {
			fnSet1[fieldName] = index
		}
	}
	for index, fieldName := range fnr2 {
		if fieldName != "" {
			fnSet2[fieldName] = index
		}
	}

	// 备注set
	fcr1 := (*records)[fieldCommentRow-1]
	fcr2 := (*newRecords)[fieldCommentRow-1]

	fcSet1 := make(map[string]int)
	fcSet2 := make(map[string]int)

	for index, Comment := range fcr1 {
		if Comment != "" {
			fcSet1[Comment] = index
		}
	}
	for index, Comment := range fcr2 {
		if Comment != "" {
			fcSet2[Comment] = index
		}
	}

	lastIndex := len(*records)

	// 扩种总sheet表，为追加新sheet表数据做准备
	for i := 0; i < len((*newRecords)[bodyStartRow-1:]); i++ {
		*records = append(*records, []string{})
	}

	for i := 0; i < max(len(fnr1), len(fcr1)); i++ {
		if index := getIndex(i, fnr1, fcr1, fnSet2, fcSet2); index >= 0 {
			// 字段存在 新sheet表中有总sheet表的字段，将数据追加到总sheet表后面
			for rowIndex, row := range (*newRecords)[bodyStartRow-1:] {
				(*records)[lastIndex+rowIndex] = append((*records)[lastIndex+rowIndex], row[index])
			}
		} else {
			// 新sheet表中没有总sheet表的字段，补充空字符串
			for rowIndex := range (*newRecords)[bodyStartRow-1:] {
				(*records)[lastIndex+rowIndex] = append((*records)[lastIndex+rowIndex], "")
			}
		}
	}

	for i := 0; i < max(len(fnr2), len(fcr2)); i++ {
		if index := getIndex(i, fnr2, fcr2, fnSet1, fcSet1); index == -1 {
			// 总sheet表中没有新sheet表的字段，将数据追加到总sheet表的后面
			for rowIndex, row := range (*newRecords)[:bodyStartRow-1] {
				// 添加表头信息
				if len(row) > i {
					(*records)[rowIndex] = append((*records)[rowIndex], row[i])
				}
			}

			for rowIndex, row := range (*newRecords)[bodyStartRow-1:] {
				// 添加数据
				if len(row) > i {
					(*records)[lastIndex+rowIndex] = append((*records)[lastIndex+rowIndex], row[i])
				}
			}
		}
	}
}

func getIndex(i int, fnr, fcr []string, fnSet, fcSet map[string]int) int {
	if index, ok := fnSet[fnr[i]]; ok {
		return index
	} else if index, ok := fcSet[fcr[i]]; ok {
		return index
	}
	return -1
}

//func combineRecords(records *[][]string, newRecords *[][]string) {
//	fieldNameRow := config.Config.Schema[config.SchemaName].FieldNameRow
//	bodyStartRow := config.Config.BodyStartRow
//
//	r1 := (*records)[fieldNameRow-1]
//	r2 := (*newRecords)[fieldNameRow-1]
//
//	set1 := make(map[string]int)
//	for index, fieldName := range r1 {
//		set1[fieldName] = index
//	}
//
//	set2 := make(map[string]int)
//	for index, fieldName := range r2 {
//		set2[fieldName] = index
//	}
//	lastIndex := len(*records)
//	for index, val1 := range r1 {
//		val2 := r2[index]
//		if val1 == val2 {
//			// 总sheet和新sheet里都有相同字段，将新sheet里的数据追加到总sheet里
//			for rowIndex, row := range (*newRecords)[bodyStartRow-1:] {
//				if len(*records) <= lastIndex+rowIndex {
//					*records = append(*records, []string{})
//				}
//				if len(row) > index {
//					(*records)[lastIndex+rowIndex] = append((*records)[lastIndex+rowIndex], row[index])
//				} else {
//					(*records)[lastIndex+rowIndex] = append((*records)[lastIndex+rowIndex], "")
//				}
//			}
//		} else if colIndex, ok := set2[val1]; ok {
//			// 新sheet里有总sheet的字段，只不过不在同一列，将新sheet里的数据追加到总sheet里
//			for rowIndex, row := range (*newRecords)[bodyStartRow-1:] {
//				if len(*records) <= lastIndex+rowIndex {
//					*records = append(*records, []string{})
//				}
//				if len(row) > colIndex {
//					(*records)[lastIndex+rowIndex] = append((*records)[lastIndex+rowIndex], row[colIndex])
//				} else {
//					(*records)[lastIndex+rowIndex] = append((*records)[lastIndex+rowIndex], "")
//				}
//			}
//		} else {
//			// 新sheet里没有总sheet的字段，创建空数据追加到总sheet里
//			for rowIndex := 0; rowIndex < len((*newRecords)[bodyStartRow-1:]); rowIndex++ {
//				if len(*records) <= lastIndex+rowIndex {
//					*records = append(*records, []string{})
//				}
//				(*records)[lastIndex+rowIndex] = append((*records)[lastIndex+rowIndex], "")
//			}
//		}
//	}
//	for index, val2 := range r2 {
//		if _, ok := set1[val2]; !ok {
//			// 总sheet里没有新sheet有的字段，将新sheet里的数据加到总sheet的最后面
//			for rowIndex, row := range (*newRecords)[0:bodyStartRow] {
//				if len(row) > index {
//					(*records)[rowIndex] = append((*records)[rowIndex], row[index])
//				} else {
//					(*records)[rowIndex] = append((*records)[rowIndex], "")
//				}
//			}
//
//			for rowIndex, row := range (*newRecords)[bodyStartRow-1:] {
//				if len(row) > index {
//					(*records)[lastIndex+rowIndex] = append((*records)[lastIndex+rowIndex], row[index])
//				} else {
//					(*records)[lastIndex+rowIndex] = append((*records)[lastIndex+rowIndex], "")
//				}
//			}
//		}
//	}
//}
