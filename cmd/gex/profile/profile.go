package profile

import (
	"fmt"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/pkg"
	"github.com/mandelsoft/cmdint/pkg/cmdint"
)

var cmdtab cmdint.ConfigurableCmdTab = cmdint.NewCmdTab("profile", nil).
	CmdDescription("garden profiles\n" +
		"list one or more profiles").
	CmdArgDescription("<command>")

func init() {
	cmdint.MainTab().
		Command("profile", cmdtab)
}

func GetCmdTab() cmdint.ConfigurableCmdTab {
	return cmdtab
}

/////////////////////////////////////////////////////////////////////////////

var filters *util.Filters = util.NewFilters()

/////////////////////////////////////////////////////////////////////////////

type _TypeHandler struct {
	data map[string]gube.Profile
}

var TypeHandler cmdline.ElementTypeHandler = &_TypeHandler{}

func (this *_TypeHandler) GetDefault(opts *cmdint.Options) *string {
	return nil
}

func (this *_TypeHandler) GetAll(ctx *context.Context, opts *cmdint.Options) ([]interface{}, error) {
	elems, err := ctx.Garden.GetProfiles()
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
	s := e.(gube.Profile)
	return s.GetName() == name, nil
}

func (this *_TypeHandler) Get(ctx *context.Context, name string) (interface{}, error) {
	if this.data == nil {
		return ctx.Garden.GetProfile(name)
	}
	s, ok := this.data[name]
	if !ok {
		return nil, fmt.Errorf("profile '%s' not found", name)
	}
	return s, nil
}
