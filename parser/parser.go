package parser

import (
	"cfg_exporter/entities"
	"cfg_exporter/reader"
	"cfg_exporter/util"
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

type IParser interface {
	ParseFromFile(path string) (*entities.Table, error)
}

type Parser struct {
	*entities.Table
	// 字段名行号
	FieldNameRow int
	// 字段类型行号
	FieldTypeRow int
	// 字段注释行号
	FieldCommentRow int
	// 字段装饰器行号
	FieldDecoratorRow int
	// 主体数据开始行号
	BodyStartRow int
}

var (
	parserRegistry       = make(map[string]func(p *Parser) IParser)
	splitMultiConsRegexp = regexp.MustCompile(`\r\n|\r|\n`)
)

// RegisterParser 注册解析器
func RegisterParser(name string, cls func(p *Parser) IParser) {
	parserRegistry[name] = cls
}

// NewParser 新建解析器
func NewParser(schemaName string) (IParser, error) {
	if parser, ok := parserRegistry[schemaName]; ok {
		return parser(&Parser{}), nil
	}
	return nil, fmt.Errorf("未找到解析器：%s", schemaName)
}

// ParseFromFile 从文件读取解析配置表
func (p *Parser) ParseFromFile(path string) (*entities.Table, error) {
	if ok := reader.CheckSupport(path); !ok {
		return nil, fmt.Errorf("配置表不支持！文件路径:%s", path)
	}

	records, err := reader.Read(path)
	if err != nil {
		return nil, fmt.Errorf("配置表读取失败！ 文件路径:%s %s", path, err)
	}

	filename := filepath.Base(path)
	name, err := util.SubTableName(filename)
	if err != nil {
		return nil, err
	}

	fmt.Printf("开始解析配置：%s\n", filename)

	table := &entities.Table{
		Path:       filepath.Dir(path),
		Filename:   filename,
		Name:       name,
		Decorators: []any{},
		Records:    records,
	}
	p.Table = table

	err = p.Parse()
	if err != nil {
		return nil, fmt.Errorf("配置表解析失败！ 文件路径:%s %s", path, err)
	}
	return table, nil
}

// Parse 解析配表
func (p *Parser) Parse() error {
	tbl := p.Table

	// 字段名 字段类型 字段装饰器 字段注释 主体数据
	fnRow, ftRow, fdRow, fcRow, recordRows := tbl.Records[p.FieldNameRow-1], tbl.Records[p.FieldTypeRow-1],
		tbl.Records[p.FieldDecoratorRow-1], tbl.Records[p.FieldCommentRow-1], tbl.Records[p.BodyStartRow-1:]

	// 字段列表
	tbl.Fields = make([]*entities.Field, 0, len(fnRow))

	// 实际数据
	tbl.DataSet = make([][]any, 0, len(recordRows))
	for i := 0; i < len(recordRows); i++ {
		tbl.DataSet = append(tbl.DataSet, make([]any, 0, len(tbl.Fields)))
	}

	for column, colIndex := 0, 0; column < len(fnRow); column++ {
		val := strings.TrimSpace(fnRow[column])
		if val != "" {
			field := &entities.Field{Column: column, ColIndex: colIndex, Name: val, Decorators: make(map[string]any)}
			if err := p.ParseFieldType(field, ftRow); err != nil {
				return err
			}
			p.ParseFieldComment(field, fcRow)
			p.ParseFieldDefaultValue(field)
			if err := p.ParseRow(field, recordRows); err != nil {
				return err
			}

			tbl.Fields = append(tbl.Fields, field)
			colIndex++
		}
	}

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

		if err := p.ParseFieldDecorator(field, fdRow, column); err != nil {
			return err
		}
	}

	if err := p.RunDecorator(); err != nil {
		return err
	}

	return nil
}

// ParseFieldType 解析字段类型
func (*Parser) ParseFieldType(field *entities.Field, ftRow []string) error {
	val := strings.ReplaceAll(ftRow[field.Column], " ", "")
	if val != "" {
		t, err := entities.NewType(val, field)
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
func (*Parser) ParseFieldComment(field *entities.Field, fcRow []string) {
	if field.Column < len(fcRow) {
		val := strings.TrimSpace(fcRow[field.Column])
		if val != "" {
			field.Comment = val
		}
	}
}

// ParseFieldDefaultValue 解析字段默认值
func (*Parser) ParseFieldDefaultValue(field *entities.Field) {
	field.DefaultValue = field.Type.GetDefaultValue()
}

// ParseRow 解析行
func (p *Parser) ParseRow(field *entities.Field, records [][]string) error {
	var value any
	var err error
	for rowIndex, recordRows := range records {
		if field.Column >= len(recordRows) {
			value = nil
		} else {
			cell := recordRows[field.Column]
			if cell == "" {
				if field.Type.GetKind() == reflect.String {
					value, err = field.Type.ParseString(cell)
				} else {
					value = nil
				}
			} else {
				value, err = field.Type.ParseString(cell)
			}
			if err != nil {
				return fmt.Errorf("字段:%s 行:%d 列：%d 错误:%s", field.Name, rowIndex+5, field.Column+1, err)
			}
		}
		p.Table.DataSet[rowIndex] = append(p.Table.DataSet[rowIndex], value)
	}
	return nil
}

// ParseFieldDecorator 解析字段装饰器
func (p *Parser) ParseFieldDecorator(field *entities.Field, fdRow []string, column int) error {
	if column < len(fdRow) {
		val := strings.TrimSpace(fdRow[column])
		if val != "" {
			parts := splitMultiConsRegexp.Split(val, -1)
			for _, part := range parts {
				if err := entities.NewDecorator(p.Table, field, part); err != nil {
					return fmt.Errorf("字段：%s 装饰器：%s %s", field.Name, part, err)
				}
			}
		}
	}
	return nil
}

// RunDecorator 运行装饰器
func (p *Parser) RunDecorator() error {
	for _, d := range p.Table.Decorators {
		if err := d.(entities.ITableDecorator).RunTableDecorator(p.Table); err != nil {
			return fmt.Errorf("装饰器：%s %s", d.(entities.IDecorator).Name(), err)
		}
	}

	for _, field := range p.Table.Fields {
		for _, d := range field.Decorators {
			if err := d.(entities.IFieldDecorator).RunFieldDecorator(p.Table, field); err != nil {
				return fmt.Errorf("字段：%s 装饰器：%s 列：%d %s", field.Name, d.(entities.IDecorator).Name(), field.Column+1, err)
			}
		}
	}
	return nil
}
