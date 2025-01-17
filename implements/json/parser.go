package json

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"cfg_exporter/parser"
)

type JParser struct {
	parser.IParser
}

func init() {
	parser.RegisterParser("json", NewParser)
}

func NewParser(p *parser.Parser) parser.IParser {
	entities.RecoverBaseTypeRegister()
	p.FieldNameRow = config.Config.Schema["json"].FieldNameRow
	p.FieldTypeRow = config.Config.FieldTypeRow
	p.FieldDecoratorRow = config.Config.FieldDecoratorRow
	p.FieldCommentRow = config.Config.FieldCommentRow
	p.BodyStartRow = config.Config.BodyStartRow
	return &JParser{p}
}

func (p *JParser) ParseFromFile(path string) (*entities.Table, error) {
	return p.IParser.ParseFromFile(path)
}
