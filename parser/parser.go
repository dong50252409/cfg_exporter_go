package parser

import (
	"cfg_exporter/entities"
	"cfg_exporter/entities/decorator"
	"cfg_exporter/entities/typesystem"
	"cfg_exporter/interfaces"
	"cfg_exporter/reader"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	splitMultiConsRegexp = regexp.MustCompile(`\r\n|\r|\n`)
)

func FromFile(path string) (*entities.Table, error) {
	if ok := reader.CheckSupport(path); !ok {
		return nil, fmt.Errorf("配置表不支持！ 文件路径:%s", path)
	}

	records, err := reader.Read(path)
	if err != nil {
		return nil, fmt.Errorf("配置表读取失败！ 文件路径:%s %s", path, err)
	}

	table := &entities.Table{
		Path:     filepath.Dir(path),
		Filename: filepath.Base(path),
		Name:     strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
	}
	err = parser(table, records)
	if err != nil {
		return nil, fmt.Errorf("配置表解析失败！ 文件路径:%s %s", path, err)
	}

	return table, nil
}

// parser 解析配表
func parser(tbl *entities.Table, records [][]string) error {
	var fields []*entities.Field

	// 获取字段名
	fieldNameRow := records[3]
	// 获取字段装饰器
	decoratorRow := records[2]
	for column, colIndex := 0, 0; column < len(fieldNameRow); column++ {
		val := fieldNameRow[column]
		field := &entities.Field{Column: column, ColIndex: colIndex, Name: val}
		parseFiledName(field, val)
		err := parseFieldDecorator(tbl, field, decoratorRow[column])
		if err != nil {
			return err
		}
		val = strings.TrimSpace(val)
		if val != "" {
			fields = append(fields)
			colIndex++
		}
	}

	// 获取字段类型
	typeRow := records[1]
	// 获取字段注释
	CommentRow := records[0]
	// 获取主体数据
	recordRows := records[5:]
	// 实际数据
	dataSet := make([]any, 0, len(recordRows))
	for i := 0; i < len(recordRows); i++ {
		dataSet = append(dataSet, make([]any, 0, len(fields)))
	}

	for _, field := range fields {
		err := parseFieldType(field, typeRow[field.Column])
		if err != nil {
			return err
		}

		if field.Column < len(CommentRow) {
			parseFieldComment(field, CommentRow[field.Column])
		}

		err = parseRow(field, dataSet, recordRows)
		if err != nil {
			return err
		}

	}
	tbl.Fields = fields
	tbl.DataSet = dataSet

	runDecorator(tbl)
	return nil
}

func parseFiledName(field *entities.Field, val string) {
	val = strings.TrimSpace(val)
	if val != "" {
		field.Name = val
	}
}

// 解析字段类型
func parseFieldType(field *entities.Field, val string) error {
	val = strings.TrimSpace(val)
	if val != "" {
		t, err := typesystem.New(val)
		if err != nil {
			return fmt.Errorf("字段：%s 类型：%s %s", field.Name, val, err)
		}
		field.Type = t
	} else {
		return fmt.Errorf("字段：%s 类型不能为空", field.Name)
	}
	return nil
}

// 解析字段注释
func parseFieldComment(field *entities.Field, val string) {
	val = strings.TrimSpace(val)
	if val != "" {
		field.Comment = val
	}
}

// 解析字段装饰器
func parseFieldDecorator(tbl *entities.Table, field *entities.Field, val string) error {
	val = strings.TrimSpace(val)
	if val != "" {
		parts := splitMultiConsRegexp.Split(val, -1)
		for _, part := range parts {
			err := decorator.New(tbl, field, part)
			if err != nil {
				return fmt.Errorf("字段：%s 装饰器：%s %s", field.Name, part, err)
			}
		}
	}
	return nil
}

func parseRow(field *entities.Field, dataSet []any, records [][]string) error {
	var value any
	var err error
	for rowIndex, recordRows := range records {
		if field.Column >= len(recordRows) {
			value = nil
		} else {
			ceil := recordRows[field.Column]
			value, err = field.Type.(interfaces.ITypeSystem).ParseString(ceil)
			if err != nil {
				return fmt.Errorf("字段:%s 行:%d 列：%d 错误:%s", field.Name, rowIndex+5, field.Column+1, err)
			}
		}
		rows := dataSet[rowIndex].([]any)
		rows = append(rows, value)
		dataSet[rowIndex] = rows
	}
	return nil
}

func runDecorator(tbl *entities.Table) error {
	for _, field := range tbl.Fields {
		for k, d := range field.Decorators {
			err := d.(interfaces.IFieldDecorator).RunFieldDecorator(tbl, field)
			if err != nil {
				return fmt.Errorf("字段：%s 装饰器：%s 列：%d %s", field.Name, k, field.Column+1, err)
			}
		}
	}
	return nil
}
