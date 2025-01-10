package erlang

import (
	"cfg_exporter/entities"
	_ "cfg_exporter/erlang/erl_type"
	"cfg_exporter/parser"
)

func FromFile(path string) (*entities.Table, error) {
	return parser.FromFile(path)
}
