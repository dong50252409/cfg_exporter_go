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
	Source             string            `toml:"source"`
	FieldCommentRow    int               `toml:"field_comment_row"`
	FieldTypeRow       int               `toml:"field_type_row"`
	FieldConstraintRow int               `toml:"field_constraint_row"`
	BodyStartRow       int               `toml:"body_start_row"`
	Schema             map[string]Schema `toml:"schema"`
}

type ImportExtension struct {
	Set map[string]any
}

type Schema struct {
	FieldNameRow int    `toml:"field_name_row"`
	Destination  string `toml:"destination"`
	FilePrefix   string `toml:"file_prefix"`
	RecordPrefix string `toml:"record_prefix"`
}

func (tc *TomlConfig) GetSchema(schemaName string) Schema {
	return tc.Schema[schemaName]
}
