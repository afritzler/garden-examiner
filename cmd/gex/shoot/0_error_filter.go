package shoot

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/pkg"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
)

func init() {
	filters.Add(&ErrorFilter{})
}

type ErrorFilter struct {
}

var _ util.Filter = &ErrorFilter{}

func (this *ErrorFilter) AddOptions(cmd cmdint.ConfigurableCmdTabCommand) cmdint.ConfigurableCmdTabCommand {
	r := cmd.FlagOption(constants.O_ERROR)
	return r
}

func (this *ErrorFilter) Match(ctx *context.Context, elem interface{}, opts *cmdint.Options) (bool, error) {
	s := elem.(gube.Shoot)
	flag := opts.IsFlag(constants.O_ERROR)

	if flag {
		if s.GetError() == "" {
			return false, nil
		}
	}
	return true, nil
}
