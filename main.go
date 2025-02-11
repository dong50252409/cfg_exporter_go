package main

import (
	"cfg_exporter/config"
	"cfg_exporter/entities"
	_ "cfg_exporter/implements/erlang"
	_ "cfg_exporter/implements/flatbuffers"
	_ "cfg_exporter/implements/json"
	_ "cfg_exporter/implements/typescript"
	"cfg_exporter/parser"
	"cfg_exporter/reader"
	"cfg_exporter/render"
	"errors"
	"flag"
	"os"
	"path/filepath"
)

func main() {

	if flag.NFlag() == 0 || config.Help {
		flag.Usage()
		return
	}
	if config.Config.UI {
		startUI()
	} else {
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
}

func run(path string) error {
	records, err := readFile(path)
	if err != nil {
		if errors.Is(err, reader.ErrorTableTempFile) {
			return nil
		}
		return err
	}

	t, err := parserTable(path, records)
	if err != nil {
		return err
	}

	err = renderTable(t)
	if err != nil {
		return err
	}
	return nil
}

func readFile(path string) ([][]string, error) {
	r, err := reader.NewReader(path)
	if err != nil {
		return nil, err
	}
	records, err := r.Read()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func parserTable(path string, records [][]string) (*entities.Table, error) {
	p, err := parser.NewParser(config.Config.SchemaName)
	if err != nil {
		return nil, err
	}

	t, err := p.ParseFromFile(path, records)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func renderTable(t *entities.Table) error {
	if r, err := render.NewRender(config.Config.SchemaName, t); err != nil {
		return err
	} else {
		if err = r.Execute(); err != nil {
			return err
		}
		if config.Config.Verify {
			if err = r.Verify(); err != nil {
				return err
			}
		}
	}
	return nil
}
