package main

import (
	"fmt"
	"log"
	"time"

	"github.com/afritzler/garden-examiner/cmd/gex/util"
)

func main() {
	doit()
	fmt.Printf("DONE\n")
	time.Sleep(60 * time.Second)
}

func doit() {
	sa, err := util.NewTempFileInput([]byte("test"))
	if err != nil {
		log.Fatalf("cannot get temporary key file name: %s", err)
	}
	defer sa.CleanupFunction()()
	time.Sleep(60 * time.Second)
}
