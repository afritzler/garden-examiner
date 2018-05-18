package verb

import (
	"fmt"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/env"
	"github.com/afritzler/garden-examiner/cmd/gex/util"

	"github.com/afritzler/garden-examiner/pkg"
)

func init() {
	NewVerb("select", cmdint.MainTab()).CmdArgDescription("clear|<type> ...").
		CmdDescription("general select command",
			"The first argument is the element type followed by",
			"an optional element name.",
			"With clear the selection can be undone. If nothing is specified",
			"the actual selection is shown",
		).
		DefaultFunction(cmd_select).
		SimpleCommand("clear", cmd_clear).
		CmdArgDescription("{project|seed|shoot}").
		CmdDescription("clear given selection")
}

func cmd_select(opts *cmdint.Options) error {
	found := 0
	if v := opts.GetOptionValue(constants.O_SEL_SHOOT); v != nil {
		fmt.Printf("SHOOT   = %s\n", *v)
		found++
	}
	if v := opts.GetOptionValue(constants.O_SEL_PROJECT); v != nil {
		fmt.Printf("PROJECT = %s\n", *v)
		found++
	}
	if v := opts.GetOptionValue(constants.O_SEL_SEED); v != nil {
		fmt.Printf("SEED    = %s\n", *v)
		found++
	}
	if found == 0 {
		return fmt.Errorf("no selection found")
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////
// clear sub command

type clear_output struct {
	*select_output
}

func (this *clear_output) Out(opts *cmdint.Options) error {
	shoot := opts.GetOptionValue(constants.O_SEL_SHOOT)
	project := opts.GetOptionValue(constants.O_SEL_PROJECT)
	seed := opts.GetOptionValue(constants.O_SEL_SEED)

	if len(opts.Arguments) > 0 {
		for _, n := range opts.Arguments {
			b, d := cmdint.SelectBest(n, "shoot", "seed", "project")
			if d > len(n)/2 {
				return fmt.Errorf("unknown selection type '%s'", n)
			}
			fmt.Printf("clearing %s selection\n", b)
			switch b {
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
	this.Write(shoot, project, seed)
	return nil
}

func cmd_clear(opts *cmdint.Options) error {
	return (&clear_output{select_output: NewSelectOutput()}).Out(opts)
}

////////////////////////////////////////////////////////////////////////////
// general select output

type select_output struct {
	*util.SingleElementOutput
}

var _ util.Output = &select_output{}

func NewSelectOutput() *select_output {
	return &select_output{util.NewSingleElementOutput()}
}

func (this *select_output) Out(ctx *context.Context) error {
	shoot := ""
	seed := ""
	project := ""
	switch e := this.Elem.(type) {
	case gube.Shoot:
		shoot = e.GetName().String()
		project = e.GetName().GetProjectName()
		seed = e.GetSeedName()
	case gube.Seed:
		seed = e.GetName()
	case gube.Project:
		project = e.GetName()
	default:
		panic(fmt.Errorf("invalid elem type for select: %T\n", this.Elem))
	}

	this.Write(&shoot, &project, &seed)
	return nil
}

func (this *select_output) Write(shoot, project, seed *string) {
	env.Warning()
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
