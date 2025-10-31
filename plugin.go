package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"plugin"
	"reflect"
	"text/tabwriter"

	"github.com/libdns/libdns"
)

func createPluginModule(ctx context.Context, build *buildCtx, debug bool) error {

	var cmd = [][]string{
		{"go", "mod", "init", "ddns-server/plugins/" + build.name},
		{"go", "clean", "-modcache"},
		{"go", "get", build.module.String()},
	}

	for i, c := 0, len(cmd); i < c; i++ {
		if err := run(ctx, build.dir, nil, debug, cmd[i]...); err != nil {
			return err
		}
	}

	return nil
}

func checkPlugin(file string) {
	mod, err := plugin.Open(file)

	if err != nil {
		log.Fatal(err)
	}

	name, err := mod.Lookup("PluginModule")

	if err != nil {
		log.Fatal(err)
	}

	version, err := mod.Lookup("PluginVersion")

	if err != nil {
		log.Fatal(err)
	}

	module, err := mod.Lookup("Plugin")

	if err != nil {
		log.Fatal(err)
	}

	var shadow = reflect.ValueOf(module)
	var object = reflect.New(shadow.Elem().Type().Elem()).Interface()
	var buf = new(bytes.Buffer)
	var writer = tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)

	_, _ = fmt.Fprintf(writer, "File\t%s\n", file)
	_, _ = fmt.Fprintf(writer, "PluginModule\t%s\n", *(name.(*string)))
	_, _ = fmt.Fprintf(writer, "PluginVersion\t%s\n", *(version.(*string)))
	_, _ = fmt.Fprintf(writer, "Plugin\t%T\n", object)
	_, _ = fmt.Fprint(writer, "LibDNS implementations:\n")

	var implements = []reflect.Type{
		reflect.TypeOf((*libdns.RecordAppender)(nil)).Elem(),
		reflect.TypeOf((*libdns.RecordGetter)(nil)).Elem(),
		reflect.TypeOf((*libdns.RecordSetter)(nil)).Elem(),
		reflect.TypeOf((*libdns.RecordDeleter)(nil)).Elem(),
		reflect.TypeOf((*libdns.ZoneLister)(nil)).Elem(),
	}

	for i, c := 0, len(implements); i < c; i++ {
		var value string

		if shadow.Elem().Type().Implements(implements[i]) {
			value = "✓"
		} else {
			value = "❌"
		}

		_, _ = fmt.Fprintf(writer, "%+v\t%s\n", implements[i], value)
	}

	_ = writer.Flush()

	var scanner = bufio.NewScanner(buf)

	for scanner.Scan() {
		log.Println(scanner.Text())
	}
}
