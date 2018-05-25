package shoot

import (
	"fmt"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/cmd/gex/verb"
	"github.com/afritzler/garden-examiner/pkg"
)

func init() {
	filters.AddOptions(verb.Add(GetCmdTab(), "shell", cmd_shell).
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
	var mapper util.ElementMapper = nil
	if opts.IsFlag("cp") {
		mapper = seed_mapper
	}
	fmt.Printf("opts: %#v\n", opts)
	return util.ExecuteOutput(opts, util.NewShellOutput(opts.GetOptionValue(constants.O_NODE), mapper), TypeHandler)
}
