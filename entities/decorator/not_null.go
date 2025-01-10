package decorator

import "cfg_exporter/entities"

// NotNull 非空
type NotNull struct {
}

func init() {
	register("not_null", newNotNull)
}

func newNotNull(_ *entities.Table, field *entities.Field, _ string) error {
	field.Decorators["not_null"] = &NotNull{}
	return nil
}

func (nn *NotNull) Check() bool {
	return true
}
