package seed

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/output"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "iaas", iaas).Raw().
		CmdDescription("run iaas specific cmd for seed or control plane in seed").
		CmdArgDescription("[--seed <seed>] [cp] {<iaas args/options>}").
		FlagOption("cp").Short('c').ArgDescription("switch to control plane").
		FlagOption(constants.O_EXPORT).Short('e').Description("set CLI environment").
		ArgOption("seed"))
}

func iaas(opts *cmdint.Options) error {
	var o output.Output
	if opts.IsFlag(constants.O_EXPORT) {
		o = output.NewIaasExportOutput(opts.Arguments, nil, context.Get(opts).CacheDirFor)
	} else {
		o = output.NewIaasOutput(opts.Arguments, nil)
	}
	return cmdline.ExecuteOutputRaw("seed", opts, o, TypeHandler)
}
