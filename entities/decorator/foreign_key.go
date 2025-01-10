package decorator

import (
	"cfg_exporter/entities"
	"cfg_exporter/util"
	"fmt"
	"strings"
)

// ForeignKey 外键引用
type ForeignKey struct {
	tableName string
	fieldName string
}

func init() {
	register("f_key", newForeignKey)
}

func newForeignKey(_ *entities.Table, field *entities.Field, str string) error {
	if args := util.SubArgs(str, ","); args != nil {
		l := strings.Split(str, ",")
		if len(l) == 2 {
			field.Decorators["f_key"] = &ForeignKey{tableName: l[0], fieldName: l[1]}
			return nil
		}
	}

	return fmt.Errorf("参数格式错误 f_key(表名,字段名)")
}

func (*ForeignKey) Check() bool {
	return true
}
