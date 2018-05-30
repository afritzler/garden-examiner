package garden

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/verb"
)

func init() {
	filters.AddOptions(verb.Add(GetCmdTab(), "select", cmd_select).
		CmdDescription("select garden").CmdArgDescription("<garden from config>").
		FlagOption(constants.O_NOFILTER).Description("ignore filters and selection"))

}

func cmd_select(opts *cmdint.Options) error {
	return cmdline.ExecuteOutput(opts, verb.NewSelectOutput(), TypeHandler)
}
