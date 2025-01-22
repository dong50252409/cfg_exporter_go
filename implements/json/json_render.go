package json

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"cfg_exporter/render"
	"encoding/json"
	"fmt"
	"github.com/stoewer/go-strcase"
	"os"
	"path/filepath"
)

type jsonRender struct {
	*entities.Table
	schema config.Schema
}

func init() {
	render.Register("json", newJSONRender)
}

func newJSONRender(table *entities.Table) render.IRender {
	return &jsonRender{table, config.Config.Schema["json"]}
}

func (r *jsonRender) Execute() error {
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
			v := rowData[fieldIndex]
			switch v {
			case nil:
				continue
			case "":
				continue
			default:
				v1 := convert(v)
				if v1 != nil {
					rowMap[strcase.LowerCamelCase(field.Name)] = v1
				}
			}
		}
		dataList = append(dataList, rowMap)
	}

	var macroMap = make(map[string]any)
	for _, macro := range r.GetMacroDecorators() {
		var childMacroMap = make(map[string]any)
		for _, macroDetail := range macro.List {
			childMacroMap[strcase.UpperSnakeCase(macroDetail.Key)] = macroDetail.Value
		}
		macroMap[strcase.LowerCamelCase(macro.MacroName)] = childMacroMap
	}

	rootMap := map[string]any{
		"dataSet": dataList,
	}
	if len(macroMap) > 0 {
		rootMap["macroSet"] = macroMap
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

func (r *jsonRender) ExportDir() string {
	return r.schema.Destination
}

func (r *jsonRender) Filename() string {
	return strcase.KebabCase(r.schema.FilePrefix+r.Name) + ".json"
}

func (r *jsonRender) ConfigName() string {
	return ""
}
