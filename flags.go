package main

import (
	"flag"
	"fmt"
	"os"
)

func init() {
	flag.String("build-dir", "", "build directory, when nothing is provided a temp directory will be created")
	flag.Bool("no-build", false, "by default it wil build the plugin, use this option with build-dir and no-cleanup to create an project template")
	flag.Bool("no-cleanup", false, "will not remove build dir when provided")
	flag.Bool("debug", false, "print debug information")
	flag.String("save-path", "/usr/share/ddns-srv", "the path where plugin will be saved on successful build")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <package...>\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), `

To replace a package you could define a package like:
	example.com/organisation/package=github.com/me/package@package-fork

This will place example.com/organisation/package for github.com/me/package@package-fork

`)
	}
}

func inputOptionDefault[T bool | string]() T {
	var ret any = new(T)

	switch ret.(type) {
	case *bool:
		*(ret.(*bool)) = false
	}

	return *((ret).(*T))
}

func inputOption[T bool | string](name string) T {

	var input = flag.Lookup(name)

	if input == nil {
		return inputOptionDefault[T]()
	}

	if getter, x := input.Value.(flag.Getter); x {
		if ret, y := getter.Get().(T); y {
			return ret
		}
	}

	return inputOptionDefault[T]()
}
