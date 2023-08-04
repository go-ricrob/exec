package task

import (
	"fmt"
	"runtime/debug"
)

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		fmt.Printf("build info %v", info)
	}
}
