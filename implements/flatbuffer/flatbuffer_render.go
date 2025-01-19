package flatbuffer

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"cfg_exporter/implements/flatbuffer/typesystem"
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

const entryTemplate = `
{{- define "entry" -}}
{{- range $tableName, $fields := .Table.GetEntries -}}
table {{ $tableName | toUpperCamelCase }}{
	{{- range $fieldName, $fieldType := $fields }}
	{{ $fieldName | toLowerCamelCase }}: {{ $fieldType }};
	{{- end }}
}
{{- end -}}
{{- end -}}
`

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
{{- "namespace" }} {{ .Table.Namespace }};

{{ template "entry" . }}

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

	for _, tmplStr := range []string{entryTemplate, dataSetTemplate, tailTemplate, fbTemplate} {
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

func (r *fbRender) GetEntries() any {
	var entries map[string]any
	for _, field := range r.Table.Fields {
		switch field.Type.(type) {
		case *typesystem.FBList:
			switch field.Type.(*typesystem.FBList).ITypeSystem.(*entities.List).T.(type) {
			case *typesystem.FBInteger, *typesystem.FBFloat, *typesystem.FBBoolean, *typesystem.FBStr, *typesystem.FBLang, *typesystem.FBRaw:
				continue
			default:
				switch m := r.GetEntries(); m.(type) {
				case map[string]string:
					// 返回的字段名和字段类型
					entries[field.Name] = m
				case map[string]map[string]string:
					// 返回的表和表所包含的字段和字段类型
					for k, v := range m.(map[string]map[string]string) {
						entries[k] = v
					}
				}
			}
		case *typesystem.FBTuple:
			switch field.Type.(*typesystem.FBTuple).ITypeSystem.(*entities.Tuple).T.(type) {
			case *typesystem.FBInteger, *typesystem.FBFloat, *typesystem.FBBoolean, *typesystem.FBStr, *typesystem.FBLang, *typesystem.FBRaw:
				continue
			default:
				switch m := r.GetEntries(); m.(type) {
				case map[string]string:
					// 返回的字段名和字段类型
					entries[field.Name] = m
				case map[string]map[string]string:
					// 返回的表和表所包含的字段和字段类型
					for k, v := range m.(map[string]map[string]string) {
						entries[k] = v
					}
				}
			}
		case *typesystem.FBMap:
			switch field.Type.(*typesystem.FBMap).ITypeSystem.(*entities.Map).ValueT.(type) {
			case *typesystem.FBInteger, *typesystem.FBFloat, *typesystem.FBBoolean, *typesystem.FBStr, *typesystem.FBLang, *typesystem.FBRaw:
				return map[string]string{
					"k": field.Type.(*typesystem.FBMap).ITypeSystem.(*entities.Map).KeyT.String(),
					"v": field.Type.(*typesystem.FBMap).ITypeSystem.(*entities.Map).ValueT.String(),
				}
			default:
				switch m := r.GetEntries(); m.(type) {
				case map[string]string:
					// 返回的字段名和字段类型
					entries[field.Name] = m
				case map[string]map[string]string:
					// 返回的表和表所包含的字段和字段类型
					for k, v := range m.(map[string]map[string]string) {
						entries[k] = v
					}
				}
			}
		default:
			continue
		}
	}
	return entries
}
