package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

var Config TomlConfig

type TomlConfig struct {
	Source            string            `toml:"source"`
	FieldCommentRow   int               `toml:"field_comment_row"`
	FieldTypeRow      int               `toml:"field_type_row"`
	FieldDecoratorRow int               `toml:"field_decorator_row"`
	BodyStartRow      int               `toml:"body_start_row"`
	Verify            bool              `toml:"verify"`
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
	Flatc           string `toml:"flatc"`
	Namespace       string `toml:"namespace"`
}

// NewTomlConfig 载入toml配置
func NewTomlConfig(path string) {
	if _, err := toml.DecodeFile(path, &Config); err != nil {
		panic(fmt.Errorf("解析配置文件失败 %s", err))
	}
}

// NewTomlConfigByFlags 根据命令行参数生成配置
func NewTomlConfigByFlags() {
	Config = TomlConfig{
		Source:            Source,
		FieldTypeRow:      FieldTypeRow,
		FieldDecoratorRow: FieldDecoratorRow,
		FieldCommentRow:   FieldCommentRow,
		BodyStartRow:      BodyStartRow,
		Verify:            Verify,
		Schema: map[string]Schema{
			SchemaName: {
				FieldNameRow:    FieldNameRow,
				Destination:     Destination,
				FilePrefix:      FilePrefix,
				TableNamePrefix: TableNamePrefix,
				Flatc:           Flatc,
				Namespace:       Namespace,
			},
		},
	}
}
