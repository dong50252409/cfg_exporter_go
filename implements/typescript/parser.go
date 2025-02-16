package typescript

import (
	"cfg_exporter/implements/typescript/ts_type"
	"cfg_exporter/parser"
	"cfg_exporter/util"
	"fmt"
	"strings"
)

type TSParse struct {
	*parser.Parser
}

func init() {
	parser.RegisterParser("typescript", NewParser)
}

func NewParser(p *parser.Parser) parser.IParser {
	return &TSParse{p}
}

// ParseFieldType 解析字段类型
func (p *TSParse) ParseFieldType() error {
	fieldTypeRow := p.Records[p.FieldTypeRow-1]
	for _, field := range p.Fields {
		if val := strings.ReplaceAll(fieldTypeRow[field.Column], " ", ""); val != "" {
			t, err := ts_type.NewType(val, field)
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
