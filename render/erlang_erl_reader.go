package render

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"cfg_exporter/interfaces"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type ErlangErlRender struct {
	*entities.Table
}

func init() {
	Register("erlang", newErlangErlRender)
}

func newErlangErlRender(table *entities.Table) interfaces.IRender {
	return &ErlangErlRender{table}
}

const hrlTemplate = `
{{template "head" .}}

{{template "record" .}}

{{template "macro" .}}

{{template "tail" .}}
`

const erlHeadTemplate = `
{{- define "head" -}}
%% Auto Create, Don't Edit
-module({{ .TableName | toLower }}).
-include("{{ .TableName | toLower }}.hrl").
-compile(export_all).
-compile(nowarn_export_all).
-compile({no_auto_import,[get/1]}).
{{- end -}}
`

const erlGetTemplate = `
{{- define "getter" -}}
{{ $pkFields := .Table.GetPrimaryKeyFields }}
{{ $pkLastIndex := len $pkFields | sub 1 }}
{{ $dsLastIndex := len .Table.DataSet | sub 1 }}
{{ $fieldLastIndex := len .Table.Fields | sub 1 }}
{{- range $rowIndex, $ds := .Table.DataSet}}
{{- range _, $field := .Table.Fields}}
get({{- range $pkIndex, $pksField := $pkFields }}{{ $ds[$pksField.ColIndex] }}{{ if lt $pkIndex $pkLastIndex }},{{end}}}) ->
    #{{ .Table.Name | toLower}}{
        {{ $field.Name | toLower}} = {{ $ds[$field.ColIndex] | toLower}}{{ if lt $rowIndex $fieldLastIndex }},{{end}}
}{{ if lt $rowIndex $fieldLastIndex }};{{end}}
{{- end -}}
get(_ID0) ->
	throw({config_error}, ?MODULE, _ID0).
{{- end -}}
`

const erlListTemplate = `
{{- define "list" -}}
list() ->
	[{{- range $index, $field := .Fields }}
		{{ $field.Name }}{{if lt $index 1}},{{end}}
	{{- end }}].
{{- end -}}
`

const erlTemplate = `
{{template "head" .}}

{{template "getter" .}}

{{template "list" .}}
`

func (r *ErlangErlRender) Execute() error {
	if err := os.MkdirAll(r.ExportDir()+"/erl", os.ModePerm); err != nil {
		return err
	}

	erlFile, err := os.Create(r.ExportDir() + "/erl/" + r.Name + ".erl")
	if err != nil {
		return err
	}
	defer func() { _ = erlFile.Close() }()

	// 必备数据
	data := map[string]any{"Table": r}

	// 用到的函数
	funcMap := template.FuncMap{
		"toUpper": strings.ToUpper,
		"sub":     func(a, b int) int { return a - b },
	}

	// 解析模板字符串
	tmpl := template.New("erl").Funcs(funcMap)

	for _, tmplStr := range []string{erlHeadTemplate, erlGetTemplate, erlListTemplate, erlTemplate} {
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

func (r *ErlangErlRender) ExportDir() string {
	erlang := config.Config.Schema["erlang"]
	return filepath.Join(erlang.Destination, erlang.RecordPrefix)
}
