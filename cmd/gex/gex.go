package main

import (
	"fmt"
	"os"
	"os/user"

	_ "github.com/afritzler/garden-examiner/cmd/gex/profile"
	_ "github.com/afritzler/garden-examiner/cmd/gex/project"
	_ "github.com/afritzler/garden-examiner/cmd/gex/seed"
	// _ "github.com/afritzler/garden-examiner/cmd/gex/select"
	_ "github.com/afritzler/garden-examiner/cmd/gex/garden"
	_ "github.com/afritzler/garden-examiner/cmd/gex/shoot"
	_ "github.com/afritzler/garden-examiner/cmd/gex/verb"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/pkg"
	"github.com/afritzler/garden-examiner/pkg/data"
	"github.com/mandelsoft/cmdint/pkg/cmdint"
	"github.com/mandelsoft/filepath/pkg/filepath"
	_ "k8s.io/client-go/tools/clientcmd"
)

func main() {
	usr, err := user.Current()
	gexdir := ""
	if err == nil {
		gexdir = filepath.Join(usr.HomeDir, ".gex")
	}
	cmdint.MainTab().CmdDescription("garden examiner").CmdArgDescription("<options> <command> <options>").
		SetupFunction(setup).
		ArgOption(constants.O_GEXDIR).Env("GEXDIR").Default(gexdir).
		ArgOption(constants.O_GEXCONFIG).Env("GEXCONFIG").
		ArgOption(constants.O_KUBECONFIG).Env("KUBECONFIG").
		ArgOption(constants.O_SEL_SHOOT).Env("GEX_SHOOT").
		ArgOption(constants.O_SEL_PROJECT).Env("GEX_PROJECT").
		ArgOption(constants.O_SEL_SEED).Env("GEX_SEED").
		ArgOption(constants.O_SEL_GARDEN).Env("GEX_GARDEN")

	cmdint.Run()
}

func setup(opts *cmdint.Options) error {
	c := &context.Context{}
	opts.Context = c

	c.Gexdir = *opts.GetOptionValue(constants.O_GEXDIR)
	gexconfig := opts.GetOptionValue(constants.O_GEXCONFIG)
	if data.IsEmpty(gexconfig) && !data.IsEmpty(c.Gexdir) {
		cfg := filepath.Join(c.Gexdir, "config")
		if _, err := os.Stat(cfg); !os.IsNotExist(err) {
			gexconfig = &cfg
		}
	}
	if !data.IsEmpty(gexconfig) {
		cfg, err := gube.NewGardenSetConfig(*gexconfig)
		if err != nil {
			return err
		}
		var gardenConfig gube.GardenConfig
		c.GardenSetConfig = cfg
		selGarden := opts.GetOptionValue(constants.O_SEL_GARDEN)
		if selGarden != nil {
			gardenConfig, err = cfg.GetConfig(*selGarden)
			c.Name = *selGarden
		} else {
			gardenConfig, err = cfg.GetConfig("")
			c.Name = cfg.GetDefault()
		}
		if err != nil {
			return err
		}
		c.GardenConfig = gardenConfig
		g, err := gardenConfig.GetGarden()
		if err != nil {
			return err
		}
		c.Garden = gube.NewCachedGarden(g)
	} else {
		configfile := opts.GetOptionValue(constants.O_KUBECONFIG)
		if data.IsEmpty(configfile) {
			return fmt.Errorf("no kubeconfig or gexconfig specified")
		}
		c.ByKubeconfig = true
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
		c.GardenSetConfig = gube.NewDefaultGardenSetConfig(g)
		c.Garden = gube.NewCachedGarden(g)
		c.Name = "default"
	}
	return nil
}
