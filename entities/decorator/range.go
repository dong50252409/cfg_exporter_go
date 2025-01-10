package decorator

import (
	"cfg_exporter/entities"
	"cfg_exporter/util"
	"fmt"
	"strconv"
)

// Range 范围
type Range struct {
	minValue int
	maxValue int
}

func init() {
	registry["range"] = newRange
}

func newRange(_ *entities.Table, field *entities.Field, str string) error {
	if args := util.SubArgs(str, ","); args != nil && len(args) == 2 {
		v1, err1 := strconv.Atoi(args[0])
		v2, err2 := strconv.Atoi(args[1])
		if err1 == nil && err2 == nil && v1 <= v2 {
			field.Decorators["range"] = &Range{minValue: v1, maxValue: v2}
			return nil
		}
	}
	return fmt.Errorf("参数格式错误 range(最小值,最大值)")
}

func (r *Range) Check() bool {
	return true
}
