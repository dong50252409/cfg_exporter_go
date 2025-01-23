package flatbuffers

import (
	"cfg_exporter/entities"
	"cfg_exporter/implements/flatbuffers/fb_type"
	"cfg_exporter/parser"
)

type FBParse struct {
	*parser.Parser
}

func init() {
	parser.RegisterParser("flatbuffers", NewParser)
}

func NewParser(p *parser.Parser) parser.IParser {
	// 用FlatBuffer类型覆盖默认类型
	register := fb_type.GetTypeRegister()
	entities.MergerTypeRegistry(register)

	return &FBParse{p}
}
