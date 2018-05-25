package seed

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/cmd/gex/verb"
)

func init() {
	filters.AddOptions(verb.Add(GetCmdTab(), "shell", cmd_shell).
		CmdDescription("run shell on node").CmdArgDescription("[<seed>] <node>")).
		ArgOption(constants.O_NODE).Short('n').ArgDescription("node name")
}

func cmd_shell(opts *cmdint.Options) error {
	return util.ExecuteOutput(opts, util.NewShellOutput(opts.GetOptionValue(constants.O_NODE), nil), TypeHandler)
}
