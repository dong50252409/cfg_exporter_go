package config

import (
	"github.com/BurntSushi/toml"
)

var Config TomlConfig

func init() {
	if _, err := toml.DecodeFile("config.toml", &Config); err != nil {
		panic(err)
	}
}

type TomlConfig struct {
	Source            string            `toml:"source"`
	FieldCommentRow   int               `toml:"field_comment_row"`
	FieldTypeRow      int               `toml:"field_type_row"`
	FieldDecoratorRow int               `toml:"field_decorator_row"`
	BodyStartRow      int               `toml:"body_start_row"`
	Schema            map[string]Schema `toml:"schema"`
}

type ImportExtension struct {
	Set map[string]any
}

type Schema struct {
	FieldNameRow    int    `toml:"field_name_row"`
	Destination     string `toml:"destination"`
	FilePrefix      string `toml:"file_prefix"`
	TableNamePrefix string `toml:"table_name_prefix"`
	Namespace       string `toml:"namespace"`
}

func (tc *TomlConfig) GetSchema(schemaName string) Schema {
	return tc.Schema[schemaName]
}
