package main

import (
	"cfg_exporter/config"
	"cfg_exporter/erlang"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()

	if flag.NFlag() == 0 || Help {
		flag.Usage()
		return
	}
	run()
}

func run() {
	err := filepath.WalkDir(config.Config.Source, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			t, err := erlang.FromFile(path)
			if err != nil {
				return err
			}
			//render.ToFile(SchemaName, t)
			fmt.Printf("%v", t)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
