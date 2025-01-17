package typescript

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"cfg_exporter/render"
	"github.com/stoewer/go-strcase"
	"os"
	"path/filepath"
	"text/template"
)

const headTemplate = `
{{- define "head" -}}
import { BaseData_Json } from "../../BaseData"
{{- end -}}
`

const classTemplate = `
{{- define "class" -}}
{{- $pkFields := .Table.GetPrimaryKeyFields -}}
{{- $pkLen := len $pkFields | add -1 -}}
export class {{ .Table.Name | toUpperCamelCase }} extends BaseData_Json {
    get({{ range $fieldIndex, $field := $pkFields }}{{ $field.Name | toLowerCamelCase }}: {{ $field.Type }}{{ if lt $fieldIndex $pkLen }}, {{ end }}{{ end }}): any {
        return super.get({{ range $fieldIndex, $field := $pkFields }}{{ $field.Name | toLowerCamelCase }}{{ if lt $fieldIndex $pkLen }}, {{ end }}{{ end }})
    }
}
{{- end -}}
`

const interfaceTemplate = `
{{- define "interface" -}}
export interface {{ .Table.ConfigName | toUpperCamelCase }} {
    {{- range $_, $field := .Table.Fields }}
    readonly {{ $field.Name | toLowerCamelCase }}{{ if not $field.IsPrimaryKey }}?{{ end }}: {{ $field.Type }};
    {{- end }}
}
{{- end -}}
`

const tsTemplate = `
{{- template "head" .}}

{{ template "class" .}}

{{ template "interface" .}}
`

type tsRender struct {
	*entities.Table
	schema config.Schema
}

func init() {
	render.Register("typescript", newtsRender)
}

func newtsRender(table *entities.Table) render.IRender {
	return &tsRender{table, config.Config.Schema["typescript"]}
}

func (r *tsRender) Execute() error {
	dir := r.ExportDir()
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	tsFilepath := filepath.Join(dir, r.Filename())
	fileIO, err := os.Create(tsFilepath)
	if err != nil {
		return err
	}
	defer func() { _ = fileIO.Close() }()

	// 必备数据
	data := map[string]any{"Table": r}

	// 解析模板字符串
	tmpl := template.New("ts").Funcs(entities.FuncMap)

	for _, tmplStr := range []string{headTemplate, classTemplate, interfaceTemplate, tsTemplate} {
		tmpl, err = tmpl.Parse(tmplStr)
		if err != nil {
			return err
		}
	}

	// 执行模板渲染并输出到文件
	err = tmpl.Execute(fileIO, data)
	if err != nil {
		return err
	}
	return nil
}

func (r *tsRender) ExportDir() string {
	return r.schema.Destination
}

func (r *tsRender) Filename() string {
	return strcase.KebabCase(r.schema.FilePrefix+r.Name) + ".ts"
}

func (r *tsRender) ConfigName() string {
	return r.schema.TableNamePrefix + r.Name
}
