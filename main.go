package main

import (
	"cfg_exporter/config"
	_ "cfg_exporter/implements/erlang"
	_ "cfg_exporter/implements/flatbuffers"
	_ "cfg_exporter/implements/json"
	_ "cfg_exporter/implements/typescript"
	"cfg_exporter/parser"
	"cfg_exporter/reader"
	"cfg_exporter/render"
	"flag"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()

	if flag.NFlag() == 0 || config.Help {
		flag.Usage()
		return
	}
	err := filepath.WalkDir(config.Config.Source, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			if err := run(path); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func run(path string) error {
	if ok := reader.CheckSupport(path); !ok {
		return nil
	}

	p, err := parser.NewParser(config.SchemaName)
	if err != nil {
		return err
	}

	t, err := p.ParseFromFile(path)
	if err != nil {
		return err
	}

	if err = render.ToFile(config.SchemaName, t); err != nil {
		return err
	}

	return nil
}
