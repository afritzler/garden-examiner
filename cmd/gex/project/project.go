package project

import (
	"fmt"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"

	"github.com/afritzler/garden-examiner/pkg"
	"github.com/mandelsoft/cmdint/pkg/cmdint"
)

var cmdtab cmdint.ConfigurableCmdTab = cmdint.NewCmdTab("project", nil).
	CmdDescription("garden projects\n" +
		"list one or more projects").
	CmdArgDescription("<command>")

func init() {
	cmdint.MainTab().
		Command("project", cmdtab)
}

func GetCmdTab() cmdint.ConfigurableCmdTab {
	return cmdtab
}

/////////////////////////////////////////////////////////////////////////////

var filters *util.Filters = util.NewFilters()

/////////////////////////////////////////////////////////////////////////////

type _TypeHandler struct {
	data map[string]gube.Project
}

var TypeHandler util.ElementTypeHandler = &_TypeHandler{}

func (this *_TypeHandler) GetDefault(opts *cmdint.Options) *string {
	shoot := opts.GetOptionValue(constants.O_SEL_PROJECT)
	return shoot
}

func (this *_TypeHandler) GetAll(ctx *context.Context, opts *cmdint.Options) ([]interface{}, error) {
	elems, err := ctx.Garden.GetProjects()
	if err != nil {
		return nil, err
	}

	this.data = elems
	a := make([]interface{}, len(elems))
	i := 0
	for _, v := range elems {
		a[i] = v
		i++
	}
	return a, nil
}

func (this *_TypeHandler) GetFilter() util.Filter {
	return filters
}
func (this *_TypeHandler) RequireScan(name string) bool {
	return false
}

func (this *_TypeHandler) MatchName(e interface{}, name string) (bool, error) {
	s := e.(gube.Project)
	return s.GetName() == name, nil
}

func (this *_TypeHandler) Get(ctx *context.Context, name string) (interface{}, error) {
	if this.data == nil {
		return ctx.Garden.GetProject(name)
	}
	s, ok := this.data[name]
	if !ok {
		return nil, fmt.Errorf("project '%s' not found", name)
	}
	return s, nil
}
