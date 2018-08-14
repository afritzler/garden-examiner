package seed

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/output"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "shell", cmd_shell).
		CmdDescription("run shell on node").CmdArgDescription("[<seed>] <node>")).
		ArgOption(constants.O_NODE).Short('n').ArgDescription("<name>").Description("node name").
		ArgOption(constants.O_POD).Short('p').ArgDescription("<name>").Description("pod name")
}

func cmd_shell(opts *cmdint.Options) error {
	return cmdline.ExecuteOutput(opts, output.NewShellOutput(opts.GetOptionValue(constants.O_NODE), opts.GetOptionValue(constants.O_POD), nil), TypeHandler)
}
