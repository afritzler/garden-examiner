package verb

import (
	"fmt"
	"os"
	"strings"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/env"
	"github.com/afritzler/garden-examiner/cmd/gex/output"

	"github.com/afritzler/garden-examiner/pkg"
	"github.com/afritzler/garden-examiner/pkg/data"

	"github.com/mandelsoft/filepath/pkg/filepath"
)

func init() {
	NewVerb("select", cmdint.MainTab()).CmdArgDescription("clear|<type> ...").
		CmdDescription("general select command",
			"The first argument is the element type followed by",
			"an optional element name.",
			"With clear the selection can be undone. If nothing is specified",
			"the actual selection is shown",
		).
		FlagOption(constants.O_DOWNLOAD).Short('d').Description("download kubeconfig").
		FlagOption(constants.O_EXPORT).Short('e').Description("export env KUBECONFIG (implies -d)").
		DefaultFunction(cmd_select).
		SimpleCommand("clear", cmd_clear).
		CmdArgDescription("{project|seed|shoot}").
		CmdDescription("clear given selection")
}

func cmd_select(opts *cmdint.Options) error {
	found := 0
	export := opts.IsFlag(constants.O_EXPORT)
	download := opts.IsFlag(constants.O_DOWNLOAD) || export
	ctx := context.Get(opts)
	if ctx.Gexdir == "" && download {
		return fmt.Errorf("No GEXDIR set")
	}

	gap := opts.GetOptionValue(constants.O_SEL_GARDEN)
	shp := opts.GetOptionValue(constants.O_SEL_SHOOT)
	sep := opts.GetOptionValue(constants.O_SEL_SEED)
	prp := opts.GetOptionValue(constants.O_SEL_PROJECT)
	if !data.IsEmpty(gap) {
		fmt.Printf("GARDEN  = %s\n", *gap)
		found++
	} else {
		fmt.Printf("GARDEN  = %s  (defaulted by kubeconfig or gexconfig)\n", ctx.Name)
	}
	if !data.IsEmpty(shp) {
		fmt.Printf("SHOOT   = %s\n", *shp)
		found++
	}
	if !data.IsEmpty(prp) {
		fmt.Printf("PROJECT = %s\n", *prp)
		found++
	}
	if !data.IsEmpty(sep) {
		fmt.Printf("SEED    = %s\n", *sep)
		found++
	}
	if found == 0 {
		return fmt.Errorf("no selection found")
	} else {
		if download {
			if !data.IsEmpty(shp) {
				a := strings.Split(*shp, "/")
				s, err := ctx.Garden.GetShoot(gube.NewShootName(a[0], a[1]))
				if err != nil {
					return err
				}
				return Download(ctx.ByKubeconfig, export, s, ctx.Gexdir, ctx.Name, "projects", a[0], a[1])
			} else {
				if !data.IsEmpty(sep) {
					s, err := ctx.Garden.GetSeed(*sep)
					if err != nil {
						return err
					}
					return Download(ctx.ByKubeconfig, export, s, ctx.Gexdir, ctx.Name, "seeds", *sep)
				} else {
					if !data.IsEmpty(gap) && data.IsEmpty(prp) {
						return Download(ctx.ByKubeconfig, export, ctx.Garden, ctx.Gexdir, ctx.Name)
					}
				}
			}
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////
// clear sub command

type clear_output struct {
	*select_output
}

func (this *clear_output) Out(opts *cmdint.Options) error {
	garden := opts.GetOptionValue(constants.O_SEL_GARDEN)
	shoot := opts.GetOptionValue(constants.O_SEL_SHOOT)
	project := opts.GetOptionValue(constants.O_SEL_PROJECT)
	seed := opts.GetOptionValue(constants.O_SEL_SEED)

	if len(opts.Arguments) > 0 {
		for _, n := range opts.Arguments {
			b, d := cmdint.SelectBest(n, "shoot", "seed", "project", "garden")
			if d > len(n)/2 {
				return fmt.Errorf("unknown selection type '%s'", n)
			}
			fmt.Printf("clearing %s selection\n", b)
			switch b {
			case "garden":
				garden = nil
			case "shoot":
				shoot = nil
			case "seed":
				seed = nil
			case "project":
				project = nil
			}
		}
	} else {
		shoot = nil
		seed = nil
		project = nil
	}
	this.Write(garden, shoot, project, seed)
	return nil
}

func cmd_clear(opts *cmdint.Options) error {
	return (&clear_output{select_output: NewSelectOutput(false, false)}).Out(opts)
}

////////////////////////////////////////////////////////////////////////////
// general select output

type select_output struct {
	*output.SingleElementOutput
	download bool
	export   bool
}

var _ output.Output = &select_output{}

func NewSelectOutput(download, export bool) *select_output {
	return &select_output{output.NewSingleElementOutput(), download || export, export}
}

func (this *select_output) Out(ctx *context.Context) error {
	var err error
	garden := ctx.Name
	shoot := ""
	seed := ""
	project := ""
	if ctx.Gexdir == "" && this.export {
		return fmt.Errorf("No GEXDIR set")
	}
	switch e := this.Elem.(type) {
	case gube.GardenConfig:
		garden = e.GetName()
		if this.download {
			Download(ctx.ByKubeconfig, this.export, e, ctx.Gexdir, garden)
		}
	case gube.Shoot:
		shoot = e.GetName().String()
		project = e.GetName().GetProjectName()
		seed = e.GetSeedName()
		if this.download {
			err = Download(ctx.ByKubeconfig, this.export, e, ctx.Gexdir, garden, "projects", project, shoot)
		}
	case gube.Seed:
		seed = e.GetName()
		if this.download {
			err = Download(ctx.ByKubeconfig, this.export, e, ctx.Gexdir, garden, "seeds", seed)
		}
	case gube.Project:
		project = e.GetName()
	default:
		panic(fmt.Errorf("invalid elem type for select: %T\n", this.Elem))
	}

	this.Write(&garden, &shoot, &project, &seed)
	return err
}

func (this *select_output) Write(garden, shoot, project, seed *string) {
	env.Warning()
	envout(garden, "GARDEN")
	envout(shoot, "SHOOT")
	envout(project, "PROJECT")
	envout(seed, "SEED")
}

func envout(value *string, key string) {
	if value == nil || *value == "" {
		env.UnSet(fmt.Sprintf("GEX_%s", key))
		fmt.Printf("%-*s cleared\n", 10, key)
	} else {
		env.Set(fmt.Sprintf("GEX_%s", key), *value)
		fmt.Printf("%-*s = \"%v\"\n", 10, key, *value)
	}
}

func Download(envbusy bool, export bool, e gube.KubeconfigProvider, path ...string) error {
	out := filepath.Join(path...)
	os.MkdirAll(out, 0700)
	out = filepath.Join(out, "kubeconfig.yaml")
	cfg, err := e.GetKubeconfig()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return fmt.Errorf("cannot create/write '%s': %s", out, err)
	}
	defer f.Close()
	_, err = f.Write(cfg)
	if err != nil {
		return fmt.Errorf("cannot create/write '%s': %s", out, err)
	}
	fmt.Printf("downloaded to %s\n", out)
	if export {
		if envbusy {
			fmt.Fprintf(os.Stderr, "Warning: not exported to environment, because current garden is determined by KUBECONFIG\n")
		} else {
			env.Set("KUBECONFIG", out)
			fmt.Printf("exported to KUBECONFIG\n")
		}
	}
	return nil
}
