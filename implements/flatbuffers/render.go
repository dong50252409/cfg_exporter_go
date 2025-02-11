package flatbuffers

import (
	"cfg_exporter/entities"
	"cfg_exporter/implements/json"
	"cfg_exporter/render"
	"fmt"
	"github.com/stoewer/go-strcase"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

type FBSRender struct {
	*render.Render
}

const dataSetTemplate = `
{{- define "tableSchema" -}}
{{- $pkFields := .Table.GetPrimaryKeyFields -}}
{{- $pkField := index $pkFields 0 -}}
table {{ .Table.Name | toUpperCamelCase }}{
    {{- range $_, $field := .Table.Fields }}
    {{ $field.Name | toLowerCamelCase }}: {{ $field.Type }}{{ if eq $field.Name $pkField.Name }}(key){{ end }};
	{{- end }}
}
{{- end -}}
`

const tailTemplate = `
{{- define "rootType" -}}
table {{ .Table.ConfigName }}List{
    {{ .Table.Name | toLowerCamelCase }}List: [{{ .Table.Name | toUpperCamelCase }}];
}

root_type {{ .Table.ConfigName }}List;
{{- end -}}
`

const fbTemplate = `
{{- if ne .Table.Schema.Namespace "" -}}
{{"namespace" }} {{ .Table.Namespace }};

{{ template "tableSchema" . }}
{{ else }}
{{- template "tableSchema" . }}
{{ end }}

{{ template "rootType" . }}
`

func init() {
	render.Register("flatbuffers", newtsRender)
}

func newtsRender(render *render.Render) render.IRender {
	return &FBSRender{render}
}

func (r *FBSRender) Execute() error {
	if err := r.Render.ExecuteBefore(); err != nil {
		return err
	}

	dir := r.ExportDir()
	fp := filepath.Join(dir, r.Filename())
	fileIO, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer func() { _ = fileIO.Close() }()

	// 必备数据
	data := map[string]any{"Table": r}

	// 解析模板字符串
	tmpl := template.New("fbs").Funcs(entities.FuncMap)

	for _, tmplStr := range []string{dataSetTemplate, tailTemplate, fbTemplate} {
		if tmpl, err = tmpl.Parse(tmplStr); err != nil {
			return err
		}
	}

	// 执行模板渲染并输出到文件
	if err = tmpl.Execute(fileIO, data); err != nil {
		return err
	}

	jsonRender, err := render.NewRender("json", r.Table)
	if err != nil {
		return err
	}
	jRender := jsonRender.(*json.JSONRender)
	if err = jRender.Execute(); err != nil {
		return err
	}

	fbFilename := filepath.Join(dir, r.Filename())
	jsonFilename := filepath.Join(jRender.ExportDir(), jRender.Filename())
	cmd := exec.Command(r.Schema.Flatc, "--no-warnings", "--unknown-json", "-o", dir, "-b", fbFilename, jsonFilename)
	if _, err = cmd.Output(); err != nil {
		return fmt.Errorf("error:%s", err)
	}

	return nil
}

func (r *FBSRender) Verify() error {
	return nil
}

func (r *FBSRender) ConfigName() string {
	return strcase.UpperCamelCase(r.Schema.TableNamePrefix + r.Name)
}

func (r *FBSRender) Filename() string {
	return strcase.SnakeCase(r.Schema.FilePrefix+r.Name) + ".fbs"
}

func (r *FBSRender) Namespace() string {
	return r.Schema.Namespace
}
