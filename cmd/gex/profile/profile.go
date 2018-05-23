package profile

import (
	"fmt"

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

type Handler struct {
	*util.BasicSelfHandler
	data map[string]gube.Profile
}

func NewHandler(o util.Output) *Handler {
	h := &Handler{}
	h.BasicSelfHandler = util.NewBasicSelfHandler(o, h)
	return h
}

func NewModeHandler(opts *cmdint.Options, o util.Outputs) (*Handler, error) {
	h := &Handler{}
	b, err := util.NewBasicModeSelfHandler(opts, o, h)
	if err != nil {
		return nil, err
	}
	h.BasicSelfHandler = b
	return h, nil
}

func (this *Handler) Doit(opts *cmdint.Options) error {
	return util.Doit(opts, this)
}

func (this *Handler) GetAll(ctx *context.Context, opts *cmdint.Options) ([]interface{}, error) {
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

func (this *Handler) GetFilter() util.Filter {
	return filters
}

func (this *Handler) MatchName(e interface{}, name string) (bool, error) {
	s := e.(gube.Profile)
	return s.GetName() == name, nil
}

func (this *Handler) Get(ctx *context.Context, name string) (interface{}, error) {
	if this.data == nil {
		return ctx.Garden.GetProfile(name)
	}
	s, ok := this.data[name]
	if !ok {
		return nil, fmt.Errorf("profile '%s' not found", name)
	}
	return s, nil
}
