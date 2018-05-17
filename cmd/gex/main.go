package main

import (
	"fmt"

	_ "github.com/afritzler/garden-examiner/cmd/gex/profile"
	_ "github.com/afritzler/garden-examiner/cmd/gex/project"
	_ "github.com/afritzler/garden-examiner/cmd/gex/seed"
	// _ "github.com/afritzler/garden-examiner/cmd/gex/select"
	_ "github.com/afritzler/garden-examiner/cmd/gex/shoot"
	_ "github.com/afritzler/garden-examiner/cmd/gex/verb"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/pkg"
	"github.com/mandelsoft/cmdint/pkg/cmdint"
	_ "k8s.io/client-go/tools/clientcmd"
)

func main() {
	cmdint.MainTab().CmdDescription("garden examiner").CmdArgDescription("<options> <command> <options>").
		SetupFunction(setup).
		ArgOption(constants.O_KUBECONFIG).Env("KUBECONFIG").
		ArgOption(constants.O_SEL_SHOOT).Env("GEX_SHOOT").
		ArgOption(constants.O_SEL_PROJECT).Env("GEX_PROJECT").
		ArgOption(constants.O_SEL_SEED).Env("GEX_SEED")
	cmdint.Run()
}

func setup(opts *cmdint.Options) error {
	configfile := opts.GetOptionValue(constants.O_KUBECONFIG)
	if configfile == nil || *configfile == "" {
		return fmt.Errorf("no kubeconfig specified")
	}
	fmt.Printf("kubeconfig is %s\n", *configfile)
	//config, err := clientcmd.BuildConfigFromFlags("", *configfile)
	//if err != nil {
	//	return err
	//}

	//g, err := gube.NewGarden(config)
	g, err := gube.NewGardenFromConfigfile(*configfile)
	if err != nil {
		return err
	}
	c := &context.Context{Garden: gube.NewCachedGarden(g)}
	opts.Context = c
	return nil
}
