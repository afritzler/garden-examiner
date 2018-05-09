package shoot

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/pkg"
)

func init() {
	filters.AddOptions(GetCmdTab().SimpleCommand("kubeconfig", kubeconfig).
		CmdDescription("get kubeconfig for shoot").
		CmdArgDescription("[<shoot>]"))
}

func kubeconfig(opts *cmdint.Options) error {
	return util.Doit(opts, NewKubeConfigHandler())
}

/////////////////////////////////////////////////////////////////////////////

type kubeconfig_output struct {
	data []string
}

func (this *kubeconfig_output) Add(ctx *context.Context, e interface{}) error {
	s := e.(gube.Shoot)
	cfg, err := s.GetKubeconfig()
	if err != nil {
		return err
	}
	this.data = append(this.data, string(cfg))
	return nil
}

func (this *kubeconfig_output) Out(ctx *context.Context) {
	for _, cfg := range this.data {
		if !strings.HasPrefix(cfg, "---\n") {
			fmt.Println("---")
		}
		fmt.Println(cfg)
	}
}

type KubeConfigHandler struct {
	*Handler
}

func NewKubeConfigHandler() util.Handler {
	return &KubeConfigHandler{&Handler{&kubeconfig_output{[]string{}}, nil}}
}
