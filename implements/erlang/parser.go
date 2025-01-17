package erlang

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"cfg_exporter/implements/erlang/typesystem"
	"cfg_exporter/parser"
)

type ErlParser struct {
	parser.IParser
}

func init() {
	parser.RegisterParser("erlang", NewParser)
}

func NewParser(p *parser.Parser) parser.IParser {
	// 用Erlang类型覆盖默认类型
	register := typesystem.GetTypeRegister()
	entities.MergerTypeRegistry(register)

	p.FieldNameRow = config.Config.Schema["erlang"].FieldNameRow
	p.FieldTypeRow = config.Config.FieldTypeRow
	p.FieldDecoratorRow = config.Config.FieldDecoratorRow
	p.FieldCommentRow = config.Config.FieldCommentRow
	p.BodyStartRow = config.Config.BodyStartRow

	return &ErlParser{p}
}

func (p *ErlParser) ParseFromFile(path string) (*entities.Table, error) {
	return p.IParser.ParseFromFile(path)
}
