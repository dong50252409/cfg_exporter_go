package json

import (
	"cfg_exporter/render"
	"encoding/json"
	"github.com/stoewer/go-strcase"
	"os"
	"path/filepath"
	"reflect"
	"sort"
)

type JSONRender struct {
	*render.Render
}

func init() {
	render.Register("json", newJSONRender)
}

func newJSONRender(render *render.Render) render.IRender {
	return &JSONRender{render}
}

func (r *JSONRender) Execute() error {
	if err := r.Render.Before(); err != nil {
		return err
	}

	jsonDir := r.ExportDir()
	if err := os.MkdirAll(jsonDir, os.ModePerm); err != nil {
		return err
	}

	fp := filepath.Join(jsonDir, r.Filename())
	fileIO, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer func() { _ = fileIO.Close() }()

	var dataList = make([]map[string]any, 0, len(r.DataSet))
	for _, rowData := range r.DataSet {
		rowMap := make(map[string]any, len(r.Fields))
		for fieldIndex, field := range r.Fields {
			switch v := rowData[fieldIndex]; v {
			case nil, "":
				continue
			default:
				v1 := convert(v)
				rowMap[strcase.LowerCamelCase(field.Name)] = v1
			}
		}
		dataList = append(dataList, rowMap)
	}
	pkFields := r.GetPrimaryKeyFields()
	sort.Slice(dataList, func(i, j int) bool {
		for _, field := range pkFields {
			switch field.Type.Kind() {
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				switch v1, v2 := dataList[i][field.Name].(int64), dataList[j][field.Name].(int64); {
				case v1 < v2:
					return true
				case v1 > v2:
					return false
				}
			case reflect.Float32, reflect.Float64:
				switch v1, v2 := dataList[i][field.Name].(float64), dataList[j][field.Name].(float64); {
				case v1 < v2:
					return true
				case v1 > v2:
					return false
				}
			case reflect.String:
				switch v1, v2 := dataList[i][field.Name].(string), dataList[j][field.Name].(string); {
				case v1 < v2:
					return true
				case v1 > v2:
					return false
				}
			default:
				return true
			}
		}
		return true
	})

	var macroMap = make(map[string]any)
	for _, macro := range r.GetMacroDecorators() {
		var childMacroMap = make(map[string]any)
		for _, macroDetail := range macro.List {
			childMacroMap[strcase.UpperSnakeCase(macroDetail.Key)] = macroDetail.Value
		}
		macroMap[strcase.LowerCamelCase(macro.MacroName)] = childMacroMap
	}

	rootMap := map[string]any{
		r.ConfigName() + "List": dataList,
	}
	if len(macroMap) > 0 {
		rootMap[r.ConfigName()+"MacroMap"] = macroMap
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

	if err := r.Render.After(); err != nil {
		return err
	}

	return nil
}

func (r *JSONRender) Verify() error {
	return nil
}

func (r *JSONRender) ConfigName() string {
	return strcase.LowerCamelCase(r.Schema.TableNamePrefix + r.Name)
}

func (r *JSONRender) Filename() string {
	return strcase.KebabCase(r.Schema.FilePrefix+r.Name) + ".json"
}
