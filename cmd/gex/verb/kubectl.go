package verb

import (
	"fmt"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"

	"github.com/mandelsoft/cmdint/pkg/cmdint"
)

func init() {
	NewVerb("kubectl", cmdint.MainTab()).CmdArgDescription("<type> ...").
		CmdDescription("general kubectl command",
			"The first argument is the element type followed by",
			"element name option and/or kubectl arguments/options.",
			"If no element option is given, it must be defaulted by the",
			"selection command or the appropriate selection option.",
			"If nothing is slected kubectl is run for the garden cluster.",
		).
		CatchUnknownCommand(catch_kubectl).Raw().
		SimpleCommand("garden", cmd_garden).Raw().
		CmdArgDescription("{<kubectl opts/args>}").
		CmdDescription("run kubectl for garden cluster")
}

func catch_kubectl(cmdtab cmdint.CmdTab, opts *cmdint.Options) error {
	fmt.Printf("catching %v\n", opts.Arguments)
	shoot := opts.GetOptionValue(constants.O_SEL_SHOOT)
	if shoot != nil && *shoot != "" {
		return cmdtab.Execute(opts, append([]string{"shoot"}, opts.Arguments...))
	}
	seed := opts.GetOptionValue(constants.O_SEL_SEED)
	if seed != nil && *seed != "" {
		return cmdtab.Execute(opts, append([]string{"seed"}, opts.Arguments...))
	}
	return cmd_garden(opts)
}

func cmd_garden(opts *cmdint.Options) error {
	fmt.Printf("using garden: %v\n", opts.Arguments)
	ctx := context.Get(opts)
	cfg, err := ctx.Garden.GetKubeconfig()
	if err != nil {
		return err
	}
	return util.Kubectl(cfg, append([]string{"-n", "garden"}, opts.Arguments...)...)
}
