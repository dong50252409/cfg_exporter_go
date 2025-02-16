package erlang

import (
	"cfg_exporter/entities"
	"fmt"
	"github.com/stoewer/go-strcase"
	"os"
	"path/filepath"
	"text/template"
)

type ERender struct {
	*ERLRender
}

const erlHeadTemplate = `
{{- define "head" -}}
%% Auto Create, Don't Edit
-module({{ .Table.ConfigName | toSnakeCase }}).
-include("{{ .Table.ConfigName | toSnakeCase }}.hrl").
-compile(export_all).
-compile(nowarn_export_all).
-compile({no_auto_import, [get/1]}).
{{- end -}}
`

const erlGetTemplate = `
{{- define "get" -}}
{{/* 声明模板渲染所需的变量 */}}
{{- $configName := .Table.ConfigName | toLower -}}
{{- $fields := .Table.Fields -}}
{{- $dataSet := .Table.DataSet -}}
{{- $dsLastIndex := index $dataSet 0 | len | add -1 -}}
{{- $pkValuesList := .Table.GetPrimaryKeyValuesByString -}}
{{- $pkLen := index $pkValuesList 0 | len -}}
{{- $pkLastIndex := $pkLen | add -1 -}}
{{- $pkSeq := seq $pkLen -}}

{{- range $rowIndex, $dataRow := $dataSet -}}
get({{ index $pkValuesList $rowIndex | joinByComma }})->
    #{{ $configName | toSnakeCase }}{
        {{- range $fieldIndex, $field := $fields }}
        {{ $field.Name | toSnakeCase }} = {{ index $dataRow $field.ColIndex | $field.Convert }}{{ if lt $fieldIndex $dsLastIndex }},{{ end }}
        {{- end }}
    };
{{ end -}}

get({{ range $pkIndex, $_ := $pkSeq }}ID{{ $pkIndex }}{{ if lt $pkIndex $pkLastIndex }}, {{ end }}{{ end }}) ->
    throw({config_error, ?MODULE, {{ range $pkIndex, $_ := $pkSeq }}ID{{ $pkIndex }}{{ if lt $pkIndex $pkLastIndex }}, {{ end }}{{ end }}}).
{{- end -}}
`
const erlListTemplate = `
{{- define "list" -}}
{{/* 声明模板渲染所需的变量 */}}
{{- $pkValuesList := .Table.GetPrimaryKeyValuesByString -}}
{{- $pkLastIndex := len $pkValuesList | add -1 -}}
{{- $pkLen := .Table.GetPrimaryKeyFields | len -}}

{{- if eq $pkLen 1 -}}
list() ->
    [
	{{- range $pkIndex, $pkValues := $pkValuesList -}}
	{{ $pkValues | joinByComma }}{{ if lt $pkIndex $pkLastIndex }}, {{ end }}
	{{- end -}}
    ].
{{- else -}}
list() ->
    [
	{{- range $pkIndex, $pkValues := $pkValuesList -}}
	{{ "{" }}{{ $pkValues | joinByComma }}{{ "}" }}{{ if lt $pkIndex $pkLastIndex }}, {{ end }}
	{{- end -}}
    ].
{{- end -}}
{{- end -}}
`

const erlTemplate = `
{{- template "head" .}}

{{ template "get" .}}

{{ template "list" .}}
`

func (r *ERender) Execute() error {
	erlDir := r.ExportDir()
	if err := os.MkdirAll(erlDir, os.ModePerm); err != nil {
		return err
	}

	fp := filepath.Join(erlDir, r.Filename())
	fileIO, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer func() { _ = fileIO.Close() }()

	// 必备数据
	data := map[string]any{"Table": r}

	// 解析模板字符串
	tmpl := template.New("erl").Funcs(entities.FuncMap)

	for _, tmplStr := range []string{erlHeadTemplate, erlGetTemplate, erlListTemplate, erlTemplate} {
		if tmpl, err = tmpl.Parse(tmplStr); err != nil {
			return err
		}
	}

	// 执行模板渲染并输出到文件
	if err = tmpl.Execute(fileIO, data); err != nil {
		return err
	}

	fmt.Printf("导出配置：%s\n", fp)

	return nil
}

// ExportDir 导出目录
func (r *ERender) ExportDir() string {
	return strcase.SnakeCase(filepath.Join(r.ERLRender.ExportDir(), "erl"))
}

// Filename 文件名
func (r *ERender) Filename() string {
	return strcase.SnakeCase(r.Schema.FilePrefix+r.Name) + ".erl"
}
