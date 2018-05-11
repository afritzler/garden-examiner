package shoot

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/util"
)

func init() {
	filters.AddOptions(GetCmdTab().SimpleCommand("kubeconfig", kubeconfig).
		CmdDescription("get kubeconfig for shoot").
		CmdArgDescription("[<shoot>]"))
}

func kubeconfig(opts *cmdint.Options) error {
	return util.Doit(opts, NewHandler(util.NewKubeconfigOutput()))
}
