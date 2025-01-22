package flatbuffers

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	"cfg_exporter/implements/flatbuffers/fb_type"
	"cfg_exporter/parser"
)

type FBParse struct {
	parser.IParser
}

func init() {
	parser.RegisterParser("flatbuffers", NewParser)
}

func NewParser(p *parser.Parser) parser.IParser {
	// 用FlatBuffer类型覆盖默认类型
	register := fb_type.GetTypeRegister()
	entities.MergerTypeRegistry(register)

	p.FieldNameRow = config.Config.Schema["flatbuffers"].FieldNameRow
	p.FieldTypeRow = config.Config.FieldTypeRow
	p.FieldDecoratorRow = config.Config.FieldDecoratorRow
	p.FieldCommentRow = config.Config.FieldCommentRow
	p.BodyStartRow = config.Config.BodyStartRow

	return &FBParse{p}
}

func (p *FBParse) ParseFromFile(path string) (*entities.Table, error) {
	return p.IParser.ParseFromFile(path)
}
