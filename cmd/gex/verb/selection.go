package verb

import (
	"fmt"

	"github.com/afritzler/garden-examiner/cmd/gex/const"

	"github.com/mandelsoft/cmdint/pkg/cmdint"
)

func catch_cluster(cmdtab cmdint.CmdTab, opts *cmdint.Options) error {
	fmt.Printf("catching %v\n", opts.Arguments)
	shoot := opts.GetOptionValue(constants.O_SEL_SHOOT)
	if shoot != nil && *shoot != "" {
		return cmdtab.Execute(opts, append([]string{"shoot"}, opts.Arguments...))
	}
	seed := opts.GetOptionValue(constants.O_SEL_SEED)
	if seed != nil && *seed != "" {
		return cmdtab.Execute(opts, append([]string{"seed"}, opts.Arguments...))
	}
	if cmdtab.GetCommand("garden") != nil {
		return cmdtab.Execute(opts, append([]string{"garden"}, opts.Arguments...))
	}
	return fmt.Errorf("No target selected")
}
