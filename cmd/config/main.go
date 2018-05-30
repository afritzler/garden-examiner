package main

import (
	"fmt"
	"os"

	"github.com/afritzler/garden-examiner/pkg"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "config file path expected \n")
		os.Exit(1)
	}
	cfg, err := gube.ReadGardenSetConfig(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not read config: %s \n", err)
		os.Exit(1)
	}
	fmt.Printf("GithubURL: %s \n", cfg.GithubURL)

}
