package main

import (
	"cfg_exporter/config"
	"cfg_exporter/erlang"
	"cfg_exporter/reader"
	"cfg_exporter/render"
	"flag"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()

	if flag.NFlag() == 0 || Help {
		flag.Usage()
		return
	}
	err := filepath.WalkDir(config.Config.Source, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			err := run(path)
			if err != nil {
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
	t, err := erlang.FromFile(path)
	if err != nil {
		return err
	}
	err = render.ToFile(SchemaName, t)
	if err != nil {
		return err
	}
	return nil
}
