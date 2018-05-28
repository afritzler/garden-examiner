package seed

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/output"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "iaas", iaas).Raw().
		CmdDescription("run iaas specific cmd for seed or control plane in seed").
		CmdArgDescription("[--seed <seed>] [cp] {<iaas args/options>}").
		FlagOption("cp").Short('c').ArgDescription("switch to control plane").
		ArgOption("seed"))
}

func iaas(opts *cmdint.Options) error {
	return cmdline.ExecuteOutputRaw("seed", opts, output.NewIaasOutput(opts.Arguments, nil), TypeHandler)
}
