package erlang

import (
	"cfg_exporter/config"
	"github.com/stoewer/go-strcase"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type erlRender struct {
	*erlangRender
}

const erlHeadTemplate = `
{{- define "head" -}}
%% Auto Create, Don't Edit
-module({{ .Table.ConfigName | toLower }}).
-include("{{ .Table.ConfigName | toLower }}.hrl").
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
{{- $dsLastIndex := index $dataSet 0 | len | sub 1 -}}
{{- $pkValuesList := .Table.GetPrimaryKeyValuesByString -}}

{{- range $rowIndex, $dataRow := $dataSet -}}
get({{ index $pkValuesList $rowIndex | joinByComma }})->
    #{{ $configName }}{{ "{" }}
        {{- range $fieldIndex, $field := $fields }}
        {{ $field.Name }} = {{ index $dataRow $field.ColIndex | $field.Convert }}{{ if lt $fieldIndex $dsLastIndex }},{{ end }}
        {{- end }}
    {{ "};" }}
{{ end -}}
get(_ID0) ->
    throw({config_error}, ?MODULE, _ID0).
{{- end -}}
`
const erlListTemplate = `
{{- define "list" -}}
{{/* 声明模板渲染所需的变量 */}}
{{- $pkValuesList := .Table.GetPrimaryKeyValuesByString -}}
{{- $pkLastIndex := len $pkValuesList | sub 1 -}}

list() ->
    [
	{{- range $pkIndex, $pkValues := $pkValuesList -}}
	{{ "{" }}{{ $pkValues | joinByComma }}{{ "}" }}{{ if lt $pkIndex $pkLastIndex }}, {{ end }}
	{{- end -}}
    ].
{{- end -}}
`

const erlTemplate = `
{{- template "head" .}}

{{ template "get" .}}

{{ template "list" .}}
`

func (r *erlRender) Execute() error {
	erlDir := r.ExportDir()
	if err := os.MkdirAll(erlDir, os.ModePerm); err != nil {
		return err
	}

	erlFilepath := filepath.Join(erlDir, r.Filename())
	erlFile, err := os.Create(erlFilepath)
	if err != nil {
		return err
	}
	defer func() { _ = erlFile.Close() }()

	// 必备数据
	data := map[string]any{"Table": r}

	// 用到的函数
	funcMap := template.FuncMap{
		"toUpper":     strings.ToUpper,
		"toLower":     strings.ToLower,
		"joinByComma": func(items []string) string { return strings.Join(items, ", ") },
		"sub":         func(a, b int) int { return b - a },
	}

	// 解析模板字符串
	tmpl := template.New("erl").Funcs(funcMap)

	for _, tmplStr := range []string{erlHeadTemplate, erlGetTemplate, erlListTemplate, erlTemplate} {
		//for _, tmplStr := range []string{erlTemplate} {
		tmpl, err = tmpl.Parse(tmplStr)
		if err != nil {
			return err
		}
	}

	// 执行模板渲染并输出到文件
	err = tmpl.Execute(erlFile, data)
	if err != nil {
		return err
	}
	return nil
}

func (r *erlRender) ExportDir() string {
	return strcase.SnakeCase(filepath.Join(r.erlangRender.ExportDir(), "erl"))
}

func (r *erlRender) Filename() string {
	return config.Config.Schema["erlang"].FilePrefix + r.Name + ".erl"
}
