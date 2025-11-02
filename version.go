package main

import (
	"fmt"
	"runtime/debug"
)

func getVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		if replace := info.Main.Replace; nil != replace {
			for {

				if x := replace.Replace; nil != x {
					replace = x
					continue
				}

				break
			}
			return fmt.Sprintf("%s(%s)", replace.Version, replace.Path)
		} else {
			return info.Main.Version
		}
	}

	return "unknown"
}
