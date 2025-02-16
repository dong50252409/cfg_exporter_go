package erlang

import (
	"cfg_exporter/entities"
	"fmt"
	"github.com/stoewer/go-strcase"
	"os"
	"path/filepath"
	"text/template"
)

type HRender struct {
	*ERLRender
}

const hrlHeadTemplate = `
{{- define "head" -}}
%% Auto Create, Don't Edit
-ifndef({{ .Table.ConfigName | toUpperSnakeCase }}_HRL).
-define({{ .Table.ConfigName | toUpperSnakeCase }}_HRL, true).
{{- end -}}
`

const hrlRecordTemplate = `
{{- define "record" -}}
-record({{ .Table.ConfigName | toSnakeCase }}, {
    {{- $lastIndex := len .Table.Fields | add -1 }}
    {{- range $index, $field := .Table.Fields }}
    {{ $field.Name | toSnakeCase }} = {{ $field.DefaultValue }} :: {{ $field.Type }}{{ if lt $index $lastIndex }},{{ end }}	% {{ $field.Comment }}
    {{- end }}
}).
{{- end -}}
`

const hrlMacroTemplate = `
{{- define "macro" -}}
{{- range $_, $macro := .Table.GetMacroDecorators -}}
%% {{ $macro.MacroName }}
{{ range $_, $macroDetail := $macro.List -}}
-define({{ $macroDetail.Key | toUpperSnakeCase }}, {{ $macroDetail.Value }}).	  % {{ $macroDetail.Comment }}
{{ end }}
{{ end -}}
{{- end -}}
`

const hrlTailTemplate = `
{{- define "tail" -}}
-endif.
{{- end -}}
`

const hrlTemplate = `
{{- template "head" .}}

{{ template "record" .}}

{{ template "macro" . -}}

{{ template "tail" .}}
`

func (r *HRender) Execute() error {
	hrlDir := r.ExportDir()
	if err := os.MkdirAll(hrlDir, os.ModePerm); err != nil {
		return err
	}

	fp := filepath.Join(hrlDir, r.Filename())
	fileIO, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer func() { _ = fileIO.Close() }()

	// 必备数据
	data := map[string]any{"Table": r}

	// 解析模板字符串
	tmpl := template.New("hrl").Funcs(entities.FuncMap)

	for _, tmplStr := range []string{hrlHeadTemplate, hrlRecordTemplate, hrlMacroTemplate, hrlTailTemplate, hrlTemplate} {
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

	fmt.Printf("导出配置：%s\n", fp)

	return nil
}

func (r *HRender) ExportDir() string {
	return strcase.SnakeCase(filepath.Join(r.ERLRender.ExportDir(), "hrl"))
}

func (r *HRender) Filename() string {
	return strcase.SnakeCase(r.Schema.FilePrefix+r.Name) + ".hrl"
}
