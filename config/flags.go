package config

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"strings"
)

var (
	Help       bool
	SchemaName string
)

func init() {
	flag.BoolVar(&Help, "h", false, "显示帮助信息.")
	flag.StringVar(&SchemaName, "s", "", "指定config.toml中的区域进行导出.")
	flag.Usage = usage
	flag.Parse()
	if _, err := toml.DecodeFile("config.toml", &Config); err != nil {
		panic(err)
	}
	Config.SchemaName = SchemaName
}

func usage() {
	var supportedSchema []string

	for key := range Config.Schema {
		supportedSchema = append(supportedSchema, key)
	}
	if _, err := fmt.Fprintf(os.Stdout, `cfg_exporter version: 0.0.1
Usage: cfg_exporter -s %s

Options:
`, strings.Join(supportedSchema, "|")); err != nil {
		panic(err)
	}
	flag.PrintDefaults()
}
