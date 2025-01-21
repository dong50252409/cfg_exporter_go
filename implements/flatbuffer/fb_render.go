package flatbuffer

import (
	"cfg_exporter/entities"
	"cfg_exporter/implements/flatbuffer/fb_type"
	"fmt"
	"github.com/stoewer/go-strcase"
	"os"
	"path/filepath"
	"text/template"
)

type fbRender struct {
	*flatbufferRender
}

const entryTemplate = `
{{- define "entry" -}}
{{- range $tableName, $fields := .Table.GetEntries }}
table {{ $tableName }}{
    {{- range $fieldName, $fieldType := $fields }}
    {{ $fieldName | toLowerCamelCase }}: {{ $fieldType }};
    {{- end }}
}
{{ end -}}
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
{{- "namespace" }} {{ .Table.Namespace -}};
{{ template "entry" . }}
{{ template "dataSet" . }}

{{ template "tail" . }}
`

func (r *fbRender) Execute() error {
	dir := r.ExportDir()
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	fp := filepath.Join(dir, r.Filename())
	fileIO, err := os.Create(fp)
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

func (r *fbRender) Filename() string {
	return strcase.SnakeCase(r.schema.FilePrefix+r.Name) + ".fbs"
}

func (r *fbRender) Namespace() string {
	return r.schema.Namespace
}

func (r *fbRender) GetEntries() map[string]map[string]string {
	var entries = make(map[string]map[string]string)
	for _, field := range r.Table.Fields {
		getNested(field.Type, entries, 0)
	}
	return entries
}

func getNested(t entities.ITypeSystem, entries map[string]map[string]string, deep int) string {
	var tableName string
	switch t.(type) {
	case *fb_type.FBInteger, *fb_type.FBFloat, *fb_type.FBBoolean, *fb_type.FBStr, *fb_type.FBLang, *fb_type.FBRaw:
		return tableName
	case *fb_type.FBTuple:
		fbType := t.(*fb_type.FBTuple)
		baseType := fbType.ITypeSystem.(*entities.Tuple)
		tableName = getTableName(baseType, entries, deep, "Tuple")
		switch baseType.T.(type) {
		case *fb_type.FBInteger, *fb_type.FBFloat, *fb_type.FBBoolean, *fb_type.FBStr, *fb_type.FBLang, *fb_type.FBRaw:
			return fmt.Sprintf("[%s]", baseType.T.String())
		case *fb_type.FBTuple, *fb_type.FBList, *fb_type.FBMap:
			entries[tableName] = make(map[string]string, 1)
			entries[tableName]["e"] = getNested(baseType.T, entries, deep+1)

		}
	case *fb_type.FBList:
		fbType := t.(*fb_type.FBList)
		baseType := fbType.ITypeSystem.(*entities.List)
		tableName = getTableName(baseType, entries, deep, "List")
		switch baseType.T.(type) {
		case *fb_type.FBInteger, *fb_type.FBFloat, *fb_type.FBBoolean, *fb_type.FBStr, *fb_type.FBLang, *fb_type.FBRaw:
			return fmt.Sprintf("[%s]", baseType.T.String())
		case *fb_type.FBTuple, *fb_type.FBList, *fb_type.FBMap:
			entries[tableName] = make(map[string]string, 1)
			entries[tableName]["e"] = getNested(baseType.T, entries, deep+1)
		}
	case *fb_type.FBMap:
		fbType := t.(*fb_type.FBMap)
		baseType := fbType.ITypeSystem.(*entities.Map)
		tableName = getTableName(baseType, entries, deep, "Map")
		switch baseType.ValueT.(type) {
		case *fb_type.FBInteger, *fb_type.FBFloat, *fb_type.FBBoolean, *fb_type.FBStr, *fb_type.FBLang, *fb_type.FBRaw:
			entries[tableName] = make(map[string]string, 1)
			entries[tableName]["k"] = baseType.KeyT.String()
			entries[tableName]["v"] = baseType.ValueT.String()
		case *fb_type.FBTuple, *fb_type.FBList, *fb_type.FBMap:
			entries[tableName] = make(map[string]string, 1)
			entries[tableName]["k"] = baseType.KeyT.String()
			entries[tableName]["v"] = getNested(baseType.ValueT, entries, deep+1)
		}
	}
	return fmt.Sprintf("[%s]", tableName)
}

func getTableName(t entities.ITypeSystem, entries map[string]map[string]string, deep int, suffix string) string {
	var field *entities.Field
	switch t.(type) {
	case *entities.Tuple:
		field = t.(*entities.Tuple).Field
	case *entities.List:
		field = t.(*entities.List).Field
	case *entities.Map:
		field = t.(*entities.Map).Field
	}
	var tableName string
	if deep == 0 {
		tableName = strcase.UpperCamelCase(fmt.Sprintf("%s_Entry", field.Name))
	} else {
		tableName = strcase.UpperCamelCase(fmt.Sprintf("%s_%sEntry", field.Name, suffix))
	}
	if _, ok := entries[tableName]; ok {
		tableName = strcase.UpperCamelCase(fmt.Sprintf("%s_%s_%s_Entry", field.Name, suffix, suffix))
	}
	return tableName
}
