package shoot

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/pkg"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
)

var cmdtab cmdint.ConfigurableCmdTab = cmdint.NewCmdTab("shoot", nil).
	CmdDescription("garden shoot clusters\n" +
		"list one or more shoot clusters").
	CmdArgDescription("<command>").
	ArgOption(constants.O_PROJECT).
	ArgOption(constants.O_SEED)

func init() {
	cmdint.MainTab().
		Command("shoot", cmdtab)
}

func GetCmdTab() cmdint.ConfigurableCmdTab {
	return cmdtab
}

/////////////////////////////////////////////////////////////////////////////

var filters *util.Filters = util.NewFilters()

/////////////////////////////////////////////////////////////////////////////

type _TypeHandler struct {
	data map[gube.ShootName]gube.Shoot
}

var TypeHandler = &_TypeHandler{}

func (this *_TypeHandler) GetAll(ctx *context.Context, opts *cmdint.Options) ([]interface{}, error) {
	elems, err := ctx.Garden.GetShoots()
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
func (this *_TypeHandler) GetDefault(opts *cmdint.Options) *string {
	return opts.LookupOptionValue(constants.O_SEL_SHOOT)
}
func (this *_TypeHandler) RequireScan(name string) bool {
	i := strings.Index(name, "/")
	return i < 0
}
func (this *_TypeHandler) MatchName(e interface{}, name string) (bool, error) {
	s := e.(gube.Shoot)
	return s.GetName().GetName() == name, nil
}
func (this *_TypeHandler) Get(ctx *context.Context, name string) (interface{}, error) {
	i := strings.Index(name, "/")
	sn := gube.NewShootName(string(name[:i]), string(name[i+1:]))
	if this.data == nil {
		//fmt.Printf("use garden %p\n", ctx.Garden)
		return ctx.Garden.GetShoot(sn)
	}
	s, ok := this.data[*sn]
	if !ok {
		return nil, fmt.Errorf("shoot '%s' not found", sn)
	}
	return s, nil
}
