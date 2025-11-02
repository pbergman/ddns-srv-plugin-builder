package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {

	flag.Parse()

	if flag.NArg() == 0 && false == inputOption[bool]("version") {
		flag.Usage()
		os.Exit(1)
	}

	cwd, err := getWorkingDirectory(
		inputOption[string]("build-dir"),
		inputOption[bool]("no-cleanup") == false,
	)

	if err != nil {
		log.Fatal(err)
	}

	if inputOption[bool]("version") {
		fmt.Println(getVersion())
		os.Exit(0)
	}

	var debug = inputOption[bool]("debug")
	var noBuild = inputOption[bool]("no-build")

	save, err := filepath.Abs(inputOption[string]("save-path"))

	if err != nil {
		log.Fatal(err)
	}

	var ctx context.Context = context.Background()

	log.Printf("build version: %s", getVersion())

	for build := range buildContext(cwd) {

		if err := createPluginModule(ctx, build, debug); err != nil {
			log.Fatal(err)
		}

		module, err := getModuleInfo(ctx, build, debug)

		if nil != err {
			log.Fatal(err)
		}

		if nil != build.replace {
			goModReplace(ctx, build, debug)
			goModTidy(ctx, build, debug)
		}

		if err := writeMainFile(build, module, debug); err != nil {
			log.Fatal(err)
		}

		goModTidy(ctx, build, debug)

		if nil != build.replace {

			// query version again
			module, err = getModuleInfo(ctx, build, debug)

			if err != nil {
				log.Fatal(err)
			}
		}

		if false == noBuild {
			checkPlugin(goBuildPlugin(save, module, build, debug))
		}
	}

}
