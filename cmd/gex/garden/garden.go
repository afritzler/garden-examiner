package garden

import (
	"fmt"

	"github.com/afritzler/garden-examiner/pkg"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
)

var cmdtab cmdint.ConfigurableCmdTab = cmdint.NewCmdTab("garden", nil).
	CmdDescription("all garden clusters\n" +
		"list one or more garden clusters").
	CmdArgDescription("<command>").
	ArgOption(constants.O_PROJECT).
	ArgOption(constants.O_GARDEN)

func init() {
	cmdint.MainTab().
		Command("garden", cmdtab)
}

func GetCmdTab() cmdint.ConfigurableCmdTab {
	return cmdtab
}

/////////////////////////////////////////////////////////////////////////////

var filters *util.Filters = util.NewFilters()

/////////////////////////////////////////////////////////////////////////////

type _TypeHandler struct {
	data map[string]gube.GardenConfig
}

var TypeHandler = &_TypeHandler{}

func (this *_TypeHandler) GetDefault(opts *cmdint.Options) *string {
	return opts.GetOptionValue(constants.O_SEL_GARDEN)
}

func (this *_TypeHandler) GetAll(ctx *context.Context, opts *cmdint.Options) ([]interface{}, error) {
	this.data = ctx.GardenSetConfig.GetConfigs()
	a := make([]interface{}, len(this.data))
	i := 0
	for _, v := range this.data {
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
	g := e.(gube.GardenConfig)
	return g.GetName() == name, nil
}

func (this *_TypeHandler) Get(ctx *context.Context, name string) (interface{}, error) {
	if this.data == nil {
		return ctx.GardenSetConfig.GetConfig(name)
	}
	g, ok := this.data[name]
	if !ok {
		return nil, fmt.Errorf("garden '%s' not found", name)
	}
	return g, nil
}
