package flatbuffer

import (
	"encoding/json"
	"fmt"
	"github.com/stoewer/go-strcase"
	"os"
	"path/filepath"
)

type jsonRender struct {
	*flatbufferRender
}

func (r *jsonRender) Execute() error {
	dir := r.ExportDir()
	fp := filepath.Join(dir, r.Filename())
	fileIO, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer func() { _ = fileIO.Close() }()

	var dataList = make([]map[string]any, 0, len(r.DataSet))
	for _, rowData := range r.DataSet {
		rowMap := make(map[string]any, len(r.Fields))
		for fieldIndex, field := range r.Fields {
			v := rowData[fieldIndex]
			switch v {
			case nil:
				continue
			case "":
				continue
			default:
				v1 := convert(v)
				switch v1.(type) {
				case []interface{}:
					if field.Type.String() == "string" {
						marshal, err := json.Marshal(v1)
						if err != nil {
							return err
						}
						rowMap[strcase.LowerCamelCase(field.Name)] = fmt.Sprintf(`%s`, marshal)
					} else {
						rowMap[strcase.LowerCamelCase(field.Name)] = v1
					}
				default:
					rowMap[strcase.LowerCamelCase(field.Name)] = v1
				}
			}
		}
		dataList = append(dataList, rowMap)
	}

	rootMap := map[string]any{
		"dataSet": dataList,
	}

	// 序列化为 JSON
	jsonData, err := json.MarshalIndent(rootMap, "", "    ")

	if err != nil {
		return err
	}

	_, err = fileIO.Write(jsonData)
	if err != nil {
		return err
	}

	fmt.Printf("导出配置：%s\n", fp)

	return nil
}

func (r *jsonRender) Filename() string {
	return strcase.SnakeCase(r.schema.FilePrefix+r.Name) + ".json"
}
