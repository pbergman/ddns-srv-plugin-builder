package main

import (
	"bytes"
	"context"
	"encoding/json"
)

type listInfo struct {
	ImportPath string `json:"ImportPath"`
	Name       string `json:"Name"`
}

func getModuleInfo(ctx context.Context, build *buildCtx, debug bool) (*listInfo, error) {
	var out = new(bytes.Buffer)

	if err := run(ctx, build.dir, out, debug, "go", "list", "-json", build.module.name); err != nil {
		return nil, err
	}

	var module *listInfo

	if err := json.NewDecoder(out).Decode(&module); err != nil {
		return nil, err
	}

	return module, nil
}
