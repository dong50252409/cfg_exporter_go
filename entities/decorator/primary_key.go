package decorator

import "cfg_exporter/entities"

// PrimaryKey 主键
type PrimaryKey struct {
	columnIndies []int
}

func init() {
	register("p_key", newPrimaryKey)
}

func newPrimaryKey(tbl *entities.Table, field *entities.Field, _ string) error {
	if v, ok := tbl.Decorators["p_key"]; ok {
		pk := v.(*PrimaryKey)
		pk.columnIndies = append(pk.columnIndies, field.ColIndex)
	}
	return nil
}

func (*PrimaryKey) Check() bool {
	return true
}
