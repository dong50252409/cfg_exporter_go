package decorator

import (
	"cfg_exporter/entities"
	"cfg_exporter/util"
	"fmt"
	"strings"
)

// ForeignKey 外键引用
type ForeignKey struct {
	TableName string
	FieldName string
}

func init() {
	register("f_key", newForeignKey)
}

func newForeignKey(_ *entities.Table, field *entities.Field, str string) error {
	args := util.SubArgs(str, ",")
	if len(args) == 2 {
		l := strings.Split(str, ",")
		if len(l) == 2 {
			field.Decorators["f_key"] = &ForeignKey{TableName: l[0], FieldName: l[1]}
			return nil
		}
	}
	return fmt.Errorf("参数格式错误 f_key(表名,字段名)")
}

func (f *ForeignKey) Name() string {
	return "f_key"
}

func (f *ForeignKey) RunFieldDecorator(tbl *entities.Table, field *entities.Field) error {
	// TODO 实现读取外键数据
	return nil
}
