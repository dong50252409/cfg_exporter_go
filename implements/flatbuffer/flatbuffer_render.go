package flatbuffer

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"cfg_exporter/render"
)

type flatbufferRender struct {
	*entities.Table
	schema config.Schema
}

func init() {
	render.Register("flatbuffer", newtsRender)
}

func newtsRender(table *entities.Table) render.IRender {
	return &flatbufferRender{table, config.Config.Schema["flatbuffer"]}
}

func (r *flatbufferRender) Execute() error {
	fb := &fbRender{r}
	if err := fb.Execute(); err != nil {
		return err
	}

	json := &jsonRender{r}
	if err := json.Execute(); err != nil {
		return err
	}
	return nil
}

func (r *flatbufferRender) ExportDir() string {
	return r.schema.Destination
}

func (r *flatbufferRender) Filename() string {
	return ""
}

func (r *flatbufferRender) ConfigName() string {
	return r.schema.TableNamePrefix + r.Name
}
