package shoot

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/pkg"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
)

func init() {
	filters.Add(&ProjectFilter{})
}

type ProjectFilter struct {
}

var _ util.Filter = &ProjectFilter{}

func (this *ProjectFilter) AddOptions(cmd cmdint.ConfigurableCmdTabCommand) cmdint.ConfigurableCmdTabCommand {
	return cmd.ArgOption(constants.O_PROJECT).Context(constants.O_SEL_PROJECT)
}

func (this *ProjectFilter) Match(ctx *context.Context, elem interface{}, opts *cmdint.Options) (bool, error) {
	s := elem.(gube.Shoot)
	project := opts.GetOptionValue(constants.O_PROJECT)

	if project != nil {
		p, err := s.GetProject()
		if err != nil || p.GetName() != *project {
			return false, nil
		}
	}
	return true, nil
}
