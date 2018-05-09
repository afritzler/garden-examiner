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

type Handler struct {
	output util.Output
	shoots *ShootAccess
}

func NewHandler(o util.Output) *Handler {
	return &Handler{o, nil}
}

func (this *Handler) RequireScan(name string) bool {
	i := strings.Index(name, "/")
	return i < 0
}

func (this *Handler) MatchName(e interface{}, name string) (bool, error) {
	s := e.(gube.Shoot)
	return s.GetName().GetName() == name, nil
}

func (this *Handler) Get(ctx *context.Context, name string) (interface{}, error) {
	i := strings.Index(name, "/")
	sn := gube.NewShootName(string(name[:i]), string(name[i+1]))
	if this.shoots == nil {
		return ctx.Garden.GetShoot(sn)
	}
	s, ok := this.shoots.data[*sn]
	if !ok {
		return nil, fmt.Errorf("shoot '%s' not found", sn)
	}
	return s, nil
}

func (this *Handler) Iterator(ctx *context.Context, opts *cmdint.Options) (util.Iterator, error) {
	if this.shoots == nil {
		shoots, err := ctx.Garden.GetShoots()
		if err != nil {
			return nil, err
		}
		this.shoots = NewShootAccess(shoots)
	}
	return util.NewIndexedIterator(this.shoots), nil
}

func (this *Handler) Match(ctx *context.Context, e interface{}, opts *cmdint.Options) (bool, error) {
	return filters.Match(ctx, e, opts)
}

func (this *Handler) Add(ctx *context.Context, e interface{}) error {
	return this.output.Add(ctx, e)
}

func (this *Handler) Out(ctx *context.Context) {
	this.output.Out(ctx)
}

type ShootAccess struct {
	data  map[gube.ShootName]gube.Shoot
	slice []interface{}
}

func NewShootAccess(data map[gube.ShootName]gube.Shoot) *ShootAccess {
	a := &ShootAccess{data: data, slice: make([]interface{}, len(data))}
	i := 0
	for _, v := range data {
		a.slice[i] = v
		i++
	}
	return a
}

func (this *ShootAccess) Size() int {
	return len(this.data)
}

func (this *ShootAccess) Get(i int) interface{} {
	return this.slice[i]
}
