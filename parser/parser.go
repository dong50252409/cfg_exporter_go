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

// ParseFromFile 从文件读取解析配置表
func ParseFromFile(path string) (*entities.Table, error) {
	if ok := reader.CheckSupport(path); !ok {
		return nil, fmt.Errorf("配置表不支持！ 文件路径:%s", path)
	}

	records, err := reader.Read(path)
	if err != nil {
		return nil, fmt.Errorf("配置表读取失败！ 文件路径:%s %s", path, err)
	}

	table := &entities.Table{
		Path:       filepath.Dir(path),
		Filename:   filepath.Base(path),
		Name:       strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		Decorators: []any{},
		Records:    records,
	}
	err = Parse(table, records)
	if err != nil {
		return nil, fmt.Errorf("配置表解析失败！ 文件路径:%s %s", path, err)
	}

	return table, nil
}

// Parse 解析配表
func Parse(tbl *entities.Table, records [][]string) error {
	var fields []*entities.Field

	// 字段名 字段类型 字段装饰器 字段注释 主体数据
	fnRow, ftRow, fdRow, fcRow, recordRows := records[3], records[1], records[2], records[0], records[5:]

	// 实际数据
	dataSet := make([][]any, 0, len(recordRows))
	for i := 0; i < len(recordRows); i++ {
		dataSet = append(dataSet, make([]any, 0, len(fields)))
	}

	for column, colIndex := 0, 0; column < len(fnRow); column++ {
		val := strings.TrimSpace(fnRow[column])
		if val != "" {
			field := &entities.Field{Column: column, ColIndex: colIndex, Name: val, Decorators: make(map[string]any)}
			err := ParseFieldType(field, ftRow)
			if err != nil {
				return err
			}
			ParseFieldComment(field, fcRow)
			err = ParseRow(field, dataSet, recordRows)
			if err != nil {
				return err
			}

			fields = append(fields, field)
			colIndex++
		}
	}
	tbl.Fields = fields
	tbl.DataSet = dataSet

	for column := 0; column < len(fdRow); column++ {
		var field *entities.Field
		if column >= len(fnRow) {
			field = &entities.Field{Column: column, Decorators: make(map[string]any)}
		} else {
			val := strings.TrimSpace(fnRow[column])
			if val != "" {
				field = tbl.GetFieldByName(fnRow[column])
			} else {
				field = &entities.Field{Column: column, Decorators: make(map[string]any)}
			}
		}

		err := ParseFieldDecorator(tbl, field, fdRow, column)
		if err != nil {
			return err
		}
	}

	err := RunDecorator(tbl)
	if err != nil {
		return err
	}
	return nil
}

// ParseFieldType 解析字段类型
func ParseFieldType(field *entities.Field, ftRow []string) error {
	val := strings.TrimSpace(ftRow[field.Column])
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

// ParseFieldComment 解析字段注释
func ParseFieldComment(field *entities.Field, fcRow []string) {
	if field.Column < len(fcRow) {
		val := strings.TrimSpace(fcRow[field.Column])
		if val != "" {
			field.Comment = val
		}
	}
}

// ParseRow 解析行
func ParseRow(field *entities.Field, dataSet [][]any, records [][]string) error {
	var value any
	var err error
	for rowIndex, recordRows := range records {
		if field.Column >= len(recordRows) {
			value = nil
		} else {
			cell := recordRows[field.Column]
			value, err = field.Type.(interfaces.ITypeSystem).ParseString(cell)
			if err != nil {
				return fmt.Errorf("字段:%s 行:%d 列：%d 错误:%s", field.Name, rowIndex+5, field.Column+1, err)
			}
		}
		rows := dataSet[rowIndex]
		rows = append(rows, value)
		dataSet[rowIndex] = rows
	}
	return nil
}

// ParseFieldDecorator 解析字段装饰器
func ParseFieldDecorator(tbl *entities.Table, field *entities.Field, fdRow []string, column int) error {
	if column < len(fdRow) {
		val := strings.TrimSpace(fdRow[column])
		if val != "" {
			parts := splitMultiConsRegexp.Split(val, -1)
			for _, part := range parts {
				err := decorator.New(tbl, field, part)
				if err != nil {
					return fmt.Errorf("字段：%s 装饰器：%s %s", field.Name, part, err)
				}
			}
		}
	}
	return nil
}

// RunDecorator 运行装饰器
func RunDecorator(tbl *entities.Table) error {
	for _, d := range tbl.Decorators {
		err := d.(interfaces.ITableDecorator).RunTableDecorator(tbl)
		if err != nil {
			return fmt.Errorf("装饰器：%s %s", d.(interfaces.IDecorator).Name(), err)
		}
	}

	for _, field := range tbl.Fields {
		for _, d := range field.Decorators {
			err := d.(interfaces.IFieldDecorator).RunFieldDecorator(tbl, field)
			if err != nil {
				return fmt.Errorf("字段：%s 装饰器：%s 列：%d %s", field.Name, d.(interfaces.IDecorator).Name(), field.Column+1, err)
			}
		}
	}
	return nil
}
