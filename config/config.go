package config

var Config TomlConfig

type TomlConfig struct {
	Source            string            `toml:"source"`
	FieldCommentRow   int               `toml:"field_comment_row"`
	FieldTypeRow      int               `toml:"field_type_row"`
	FieldDecoratorRow int               `toml:"field_decorator_row"`
	BodyStartRow      int               `toml:"body_start_row"`
	UI                bool              `toml:"ui"`
	Verify            bool              `toml:"verify"`
	SchemaName        string            `toml:"default_schema"`
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

func (tc *TomlConfig) GetSchema(schemaName string) Schema {
	return tc.Schema[schemaName]
}
