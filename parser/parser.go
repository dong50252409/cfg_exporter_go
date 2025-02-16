package parser

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"cfg_exporter/util"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

type IParser interface {
	ParseFieldName() error
	ParseFieldType() error
	ParseFieldDecorators() error
	ParseFieldComment()
	ParseFieldDefaultValue()
	ParseRow() error
	RunDecorators() error
	GetTable() *entities.Table
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
func NewParser(path string, schemaName string, records [][]string) (IParser, error) {
	if parser, ok := parserRegistry[schemaName]; ok {
		filename := filepath.Base(path)
		name, err := util.SubTableName(filename)
		if err != nil {
			return nil, err
		}
		p := &Parser{
			Table: &entities.Table{
				Path:       filepath.Dir(path),
				Filename:   filename,
				Name:       name,
				Decorators: make([]entities.ITableDecorator, 0),
				Records:    records,
			},
			FieldNameRow:      config.Config.Schema[schemaName].FieldNameRow,
			FieldTypeRow:      config.Config.FieldTypeRow,
			FieldDecoratorRow: config.Config.FieldDecoratorRow,
			FieldCommentRow:   config.Config.FieldCommentRow,
			BodyStartRow:      config.Config.BodyStartRow,
		}

		return parser(p), nil
	}
	return nil, fmt.Errorf("未找到解析器：%s", schemaName)
}

func CloneParser(schemaName string, tbl *entities.Table) (IParser, error) {
	if parser, ok := parserRegistry[schemaName]; ok {
		p := &Parser{
			Table:             tbl,
			FieldNameRow:      config.Config.Schema[schemaName].FieldNameRow,
			FieldTypeRow:      config.Config.FieldTypeRow,
			FieldDecoratorRow: config.Config.FieldDecoratorRow,
			FieldCommentRow:   config.Config.FieldCommentRow,
			BodyStartRow:      config.Config.BodyStartRow,
		}

		iParser := parser(p)
		if err := iParser.ParseFieldType(); err != nil {
			return nil, err
		}
		return iParser, nil
	}
	return nil, fmt.Errorf("未找到解析器：%s", schemaName)
}

// Parse 解析配表
func (p *Parser) Parse() error {
	err := p.ParseFieldName()
	if err != nil {
		return err
	}

	err = p.ParseFieldType()
	if err != nil {
		return err
	}

	err = p.ParseFieldDecorators()
	if err != nil {
		return err
	}

	p.ParseFieldComment()

	err = p.ParseRow()
	if err != nil {
		return err
	}

	p.ParseFieldDefaultValue()

	return nil
}

// ParseFieldName 解析字段名
func (p *Parser) ParseFieldName() error {
	fieldNameRow := p.Records[p.FieldNameRow-1]
	p.Fields = make([]*entities.Field, 0, len(fieldNameRow))

	for column, colIndex := 0, 0; column < len(fieldNameRow); column++ {
		if val := strings.TrimSpace(fieldNameRow[column]); val != "" {
			field := &entities.Field{Column: column, ColIndex: colIndex, Name: val, Decorators: make(map[string]entities.IFieldDecorator)}
			p.Fields = append(p.Fields, field)
			colIndex++
		}
	}

	// 检查是否有重复字段名
	fieldSet := make(map[string]*entities.Field)
	for _, field := range p.Fields {
		if f, ok := fieldSet[field.Name]; ok {
			return fmt.Errorf("单元格：%s、%s\n错误：%s 重复定义", util.ToCell(p.FieldNameRow, f.Column), util.ToCell(p.FieldNameRow, field.Column), f.Name)
		} else {
			fieldSet[field.Name] = field
		}
	}
	return nil
}

// ParseFieldType 解析字段类型
func (p *Parser) ParseFieldType() error {
	fieldTypeRow := p.Records[p.FieldTypeRow-1]

	for _, field := range p.Fields {
		if val := strings.ReplaceAll(fieldTypeRow[field.Column], " ", ""); val != "" {
			t, err := entities.NewType(val, field)
			if err != nil {
				return fmt.Errorf("单元格：%s\n错误：%s", util.ToCell(p.FieldTypeRow, field.Column), err)
			}
			field.Type = t
		} else {
			return fmt.Errorf("单元格：%s\n错误：类型不能为空", util.ToCell(p.FieldTypeRow, field.Column))
		}
	}
	return nil
}

// ParseFieldDecorators 解析装饰器信息
func (p *Parser) ParseFieldDecorators() error {
	fieldDecoratorRow := p.Records[p.FieldDecoratorRow-1]

	// 第一列默认为主键列
	if fieldDecoratorRow == nil {
		fieldDecoratorRow = []string{"p_key\n"} // 没有任何装饰器
	} else if !strings.Contains(fieldDecoratorRow[0], "p_key") {
		fieldDecoratorRow[0] = "p_key\n" + fieldDecoratorRow[0]
	}

	for column, val := range fieldDecoratorRow {
		if val = strings.TrimSpace(val); val != "" {
			var field *entities.Field
			if field = p.GetFieldByColumn(column); field == nil {
				field = &entities.Field{Column: column}
			}

			parts := splitMultiConsRegexp.Split(val, -1)
			for _, part := range parts {
				if err := entities.NewDecorator(p.Table, field, part); err != nil {
					return fmt.Errorf("单元格：%s\n错误：%s 装饰器 %s", field.Name, part, err)
				}
			}
		}
	}
	return nil
}

// ParseFieldComment 解析字段注释
func (p *Parser) ParseFieldComment() {
	fileCommentRow := p.Records[p.FieldCommentRow-1]
	for _, field := range p.Fields {
		if len(fileCommentRow) > field.Column {
			field.Comment = strings.TrimSpace(fileCommentRow[field.Column])
		}
	}
}

// ParseFieldDefaultValue 解析字段默认值
func (p *Parser) ParseFieldDefaultValue() {
	for _, field := range p.Fields {
		field.DefaultValue = field.Type.DefaultValue()
	}
}

// ParseRow 解析行
func (p *Parser) ParseRow() error {
	records := p.Records[p.BodyStartRow-1:]
	// 初始化数据
	p.DataSet = make([][]interface{}, len(records))
	for i := 0; i < len(records); i++ {
		p.DataSet[i] = make([]interface{}, len(p.Fields))
	}

	maxRow := 0 // 记录最大有效行数
	for _, field := range p.Fields {
		rowIndex := 0
		for row := 0; row < len(records); row++ {
			recordRows := records[row]
			if len(recordRows) > field.Column {
				if cell := recordRows[field.Column]; cell != "" {
					value, err := field.Type.ParseString(cell)
					if err != nil {
						return fmt.Errorf("单元格：%s\n错误：%s", util.ToCell(rowIndex+p.BodyStartRow, field.Column), err)
					}
					p.DataSet[rowIndex][field.ColIndex] = value
					rowIndex++
				}
			}
		}
		maxRow = max(maxRow, rowIndex)
	}
	p.DataSet = p.DataSet[:maxRow] // 去除无效行

	return nil
}

// RunDecorators 运行装饰器
func (p *Parser) RunDecorators() error {
	if len(p.Fields) == 0 {
		return nil
	}

	for _, d := range p.Decorators {
		if err := d.RunTableDecorator(p.Table); err != nil {
			return fmt.Errorf("装饰器：%s\n错误：%s", d.Name(), err)
		}
	}

	for _, field := range p.Fields {
		for _, d := range field.Decorators {
			if err := d.RunFieldDecorator(p.Table, field); err != nil {
				return fmt.Errorf("装饰器：%s\n错误：%s", d.Name(), err)
			}
		}
	}
	return nil
}

// GetTable 获取配表
func (p *Parser) GetTable() *entities.Table {
	return p.Table
}
