package main

import (
	"flag"
	"iter"
	"log"
	"strings"
)

type buildCtxModule struct {
	name    string
	version string
}

func (b *buildCtxModule) String() string {

	if b.version == "" {
		return b.name
	}

	return b.name + "@" + b.version
}

type buildCtx struct {
	name    string
	module  *buildCtxModule
	replace *buildCtxModule
	dir     *workingDirectory
}

func buildContext(root *workingDirectory) iter.Seq[*buildCtx] {
	return func(yield func(*buildCtx) bool) {
		for _, input := range flag.Args() {

			var replace *buildCtxModule
			var module *buildCtxModule = &buildCtxModule{name: input}

			if idx := strings.Index(input, "="); idx > 0 {
				replace = &buildCtxModule{name: input[idx+1:]}
				module.name = input[:idx]
			}

			for _, x := range []*buildCtxModule{module, replace} {

				if nil == x {
					continue
				}

				if idx := strings.Index(x.name, "@"); idx > 0 {
					x.version = x.name[idx+1:]
					x.name = x.name[:idx]
				}
			}

			name := module.name
			name = strings.Replace(name, "/", "_", -1)
			name = strings.Replace(name, ".", "_", -1)

			build, err := root.ChangeDirectory(name)

			if err != nil {
				log.Fatal(err)
			}

			log.Printf("created build directory: \"%s\" for module \"%s\"", build.Path(), module)

			if false == yield(&buildCtx{name, module, replace, build}) {
				return
			}
		}
	}
}
