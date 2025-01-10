package decorator

import "cfg_exporter/entities"

type Unique struct {
}

func init() {
	register("u_key", newUnique)
}

func newUnique(_ *entities.Table, field *entities.Field, _ string) error {
	field.Decorators["u_key"] = &Unique{}
	return nil
}

func (u *Unique) Check() bool {
	return true
}
