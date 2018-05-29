package shoot

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/pkg"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
)

func init() {
	filters.Add(&ProfileFilter{})
}

type ProfileFilter struct {
}

var _ util.Filter = &ProfileFilter{}

func (this *ProfileFilter) AddOptions(cmd cmdint.ConfigurableCmdTabCommand) cmdint.ConfigurableCmdTabCommand {
	return cmd.ArgOption(constants.O_PROFILE)
}

func (this *ProfileFilter) Match(ctx *context.Context, elem interface{}, opts *cmdint.Options) (bool, error) {
	s := elem.(gube.Shoot)
	profile := opts.GetOptionValue(constants.O_PROFILE)

	if profile != nil {
		if s.GetProfileName() != *profile {
			return false, nil
		}
	}
	return true, nil
}
