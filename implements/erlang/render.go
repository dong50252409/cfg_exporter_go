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

type erlangRender struct {
	*render.Render
	hrlRender *hrlRender
	erlRender *erlRender
}

func init() {
	render.Register("erlang", newErlangRender)
}

func newErlangRender(render *render.Render) render.IRender {
	r := &erlangRender{render, &hrlRender{}, &erlRender{}}
	r.hrlRender.erlangRender = r
	r.erlRender.erlangRender = r
	return r
}

// Execute 执行导出
func (r *erlangRender) Execute() error {
	if err := r.hrlRender.Execute(); err != nil {
		return err
	}
	if err := r.erlRender.Execute(); err != nil {
		return err
	}
	return nil
}

// Verify 验证导出结果
func (r *erlangRender) Verify() error {
	hrlDir := r.hrlRender.ExportDir()
	erl := filepath.Join(r.erlRender.ExportDir(), r.erlRender.Filename())
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
func (r *erlangRender) ConfigName() string {
	return strcase.SnakeCase(r.Schema.TableNamePrefix + r.Name)
}
