package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const VERSION = "0.1.0"

var (
	Help              bool
	Source            string
	Destination       string
	FieldNameRow      int
	FieldTypeRow      int
	FieldDecoratorRow int
	FieldCommentRow   int
	BodyStartRow      int
	Verify            bool
	TomlPath          string
	SchemaName        string
	FilePrefix        string
	TableNamePrefix   string
	Flatc             string
	Namespace         string
)

func init() {
	flag.BoolVar(&Help, "h", false, "显示帮助信息.")
	flag.StringVar(&Source, "src", "", "指定配置表文件路径.")
	flag.StringVar(&Destination, "dest", "", "指定导出文件路径.")
	flag.IntVar(&FieldNameRow, "fnr", 2, "指定字段名所在行.")
	flag.IntVar(&FieldTypeRow, "ftr", 3, "指定字段类型所在行.")
	flag.IntVar(&FieldDecoratorRow, "fdr", 4, "指定字段装饰器所在行.")
	flag.IntVar(&FieldCommentRow, "fcr", 1, "指定字段注释所在行.")
	flag.IntVar(&BodyStartRow, "bsr", 6, "指定表体开始行.")
	flag.BoolVar(&Verify, "verify", false, "是否对生成的文件进行校验.")
	flag.StringVar(&TomlPath, "tp", "config.toml", "指定config.toml文件路径.")
	flag.StringVar(&SchemaName, "sn", "", "指定config.toml中的区域进行导出.")
	flag.StringVar(&FilePrefix, "fp", "", "指定导出文件的前缀.")
	flag.StringVar(&TableNamePrefix, "tnp", "", "指定导出表的前缀.")
	flag.StringVar(&Flatc, "flatc", "", "指定flatc路径.")
	flag.StringVar(&Namespace, "ns", "", "指定flatbuffers的namespace.")
	flag.Usage = usage
	flag.Parse()

	if TomlPath != "" && SchemaName != "" {
		if _, err := os.Stat(TomlPath); os.IsNotExist(err) {
			NewTomlConfig(TomlPath)
		} else {
			panic(fmt.Errorf("配置文件不存在 %s", TomlPath))
		}
	} else {
		NewTomlConfigByFlags()
	}
}

func usage() {
	supportedSchema := []string{
		"erlang",
		"flatbuffers",
		"json",
		"typescript",
	}

	if _, err := fmt.Fprintf(os.Stdout, `cfg_exporter version: %s
Usage: cfg_exporter -s %s

Options:
`, VERSION, strings.Join(supportedSchema, "|")); err != nil {
		panic(err)
	}
	flag.PrintDefaults()
}
