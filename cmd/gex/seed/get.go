package seed

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/cmd/gex/verb"
	"github.com/afritzler/garden-examiner/pkg"
)

func init() {
	filters.AddOptions(verb.Add(GetCmdTab(), "get", get).CmdDescription("get seed(s)").CmdArgDescription("[<seed>]")).
		ArgOption(constants.O_OUTPUT).Short('o')
}

func get(opts *cmdint.Options) error {
	h, err := NewGetHandler(opts)
	if err != nil {
		return err
	}
	return util.Doit(opts, h)
}

/////////////////////////////////////////////////////////////////////////////

type get_output struct {
	*util.TableOutput
}

func (this *get_output) Add(ctx *context.Context, e interface{}) error {
	s := e.(gube.Seed)
	c := s.GetCloud()
	p, err := ctx.GetProfile(c.Profile)
	i := "unknown"
	if err == nil {
		i = p.GetInfrastructure()
	}
	shoot := ""
	sn := s.GetShoot()
	if sn != nil {
		shoot = sn.GetName()
	}
	this.AddLine(
		[]string{s.GetName(), i, c.Region, c.Profile, shoot},
	)
	return nil
}

type GetHandler struct {
	*Handler
}

func NewGetHandler(opts *cmdint.Options) (util.Handler, error) {

	o, err := util.GetOutput(opts, &get_output{
		util.NewTableOutput([][]string{
			[]string{"Seed", "Infra", "Region", "Profile", "Shoot"},
		}),
	})
	if err != nil {
		return nil, err
	}
	return &GetHandler{NewHandler(o)}, nil
}
