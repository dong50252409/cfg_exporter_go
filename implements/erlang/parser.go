package erlang

import (
	"cfg_exporter/implements/erlang/erl_type"
	"cfg_exporter/parser"
	"cfg_exporter/util"
	"fmt"
	"strings"
)

type ErlParser struct {
	*parser.Parser
}

func init() {
	parser.RegisterParser("erlang", NewParser)
}

func NewParser(p *parser.Parser) parser.IParser {
	return &ErlParser{p}
}

// ParseFieldType 解析字段类型
func (p *ErlParser) ParseFieldType() error {
	fieldTypeRow := p.Records[p.FieldTypeRow-1]
	for _, field := range p.Fields {
		if val := strings.ReplaceAll(fieldTypeRow[field.Column], " ", ""); val != "" {
			t, err := erl_type.NewType(val, field)
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
