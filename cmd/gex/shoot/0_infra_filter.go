package shoot

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/pkg"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
)

func init() {
	filters.Add(&InfraFilter{})
}

type InfraFilter struct {
}

var _ util.Filter = &InfraFilter{}

func (this *InfraFilter) AddOptions(cmd cmdint.ConfigurableCmdTabCommand) cmdint.ConfigurableCmdTabCommand {
	return cmd.ArgOption(constants.O_INFRA)
}

func (this *InfraFilter) Match(ctx *context.Context, elem interface{}, opts *cmdint.Options) (bool, error) {
	s := elem.(gube.Shoot)
	infra := opts.GetOptionValue(constants.O_INFRA)

	if infra != nil {
		if s.GetInfrastructure() != *infra {
			return false, nil
		}
	}
	return true, nil
}
