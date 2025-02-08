package erlang

import (
	"cfg_exporter/render"
	"fmt"
	"github.com/stoewer/go-strcase"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type ERLRender struct {
	*render.Render
	hrlRender *HRender
	erlRender *ERender
}

func init() {
	render.Register("erlang", newErlangRender)
}

func newErlangRender(render *render.Render) render.IRender {
	r := &ERLRender{render, &HRender{}, &ERender{}}
	r.hrlRender.ERLRender = r
	r.erlRender.ERLRender = r
	return r
}

// Execute 执行导出
func (r *ERLRender) Execute() error {
	if err := r.Render.ExecuteBefore(); err != nil {
		return err
	}
	if err := r.hrlRender.Execute(); err != nil {
		return err
	}
	if err := r.erlRender.Execute(); err != nil {
		return err
	}
	return nil
}

// Verify 验证导出结果
func (r *ERLRender) Verify() error {
	hrlDir := r.hrlRender.ExportDir()
	erl := filepath.Join(r.erlRender.ExportDir(), r.erlRender.Filename())
	fmt.Printf("开始验证生成结果：%s\n", erl)
	var out string
	switch runtime.GOOS {
	case "windows":
		out = os.Getenv("TEMP")
	default:
		out = "/dev/null"
	}
	cmd := exec.Command("erlc", "-Werror", "-Wall", "-o", out, "-I", hrlDir, erl)
	result, err := cmd.Output()
	if err != nil {
		fmt.Printf("%s", result)
		return err
	}
	return nil
}

// ConfigName 配置名
func (r *ERLRender) ConfigName() string {
	return strcase.SnakeCase(r.Schema.TableNamePrefix + r.Name)
}
