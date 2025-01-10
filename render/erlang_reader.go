package render

import (
	"cfg_exporter/config"
	"cfg_exporter/model"
	"cfg_exporter/typesystem"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type ErlangRender struct {
	*renderTable
}

func init() {
	register("erlang", newErlangRender)
}

func newErlangRender(table renderTable) model.Render {
	return &ErlangRender{renderTable: &table}
}

const hrlHeadTemplate = `
{{- define "head" -}}
%% Auto Create, Don't Edit
-ifndef({{ .TableName | toUpper }}_HRL).
-define({{ .TableName | toUpper }}_HRL, true).
{{- end -}}
`

const hrlRecordTemplate = `
{{- define "record" -}}
-record({{.TableName}}, {
	{{- $lastIndex := len .Table.Fields | sub 1 }}
	{{- range $index, $field := .Fields }}
	{{ $field.Name }} = {{ defaultValue $field}} :: {{ fieldType $field.Type }}{{if lt $index $lastIndex}},{{end}}    % {{ $field.Comment }}
	{{- end }}
}).
{{- end -}}
`

const hrlMacroTemplate = `
{{- define "macro" -}}
{{- range _, $macroField := .Table.GetMacroFields }}
{{- range $rowIndex, $ds := .Table.DataSet }}
{{ $macroType := $macroField.Type.(*typesystem.Macro) }}
-define({{ $ds[$macroField.ColIndex] | toUpper}}, {{ $ds[.Table.GetFieldByName($macroType.ValueFieldName).ColIndex] }}).    % {{ $ds[.Table.GetFieldByName($macroType.ContentFieldName).ColIndex] }}).
{{- end }}
{{- end }}
{{- end -}}
`

const hrlTailTemplate = `
{{- define "tail" -}}
-endif.
{{- end -}}
`

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

func (r *ErlangRender) Execute() error {
	err := r.renderHrl()
	if err != nil {
		return err
	}

	err = r.renderErl()
	if err != nil {
		return err
	}

	return nil
}

func (r *ErlangRender) ExportDir() string {
	erlang := config.Config.Schema["erlang"]
	return filepath.Join(erlang.Destination, erlang.RecordPrefix)
}

// 渲染hrl
func (r *ErlangRender) renderHrl() error {
	hrlDir := filepath.Join(r.ExportDir(), "hrl")
	if err := os.MkdirAll(hrlDir, os.ModePerm); err != nil {
		return err
	}
	hrlFilepath := filepath.Join(hrlDir, r.Name+".hrl")
	hrlFile, err := os.Create(hrlFilepath)
	if err != nil {
		return err
	}
	defer func() { _ = hrlFile.Close() }()

	// 必备数据
	data := map[string]any{"Table": r}

	// 用到的函数
	funcMap := template.FuncMap{
		"toUpper":      strings.ToUpper,
		"toLower":      strings.ToLower,
		"defaultValue": r.DefaultValue,
		"fieldType":    r.FieldType,
		"sub":          func(a, b int) int { return a - b },
	}

	// 解析模板字符串
	tmpl := template.New("hrl").Funcs(funcMap)

	for _, tmplStr := range []string{hrlHeadTemplate, hrlRecordTemplate, hrlMacroTemplate, hrlTailTemplate, hrlTemplate} {
		tmpl, err = tmpl.Parse(tmplStr)
		if err != nil {
			return err
		}
	}

	// 执行模板渲染并输出到文件
	err = tmpl.Execute(hrlFile, data)
	if err != nil {
		return err
	}
	return nil
}

func (r *ErlangRender) renderErl() error {
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
		"toUpper":      strings.ToUpper,
		"defaultValue": r.DefaultValue,
		"fieldType":    r.FieldType,
		"sub":          func(a, b int) int { return a - b },
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

func (r *ErlangRender) DefaultValue(field model.Field) any {
	switch field.Type.(type) {
	case *typesystem.Integer:
		return 0
	case *typesystem.Float:
		return 0
	case *typesystem.Boolean:
		return "false"
	case *typesystem.String:
		return `<<""/utf8>>`
	case *typesystem.Lang:
		return `<<""/utf8>>`
	case *typesystem.Tuple:
		return "{}"
	case *typesystem.List:
		return "[]"
	case *typesystem.Map:
		return "#{}"
	case *typesystem.Macro:
		return `<<""/utf8>>`
	default:
		return "undefined"
	}
}

func (r *ErlangRender) FieldType(val model.TypeSystem) string {
	switch val.(type) {
	case *typesystem.Integer:
		return "integer()"
	case *typesystem.Float:
		return "float()"
	case *typesystem.Boolean:
		return "boolean()"
	case *typesystem.String:
		return "binary()"
	case *typesystem.Lang:
		return "binary()"
	case *typesystem.Tuple:
		return "tuple()"
	case *typesystem.List:
		return "list()"
	case *typesystem.Map:
		return "map()"
	case *typesystem.Macro:
		return "binary()"
	default:
		return "term()"
	}
}
