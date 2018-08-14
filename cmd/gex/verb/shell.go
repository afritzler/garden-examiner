package verb

import (
	"fmt"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/output"

	"github.com/mandelsoft/cmdint/pkg/cmdint"
)

func init() {
	NewVerb("shell", cmdint.MainTab()).CmdArgDescription("<type> ...").
		CmdDescription("run shell on cluster node",
			"The first argument is the element type followed by",
			"element name option and/or kubectl arguments/options.",
			"If no element option is given, it must be defaulted by the",
			"selection command or the appropriate selection option.",
			"If nothing is selected shell is run for the garden cluster.",
		).
		ArgOption(constants.O_NODE).Short('n').ArgDescription("<name>").Description("node name").
		ArgOption(constants.O_POD).Short('p').ArgDescription("<name>").Description("pod name").
		CatchUnknownCommand(catch_cluster).
		SimpleCommand("garden", cmd_shell_garden).
		ArgOption(constants.O_NODE).Short('n').ArgDescription("<name>").Description("node name").
		ArgOption(constants.O_POD).Short('p').ArgDescription("<name>").Description("pod name").
		CmdDescription("run shell for garden cluster")
}

func cmd_shell_garden(opts *cmdint.Options) error {
	fmt.Printf("using garden: %v\n", opts.Arguments)
	ctx := context.Get(opts)
	out := output.NewShellOutput(opts.GetOptionValue(constants.O_NODE), opts.GetOptionValue(constants.O_POD), nil)
	out.Add(ctx, ctx.Garden)
	return out.Out(ctx)
}
