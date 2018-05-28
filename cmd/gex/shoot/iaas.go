package shoot

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/output"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "iaas", cmd_iaas).Raw().
		CmdDescription("run iaas specific cmd for shoot or control plane in seed").
		CmdArgDescription("[--shoot <shoot>] [cp] {<iaas args/options>}").
		FlagOption("cp").Short('c').ArgDescription("switch to control plane").
		ArgOption("shoot"))
}

func cmd_iaas(opts *cmdint.Options) error {
	var mapper output.ElementMapper = nil
	if opts.IsFlag("cp") {
		mapper = seed_mapper
	}
	return cmdline.ExecuteOutputRaw("shoot", opts, output.NewIaasOutput(opts.Arguments, mapper), TypeHandler)
}
