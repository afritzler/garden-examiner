package seed

import (
	"fmt"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/shoot"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/pkg"
	"github.com/mandelsoft/cmdint/pkg/cmdint"
)

var cmdtab cmdint.ConfigurableCmdTab = cmdint.NewCmdTab("seed", nil).
	CmdDescription("garden seed clusters\n" +
		"list one or more seed clusters").
	CmdArgDescription("<command>")

func init() {
	cmdint.MainTab().
		Command("seed", cmdtab)
}

func GetCmdTab() cmdint.ConfigurableCmdTab {
	return cmdtab
}

/////////////////////////////////////////////////////////////////////////////

var filters *util.Filters = util.NewFilters()

/////////////////////////////////////////////////////////////////////////////

type Handler struct {
	output util.Output
	seeds  *SeedAccess
	shoots *shoot.Handler
}

func NewHandler(o util.Output) *Handler {
	return &Handler{o, nil, nil}
}

func (this *Handler) RequireScan(name string) bool {
	return false
}

func (this *Handler) MatchName(e interface{}, name string) (bool, error) {
	s := e.(gube.Seed)
	return s.GetName() == name, nil
}

func (this *Handler) Get(ctx *context.Context, name string) (interface{}, error) {
	if this.seeds == nil {
		return ctx.Garden.GetSeed(name)
	}
	s, ok := this.seeds.data[name]
	if !ok {
		return nil, fmt.Errorf("seed '%s' not found", name)
	}
	return s, nil
}

func (this *Handler) Iterator(ctx *context.Context, opts *cmdint.Options) (util.Iterator, error) {
	if this.seeds == nil {
		seeds, err := ctx.Garden.GetSeeds()
		if err != nil {
			return nil, err
		}
		this.seeds = NewSeedAccess(seeds)
	}
	return util.NewIndexedIterator(this.seeds), nil
}

func (this *Handler) Match(ctx *context.Context, e interface{}, opts *cmdint.Options) (bool, error) {
	s := e.(gube.Seed)
	project := opts.GetOptionValue(constants.O_PROJECT)
	seed := opts.GetOptionValue(constants.O_SEED)

	if project != nil {
		if this.shoots == nil {
			this.shoots = shoot.NewHandler(nil)
		}
		i, err := this.shoots.Iterator(ctx, opts)
		if err != nil {
			return false, err
		}
		for i.HasNext() {
			shoot := i.Next().(gube.Shoot)
			if shoot.GetSeedName() == s.GetName() {
				return true, nil
			}
		}
	}
	if seed != nil {
		if s.GetName() != *seed {
			return false, nil
		}
	}
	return true, nil
}

func (this *Handler) Add(ctx *context.Context, e interface{}) error {
	return this.output.Add(ctx, e)
}

func (this *Handler) Out(ctx *context.Context) {
	this.output.Out(ctx)
}

type SeedAccess struct {
	data  map[string]gube.Seed
	slice []interface{}
}

func NewSeedAccess(data map[string]gube.Seed) *SeedAccess {
	a := &SeedAccess{data: data, slice: make([]interface{}, len(data))}
	i := 0
	for _, v := range data {
		a.slice[i] = v
		i++
	}
	return a
}

func (this *SeedAccess) Size() int {
	return len(this.data)
}

func (this *SeedAccess) Get(i int) interface{} {
	return this.slice[i]
}
