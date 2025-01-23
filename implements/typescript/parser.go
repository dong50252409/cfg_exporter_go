package typescript

import (
	"cfg_exporter/entities"
	"cfg_exporter/implements/typescript/ts_type"
	"cfg_exporter/parser"
)

type TSParse struct {
	*parser.Parser
}

func init() {
	parser.RegisterParser("typescript", NewParser)
}

func NewParser(p *parser.Parser) parser.IParser {
	// 用TypeScript类型覆盖默认类型
	register := ts_type.GetTypeRegister()
	entities.MergerTypeRegistry(register)
	return &TSParse{p}
}
