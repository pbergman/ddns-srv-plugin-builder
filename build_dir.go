package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type workingDirectory struct {
	path    string
	cleanup bool
}

func (b *workingDirectory) ChangeDirectory(path string) (*workingDirectory, error) {

	if false == filepath.IsAbs(path) {
		path = filepath.Join(b.path, path)
	}

	return getWorkingDirectory(path, b.cleanup)
}

func (b *workingDirectory) Close() error {
	if false == b.cleanup {
		return nil
	}

	return os.RemoveAll(b.path)
}

func (b *workingDirectory) Path() string {
	return b.path
}

func getWorkingDirectory(dir string, cleanup bool) (*workingDirectory, error) {
	if dir != "" {

		if false == filepath.IsAbs(dir) {

			x, err := filepath.Abs(dir)

			if err != nil {
				return nil, err
			}

			dir = x
		}

		stat, err := os.Stat(dir)

		if err != nil {

			if os.IsNotExist(err) {

				if err := os.MkdirAll(dir, 0700); err != nil {
					return nil, err
				}

				return &workingDirectory{path: dir, cleanup: cleanup}, nil

			} else {
				return nil, err
			}
		}

		if stat.IsDir() {
			return &workingDirectory{path: dir, cleanup: cleanup}, nil
		}

		return nil, fmt.Errorf("build dir '%s %s' exists but is not a directory", dir, stat.Mode().String())
	}

	tmp, err := os.MkdirTemp("", "plugin")

	if err != nil {
		return nil, err
	}

	return &workingDirectory{path: tmp, cleanup: cleanup}, nil
}
