package verb

import (
	"fmt"

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
		CatchUnknownCommand(catch_cluster).Raw().
		SimpleCommand("garden", cmd_kubectl_garden).Raw().
		CmdArgDescription("{<kubectl opts/args>}").
		CmdDescription("run kubectl for garden cluster")
}

func cmd_kubectl_garden(opts *cmdint.Options) error {
	fmt.Printf("using garden: %v\n", opts.Arguments)
	ctx := context.Get(opts)
	cfg, err := ctx.Garden.GetKubeconfig()
	if err != nil {
		return err
	}
	return util.Kubectl(cfg, nil, append([]string{"-n", "garden"}, opts.Arguments...)...)
}
