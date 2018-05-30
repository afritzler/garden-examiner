package cmdint

import (
	"fmt"
	"os"
)

var maintab *_CmdTab = NewCmdTab(os.Args[0])

func MainTab() *_CmdTab {
	return maintab
}

func Run() {
	err := maintab.Execute(nil, os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
