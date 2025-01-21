package typescript

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"cfg_exporter/implements/typescript/ts_type"
	"cfg_exporter/parser"
)

type TSParse struct {
	parser.IParser
}

func init() {
	parser.RegisterParser("typescript", NewParser)
}

func NewParser(p *parser.Parser) parser.IParser {
	// 用TypeScript类型覆盖默认类型
	register := ts_type.GetTypeRegister()
	entities.MergerTypeRegistry(register)

	p.FieldNameRow = config.Config.Schema["typescript"].FieldNameRow
	p.FieldTypeRow = config.Config.FieldTypeRow
	p.FieldDecoratorRow = config.Config.FieldDecoratorRow
	p.FieldCommentRow = config.Config.FieldCommentRow
	p.BodyStartRow = config.Config.BodyStartRow

	return &TSParse{p}
}

func (p *TSParse) ParseFromFile(path string) (*entities.Table, error) {
	return p.IParser.ParseFromFile(path)
}
