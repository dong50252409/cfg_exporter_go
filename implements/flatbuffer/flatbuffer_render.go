package flatbuffer

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"cfg_exporter/render"
	"github.com/stoewer/go-strcase"
	"os"
	"path/filepath"
	"text/template"
)

type fbRender struct {
	*entities.Table
	schema config.Schema
}

const dataSetTemplate = `
{{- define "dataSet" -}}
table DataSet{
    {{- range $_, $field := .Table.Fields }}
    {{ $field.Name | toLowerCamelCase }}: {{ $field.Type }};
	{{- end }}
}
{{- end -}}
`

const tailTemplate = `
{{- define "tail" -}}
table Root{
    dataSet: [DataSet];
}

root_type Root;
{{- end -}}
`

const fbTemplate = `
namespace {{ .Table.Namespace }};

{{ template "dataSet" . }}

{{ template "tail" . }}
`

func init() {
	render.Register("flatbuffer", newtsRender)
}

func newtsRender(table *entities.Table) render.IRender {
	return &fbRender{table, config.Config.Schema["flatbuffer"]}
}

func (r *fbRender) Execute() error {
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
	tmpl := template.New("fb").Funcs(entities.FuncMap)

	for _, tmplStr := range []string{dataSetTemplate, tailTemplate, fbTemplate} {
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

func (r *fbRender) ExportDir() string {
	return r.schema.Destination
}

func (r *fbRender) Filename() string {
	return strcase.SnakeCase(r.schema.FilePrefix+r.Name) + ".fbs"
}

func (r *fbRender) ConfigName() string {
	return r.schema.TableNamePrefix + r.Name
}

func (r *fbRender) Namespace() string {
	return r.schema.Namespace
}
