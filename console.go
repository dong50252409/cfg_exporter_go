//go:build console

package main

import (
	"cfg_exporter/config"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if flag.NFlag() == 0 || config.Help {
		flag.Usage()
	} else {
		err := filepath.WalkDir(config.Config.Source,
			func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() {
					if err = run(path, config.SchemaName); err != nil {
						return err
					}
				}
				return nil
			})

		if err != nil {
			fmt.Println(err)
		}
	}
}
