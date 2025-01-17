package erlang

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"cfg_exporter/render"
	"github.com/stoewer/go-strcase"
)

type erlangRender struct {
	*entities.Table
}

func init() {
	render.Register("erlang", newErlangRender)
}

func newErlangRender(table *entities.Table) render.IRender {
	return &erlangRender{table}
}

func (r *erlangRender) Execute() error {
	h := &hrlRender{r}
	err := h.Execute()
	if err != nil {
		return err
	}

	e := &erlRender{r}
	err = e.Execute()
	if err != nil {
		return err
	}

	return nil
}

func (r *erlangRender) ExportDir() string {
	erlang := config.Config.Schema["erlang"]
	return erlang.Destination
}

func (r *erlangRender) Filename() string {
	return ""
}

func (r *erlangRender) ConfigName() string {
	return strcase.SnakeCase(config.Config.Schema["erlang"].TableNamePrefix + r.Name)
}
