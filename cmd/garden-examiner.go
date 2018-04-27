package main

import (
	"fmt"
	"os"

	landscaper "github.com/afritzler/garden-examiner/pkg"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	fmt.Println("hello gubernetix!")
	config, err := clientcmd.BuildConfigFromFlags("", os.Args[1])
	if err != nil {
		panic(err)
	}

	g, err := landscaper.NewGarden(config)
	if err != nil {
		panic(err)
	}
	shoots, err := g.GetShoots()
	if err != nil {
		panic(err)
	}
	for _, s := range shoots {
		c, err := s.GetNodeCount()
		if err != nil {
			fmt.Printf("%30s / %30s: %s\n", s.GetName(), s.GetSeedName(), err)
		} else {
			fmt.Printf("%30s / %30s: %d\n", s.GetName(), s.GetSeedName(), c)
		}
	}
}
