package seed

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/pkg"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
)

func init() {
	filters.Add(&SeedFilter{})
}

type SeedFilter struct {
}

var _ util.Filter = &SeedFilter{}

func (this *SeedFilter) AddOptions(cmd cmdint.ConfigurableCmdTabCommand) cmdint.ConfigurableCmdTabCommand {
	return cmd.ArgOption(constants.O_SEED)
}

func (this *SeedFilter) Match(ctx *context.Context, elem interface{}, opts *cmdint.Options) (bool, error) {
	s := elem.(gube.Seed)
	seed := opts.GetOptionValue(constants.O_SEED)
	if seed != nil {
		if s.GetName() != *seed {
			return false, nil
		}
	}
	return true, nil
}
