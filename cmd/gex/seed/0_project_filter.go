package seed

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
	s := elem.(gube.Seed)
	project := opts.GetOptionValue(constants.O_PROJECT)

	if project != nil {
		shoots, err := ctx.Garden.GetShoots()
		if err != nil {
			return false, err
		}
		for n, sh := range shoots {
			if n.GetProjectName() == *project {
				if sh.GetSeedName() == s.GetName() {
					return true, nil
				}
			}
		}
	}
	return true, nil
}
