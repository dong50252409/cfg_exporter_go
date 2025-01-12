package decorator

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"cfg_exporter/entities/typesystem"
	"cfg_exporter/util"
	"fmt"
	"strconv"
)

// Range 范围
type Range struct {
	minValue any
	maxValue any
}

func init() {
	registry["range"] = newRange
}

func newRange(_ *entities.Table, field *entities.Field, str string) error {
	args := util.SubArgs(str, ",")
	if len(args) == 2 {
		switch field.Type.(type) {
		case *typesystem.Integer:
			v1, err1 := strconv.Atoi(args[0])
			v2, err2 := strconv.Atoi(args[1])
			if err1 == nil && err2 == nil && v1 <= v2 {
				field.Decorators["range"] = &Range{minValue: v1, maxValue: v2}
				return nil
			}
		case *typesystem.Float:
			v1, err1 := strconv.ParseFloat(args[0], 64)
			v2, err2 := strconv.ParseFloat(args[1], 64)
			if err1 == nil && err2 == nil && v1 <= v2 {
				field.Decorators["range"] = &Range{minValue: v1, maxValue: v2}
				return nil
			}
		default:
			return fmt.Errorf("类型无法使用此装饰器")
		}
	}
	return fmt.Errorf("参数格式错误 range(最小值,最大值)")
}

func (r *Range) RunFieldDecorator(tbl *entities.Table, field *entities.Field) error {
	for corIndex, row := range tbl.DataSet {
		err := r.Equal(corIndex, row[field.Column], field.Type)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Range) Equal(rowIndex int, v any, t any) error {
	switch t.(type) {
	case *typesystem.Integer:
		if !(r.minValue.(int) >= v.(int) && r.maxValue.(int) <= v.(int)) {
			return fmt.Errorf("第 %d 行 数值必须在%d到%d之间", rowIndex+config.Config.BodyStartRow, r.minValue, r.maxValue)
		}
		return nil
	case *typesystem.Float:
		if !(r.minValue.(float64) >= v.(float64) && r.maxValue.(float64) <= v.(float64)) {
			return fmt.Errorf("第 %d 行 数值必须在%d到%d之间", rowIndex+config.Config.BodyStartRow, r.minValue, r.maxValue)
		}
		return nil
	default:
		return fmt.Errorf("类型无法使用此装饰器")
	}
}
