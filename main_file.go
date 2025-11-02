package main

import (
	"io"
	"os"
	"path/filepath"
	"text/template"
)

func writeMainFile(build *buildCtx, info *listInfo, debug bool) error {

	tmpl, err := template.New(info.Name).Parse(`
package main

import (
    "{{ .ImportPath }}"
)

var (
	Plugin *{{ .Name }}.Provider
)
`)

	if err != nil {
		return err
	}

	fd, err := os.OpenFile(filepath.Join(build.dir.Path(), "main.go"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		return err
	}

	defer fd.Close()

	var writer io.Writer = fd

	if debug {
		writer = io.MultiWriter(fd, os.Stdout)
	}

	if err := tmpl.Execute(writer, info); err != nil {
		return err
	}

	return nil
}
