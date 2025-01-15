package erlang

import (
	"cfg_exporter/config"
	"github.com/stoewer/go-strcase"
	"os"
	"path/filepath"
	"strings"
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
	{{- $lastIndex := len .Table.Fields | sub 1 }}
	{{- range $index, $field := .Table.Fields }}
	{{ $field.Name }} = {{ $field.DefaultValue }} :: {{ $field.Type }}{{ if lt $index $lastIndex }},{{ end }}	% {{ $field.Comment }}
	{{- end }}
}).
{{- end -}}
`

const hrlMacroTemplate = `
{{- define "macro" -}}
{{- range $_, $macro := .Table.GetMacroFields -}}
{{- range $_, $macroDetail := $macro.List -}}
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
	hrlFile, err := os.Create(hrlFilepath)
	if err != nil {
		return err
	}
	defer func() { _ = hrlFile.Close() }()

	// 必备数据
	data := map[string]any{"Table": r}

	// 用到的函数
	funcMap := template.FuncMap{
		"toUpper": strings.ToUpper,
		"toLower": strings.ToLower,
		"sub":     func(a, b int) int { return b - a },
	}

	// 解析模板字符串
	tmpl := template.New("hrl").Funcs(funcMap)

	for _, tmplStr := range []string{hrlHeadTemplate, hrlRecordTemplate, hrlMacroTemplate, hrlTailTemplate, hrlTemplate} {
		//for _, tmplStr := range []string{hrlHeadTemplate, hrlRecordTemplate, hrlMacroTemplate} {
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

func (r *hrlRender) ExportDir() string {
	return strcase.SnakeCase(filepath.Join(r.erlangRender.ExportDir(), "hrl"))
}

func (r *hrlRender) Filename() string {
	return config.Config.Schema["erlang"].FilePrefix + r.Name + ".hrl"
}
