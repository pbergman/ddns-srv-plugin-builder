package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"os/exec"
	"path/filepath"
)

func run(ctx context.Context, cwd *workingDirectory, out io.Writer, debug bool, run ...string) error {

	var name string
	var args []string

	if len(run) == 1 {
		name, args = run[0], ([]string)(nil)
	} else {
		name, args = run[0], run[1:]
	}

	var stderr = new(bytes.Buffer)

	cmd := exec.CommandContext(ctx, name, args...)

	log.Printf("exec: %s", cmd)

	cmd.Stderr = io.MultiWriter(&output{prefix: "> ", out: log.Default()}, stderr)
	cmd.Dir = cwd.Path()

	if debug {
		cmd.Stdout = &output{prefix: "> ", out: log.Default()}
	}

	if out != nil {
		if cmd.Stdout != nil {
			cmd.Stdout = io.MultiWriter(cmd.Stdout, out)
		} else {
			cmd.Stdout = out
		}
	}

	if err := cmd.Run(); err != nil {

		if x, ok := err.(*exec.ExitError); ok {
			x.Stderr = stderr.Bytes()
		}

		return errors.Join(errors.New(stderr.String()), err)
	}

	return nil
}

func goModTidy(ctx context.Context, build *buildCtx, debug bool) {
	if err := run(ctx, build.dir, nil, debug, "go", "mod", "tidy"); err != nil {
		log.Fatal(err)
	}
}

func goModReplace(ctx context.Context, build *buildCtx, debug bool) {
	if err := run(ctx, build.dir, nil, debug, "go", "mod", "edit", "-replace", build.module.String()+"="+build.replace.String()); err != nil {
		log.Fatal(err)
	}
}

func goBuildPlugin(path string, module *listInfo, build *buildCtx, debug bool) string {
	var args = []string{
		"go",
		"build",
		"-buildmode", "plugin",
		"-ldflags", "-s -w",
		"-o", filepath.Join(path, build.name+".so"),
	}

	if err := run(context.Background(), build.dir, nil, debug, args...); err != nil {
		log.Fatal(err)
	}

	log.Printf("build plugin for module %s\n", module.Name)

	return args[7]
}
