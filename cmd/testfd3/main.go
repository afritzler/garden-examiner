package main

import (
	"fmt"
	"os"

	"github.com/afritzler/garden-examiner/cmd/gex/env"
)

func main() {
	f, _ := os.OpenFile("/dev/fd/3", os.O_WRONLY, 0755)
	fmt.Printf("f=%#v\n", f)
	if f == nil {
		fmt.Printf("fd 3 NOT open\n")
	} else {
		fmt.Printf("fd 3 open\n")
	}

	fmt.Printf("MOD %t\n", env.IsModifiable())
}
