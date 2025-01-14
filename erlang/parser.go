package erlang

import (
	"cfg_exporter/entities"
	_ "cfg_exporter/erlang/typesystem"
	"cfg_exporter/parser"
)

func FromFile(path string) (*entities.Table, error) {
	return parser.ParseFromFile(path)
}
