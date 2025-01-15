package erlang

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"github.com/stoewer/go-strcase"
	"os"
	"path/filepath"
	"text/template"
)

type hrlRender struct {
	*erlangRender
}

const hrlHeadTemplate = `
{{- define "head" -}}
%% Auto Create, Don't Edit
-ifndef({{ .Table.ConfigName | toUpper }}_HRL).
-define({{ .Table.ConfigName | toUpper }}_HRL, true).
{{- end -}}
`

const hrlRecordTemplate = `
{{- define "record" -}}
-record({{.Table.ConfigName}}, {
	{{- $lastIndex := len .Table.Fields | add -1 }}
	{{- range $index, $field := .Table.Fields }}
	{{ $field.Name }} = {{ $field.DefaultValue }} :: {{ $field.Type }}{{ if lt $index $lastIndex }},{{ end }}	% {{ $field.Comment }}
	{{- end }}
}).
{{- end -}}
`

const hrlMacroTemplate = `
{{- define "macro" -}}
{{- range $_, $macro := .Table.GetMacroDecorators -}}
%% {{ $macro.MacroName }}
{{ range $_, $macroDetail := $macro.List -}}
-define({{ $macroDetail.Key | toUpper}}, {{ $macroDetail.Value }}).	  % {{ $macroDetail.Comment }}
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

func (r *hrlRender) Execute() error {
	hrlDir := r.ExportDir()
	if err := os.MkdirAll(hrlDir, os.ModePerm); err != nil {
		return err
	}
	hrlFilepath := filepath.Join(hrlDir, r.Filename())
	fileIO, err := os.Create(hrlFilepath)
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
	return nil
}

func (r *hrlRender) ExportDir() string {
	return strcase.SnakeCase(filepath.Join(r.erlangRender.ExportDir(), "hrl"))
}

func (r *hrlRender) Filename() string {
	return config.Config.Schema["erlang"].FilePrefix + r.Name + ".hrl"
}
