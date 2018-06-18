package shoot

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/output"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "iaas", cmd_iaas).Raw().
		CmdDescription("run iaas specific cmd for shoot or control plane in seed").
		CmdArgDescription("[--shoot <shoot>] [cp] {<iaas args/options>}").
		FlagOption("cp").Short('c').ArgDescription("switch to control plane").
		FlagOption(constants.O_EXPORT).Short('e').Description("set CLI environment").
		ArgOption("shoot"))
}

func cmd_iaas(opts *cmdint.Options) error {
	var mapper output.ElementMapper = nil
	if opts.IsFlag("cp") {
		mapper = seed_mapper
	}
	var o output.Output
	if opts.IsFlag(constants.O_EXPORT) {
		o = output.NewIaasExportOutput(opts.Arguments, mapper, context.Get(opts).CacheDirFor)
	} else {
		o = output.NewIaasOutput(opts.Arguments, mapper)
	}
	return cmdline.ExecuteOutputRaw("shoot", opts, o, TypeHandler)
}
