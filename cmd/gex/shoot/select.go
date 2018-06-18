package shoot

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/verb"
)

func init() {
	filters.AddOptions(verb.Add(GetCmdTab(), "select", cmd_select).
		CmdDescription("select shoot cluster").CmdArgDescription("<shoot>").
		FlagOption(constants.O_DOWNLOAD).Short('d').Description("download kubeconfig").
		FlagOption(constants.O_EXPORT).Short('e').Description("export env KUBECONFIG (implies -d)"))

}

func cmd_select(opts *cmdint.Options) error {
	return cmdline.ExecuteOutput(opts, verb.NewSelectOutput(opts.IsFlag(constants.O_DOWNLOAD), opts.IsFlag(constants.O_EXPORT)), TypeHandler)
}
