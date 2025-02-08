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
	UI         bool
	Verify     bool
	SchemaName string
)

func init() {
	flag.BoolVar(&Help, "h", false, "显示帮助信息.")
	flag.BoolVar(&UI, "ui", true, "启动窗体应用.")
	flag.BoolVar(&Verify, "verify", true, "是否对导出的配置进行验证.")
	flag.StringVar(&SchemaName, "s", "", "指定config.toml中的区域进行导出.")
	flag.Usage = usage
	if _, err := toml.DecodeFile("config.toml", &Config); err != nil {
		panic(err)
	}
	Config.UI = UI
	Config.Verify = Verify
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
