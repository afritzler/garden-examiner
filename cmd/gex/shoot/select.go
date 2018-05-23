package shoot

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/cmd/gex/verb"
)

func init() {
	filters.AddOptions(verb.Add(GetCmdTab(), "select", cmd_select).
		CmdDescription("select shoot cluster").CmdArgDescription("<shoot>"))

}

func cmd_select(opts *cmdint.Options) error {
	h, err := NewSelectHandler(opts)
	if err != nil {
		return err
	}
	return util.Doit(opts, h)
}

func NewSelectHandler(opts *cmdint.Options) (util.Handler, error) {
	return NewHandler(verb.NewSelectOutput()), nil
}
