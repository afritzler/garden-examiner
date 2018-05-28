package shoot

import (
	"fmt"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/output"
	"github.com/afritzler/garden-examiner/pkg"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "shell", cmd_shell).
		CmdDescription("run shell on node").CmdArgDescription("[<shoot>] <node>")).
		ArgOption(constants.O_NODE).Short('n').ArgDescription("node name").
		FlagOption("cp").Short('c').ArgDescription("switch to control plane")
}

func seed_mapper(ctx *context.Context, e interface{}) (interface{}, []string, error) {
	shoot := e.(gube.Shoot)

	seed, err := ctx.Garden.GetSeed(shoot.GetSeedName())
	if err != nil {
		return nil, nil, err
	}
	return seed, nil, nil
}

func cmd_shell(opts *cmdint.Options) error {
	var mapper output.ElementMapper = nil
	if opts.IsFlag("cp") {
		mapper = seed_mapper
	}
	fmt.Printf("opts: %#v\n", opts)
	return cmdline.ExecuteOutput(opts, output.NewShellOutput(opts.GetOptionValue(constants.O_NODE), mapper), TypeHandler)
}
