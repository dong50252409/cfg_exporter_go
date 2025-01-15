package json

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"cfg_exporter/render"
	"encoding/json"
	"github.com/stoewer/go-strcase"
	"os"
	"path/filepath"
)

type jsonRender struct {
	*entities.Table
}

func init() {
	render.Register("json", newJSONRender)
}

func newJSONRender(table *entities.Table) render.IRender {
	return &jsonRender{table}
}

func (r *jsonRender) Execute() error {
	jsonDir := r.ExportDir()
	if err := os.MkdirAll(jsonDir, os.ModePerm); err != nil {
		return err
	}

	jsonFilename := filepath.Join(jsonDir, r.Filename())
	fileIO, err := os.Create(jsonFilename)
	if err != nil {
		return err
	}
	defer func() { _ = fileIO.Close() }()

	var dataList = make([]map[string]interface{}, 0, len(r.DataSet))
	for _, rowData := range r.DataSet {
		rowMap := make(map[string]interface{}, len(r.Fields))
		for fieldIndex, field := range r.Fields {
			v := rowData[fieldIndex]
			switch v {
			case nil:
				continue
			case "":
				continue
			default:
				rowMap[field.Name] = convert(v)
			}
		}
		dataList = append(dataList, rowMap)
	}

	// 序列化为 JSON
	jsonData, err := json.MarshalIndent(dataList, "   ", "   ")
	if err != nil {
		return err
	}

	_, err = fileIO.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func (r *jsonRender) ExportDir() string {
	erlang := config.Config.Schema["json"]
	return erlang.Destination
}

func (r *jsonRender) Filename() string {
	return config.Config.Schema["json"].FilePrefix + r.Name + ".json"
}

func (r *jsonRender) ConfigName() string {
	return strcase.SnakeCase(config.Config.Schema["json"].RecordPrefix + r.Name)
}
