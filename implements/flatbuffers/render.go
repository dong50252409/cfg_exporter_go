package flatbuffers

import (
	"cfg_exporter/render"
	"fmt"
	"os/exec"
	"path/filepath"
)

type flatbuffersRender struct {
	*render.Render
}

func init() {
	render.Register("flatbuffers", newtsRender)
}

func newtsRender(render *render.Render) render.IRender {
	return &flatbuffersRender{render}
}

func (r *flatbuffersRender) Execute() error {
	fb := &fbsReader{r}
	if err := fb.Execute(); err != nil {
		return err
	}

	json := &jsonRender{r}
	if err := json.Execute(); err != nil {
		return err
	}

	dir := r.ExportDir()
	fbFilename := filepath.Join(dir, fb.Filename())
	jsonFilename := filepath.Join(dir, json.Filename())
	cmd := exec.Command(r.Schema.Flatc, "--no-warnings", "-o", dir, "-b", fbFilename, jsonFilename)
	// 获取命令的输出
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error:%s", err)
	}
	return nil
}

func (r *flatbuffersRender) ConfigName() string {
	return r.Schema.TableNamePrefix + r.Name
}
