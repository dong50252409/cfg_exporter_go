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

type ErlangHrlRender struct {
	*entities.Table
}

func init() {
	Register("erlang", newErlangHrlRender)
}

func newErlangHrlRender(table *entities.Table) interfaces.IRender {
	return &ErlangHrlRender{table}
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

func (r *ErlangHrlRender) Execute() error {
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
		"toUpper": strings.ToUpper,
		"toLower": strings.ToLower,
		"sub":     func(a, b int) int { return a - b },
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

func (r *ErlangHrlRender) ExportDir() string {
	erlang := config.Config.Schema["erlang"]
	return filepath.Join(erlang.Destination, erlang.RecordPrefix)
}
