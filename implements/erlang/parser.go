package erlang

import (
	"cfg_exporter/entities"
	"cfg_exporter/implements/erlang/erl_type"
	"cfg_exporter/parser"
)

type ErlParser struct {
	*parser.Parser
}

func init() {
	parser.RegisterParser("erlang", NewParser)
}

func NewParser(p *parser.Parser) parser.IParser {
	// 用Erlang类型覆盖默认类型
	register := erl_type.GetTypeRegister()
	entities.MergerTypeRegistry(register)

	return &ErlParser{p}
}
