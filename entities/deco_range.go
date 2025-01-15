package entities

import (
	"cfg_exporter/config"
	"cfg_exporter/util"
	"fmt"
	"reflect"
	"strconv"
)

// Range 范围
type Range struct {
	minValue any
	maxValue any
}

func init() {
	decoratorRegister("range", newRange)
}

func newRange(_ *Table, field *Field, str string) error {
	args := util.SubArgs(str, ",")
	if len(args) == 2 {
		switch field.Type.(ITypeSystem).GetKind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v1, err1 := strconv.ParseInt(args[0], 10, 64)
			v2, err2 := strconv.ParseInt(args[1], 10, 64)
			if err1 == nil && err2 == nil && v1 <= v2 {
				field.Decorators["range"] = &Range{minValue: v1, maxValue: v2}
				return nil
			}
		case reflect.Float32, reflect.Float64:
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

func (r *Range) RunFieldDecorator(tbl *Table, field *Field) error {
	for corIndex, row := range tbl.DataSet {
		v := row[field.Column]
		if v != nil {
			err := r.Equal(corIndex, row[field.Column], field.Type)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (*Range) Name() string {
	return "range"
}

func (r *Range) Equal(rowIndex int, v any, t any) error {
	switch t.(ITypeSystem).GetKind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !(r.minValue.(int64) <= v.(int64) && v.(int64) <= r.maxValue.(int64)) {
			return fmt.Errorf("第 %d 行 数值必须在%d到%d之间", rowIndex+config.Config.BodyStartRow, r.minValue, r.maxValue)
		}
		return nil
	case reflect.Float32, reflect.Float64:
		if !(r.minValue.(float64) <= v.(float64) && v.(float64) <= r.maxValue.(float64)) {
			return fmt.Errorf("第 %d 行 数值必须在%d到%d之间", rowIndex+config.Config.BodyStartRow, r.minValue, r.maxValue)
		}
		return nil
	default:
		return fmt.Errorf("类型无法使用此装饰器")
	}
}
