package typescript

import (
	"cfg_exporter/entities"
	"cfg_exporter/parser"
	"cfg_exporter/render"
	"fmt"
	"github.com/stoewer/go-strcase"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const headTemplate = `
{{- define "head" -}}
import * as flatbuffers from '../../../../Plugins/FlatBuffers';
import { BaseData_Fbs } from '../../BaseData';
import { cacheObjRes, cacheStrRes, FbsData, FbsDataList } from '../../FbsDataListView';
import { {{ .Table.ConfigName | toUpperCamelCase }} } from '../FbsCls/{{ .Table.Name | toKebabCase }}';
import { {{ .Table.ConfigName | toUpperCamelCase }}List } from '../FbsCls/{{ .Table.Name | toKebabCase }}-list';
{{- end -}}
`

const baseClassTemplate = `
{{- define "base_class" -}}
export class {{ .Table.Name | toUpperCamelCase }} extends BaseData_Fbs<Cfg{{ .Table.Name | toUpperCamelCase }}> {
    public getFbsDataList(data: ArrayBuffer): FbsDataList<{{ .Table.InnerConfigName }}> {
        return new {{ .Table.InnerConfigName }}List(data);
    }

	{{- $pkFields := .Table.GetPrimaryKeyFields }}
	{{ $pkLen := $pkFields | len | add -1 }}
    get({{ range $pkIndex, $pkField := $pkFields }}{{ $pkField.Name | toLowerCamelCase }}: {{ $pkField.Type }}{{ if lt $pkIndex $pkLen }}, {{ end }}{{ end }}): {{ .Table.InnerConfigName }} {
	    return super.get({{ range $pkIndex, $pkField := $pkFields }}{{ $pkField.Name | toLowerCamelCase }}{{ if lt $pkIndex $pkLen }}, {{ end }}{{ end }});
    }
}
{{- end -}}
`

const innerClass1Template = `
{{- define "inner_class1" -}}
export class {{ .Table.InnerConfigName }} extends FbsData {
    private _fbs: {{ .Table.ConfigName | toUpperCamelCase }};

    __init(fbs: {{ .Table.ConfigName | toUpperCamelCase }}) {
		this._fbs = fbs;
    }

    public clone(): {{ .Table.InnerConfigName }} {
	    let newFD: {{ .Table.InnerConfigName }} = new {{ .Table.InnerConfigName }}();
	    newFD._fbs = new {{ .Table.ConfigName | toUpperCamelCase }}();
        newFD._fbs.bb = this._fbs.bb;
        newFD._fbs.bb_pos = this._fbs.bb_pos;
	    return newFD;
    }

	{{- $pkFields := .Table.GetPrimaryKeyFields }}
	{{ $pkLen := $pkFields | len | add -1 }}
    public checkKeyPrefix(): string{
        return {{ "\x60" }}{{ .Table.InnerConfigName }}.{{ range $pkIndex, $pkField := $pkFields }}${this.{{ $pkField.Name | toLowerCamelCase }}}{{ if lt $pkIndex $pkLen }}.{{ end }}{{ end }}{{ "\x60" }};
    }

    {{ range $index, $field := .Table.Fields }}
	{{ $field.Type.Decorator }}
	public get {{ $field.Name | toLowerCamelCase }}(): {{ $field.Type }} {
		return this._fbs.{{ $field.Name | toLowerCamelCase }}();
	}
    {{ end }}
}
{{- end -}}
`
const innerClass2Template = `
{{- define "inner_class2" -}}
export class {{ .Table.InnerConfigName }}List extends FbsDataList<{{ .Table.InnerConfigName }}> {
	private _fbsList: {{ .Table.ConfigName | toUpperCamelCase }}List;

	__init(data: ArrayBuffer) {
		this._fbsList = {{ .Table.ConfigName | toUpperCamelCase }}List.getRootAs{{ .Table.ConfigName | toUpperCamelCase }}List(new flatbuffers.ByteBuffer(new Uint8Array(data)));
	}

    public get length(): number {
        return this._fbsList.{{ .Table.Name | toLowerCamelCase }}Length();
    }

    public getFbsData(index: number, obj?: {{ .Table.InnerConfigName | toUpperCamelCase }}): {{ .Table.InnerConfigName | toUpperCamelCase }} {
        obj = obj ? obj : new {{ .Table.InnerConfigName | toUpperCamelCase }}();
        obj.__init(this._fbsList.{{ .Table.Name | toLowerCamelCase }}(index) as {{ .Table.ConfigName | toUpperCamelCase}});
        return obj;
    }
}
{{- end -}}
`

const enumTemplate = `
{{- define "enum" -}}
{{ range $_, $macro := .Table.GetMacroDecorators }}
export enum {{ $macro.MacroName | toUpperCamelCase }}Enum { 
	{{- range $_, $macroDetail := $macro.List }}
    /** {{ $macroDetail.Comment }} */
    {{ $macroDetail.Key | toUpperSnakeCase }} = {{ $macroDetail.Value }},
	{{- end }}
}
{{ end }}
{{- end -}}
`

const tsTemplate = `
{{- template "head" .}}

{{ template "base_class" .}}

{{ template "inner_class1" .}}

{{ template "inner_class2" .}}

{{ template "enum" . -}}
`

type TSRender struct {
	*render.Render
}

func init() {
	render.Register("typescript", newtsRender)
}

func newtsRender(render *render.Render) render.IRender {
	return &TSRender{render}
}

func (r *TSRender) Execute() error {
	if err := r.Render.ExecuteBefore(); err != nil {
		return err
	}

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
	tmpl := template.New("typescript").Funcs(entities.FuncMap)

	for _, tmplStr := range []string{headTemplate, baseClassTemplate, innerClass1Template, innerClass2Template, enumTemplate, tsTemplate} {
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

	fbsParser, err := parser.CloneParser("flatbuffers", r.Table)
	if err != nil {
		return err
	}
	fbsRender, err := render.NewRender("flatbuffers", fbsParser.GetTable())
	if err != nil {
		return err
	}
	if err = fbsRender.Execute(); err != nil {
		return err
	}

	fbFilename := filepath.Join(dir, r.Filename())
	cmd := exec.Command(r.Schema.Flatc, "--ts ", "--no-warnings", "--ts-omit-entrypoint", "-o", dir, fbFilename)
	if _, err = cmd.Output(); err != nil {
		return fmt.Errorf("error:%s", err)
	}

	return nil
}

func (r *TSRender) Verify() error {
	return nil
}

func (r *TSRender) Filename() string {
	return strcase.KebabCase(r.Schema.FilePrefix+r.Name) + ".ts"
}

func (r *TSRender) ConfigName() string {
	return r.Schema.TableNamePrefix + r.Name
}

func (r *TSRender) InnerConfigName() string {
	return strcase.UpperCamelCase("cfg_" + r.Name)
}
