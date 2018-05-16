package cmd_select

import (
	"fmt"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
)

func init() {
	GetCmdTab().SimpleCommand("clear", clear).CmdDescription("clear selection").CmdArgDescription("{project|seed|shoot}")
}

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

func clear(opts *cmdint.Options) error {
	return (&clear_output{select_output: NewSelectOutput()}).Out(opts)
}
