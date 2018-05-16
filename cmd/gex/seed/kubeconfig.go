package seed

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/cmd/gex/verb"
)

func init() {
	filters.AddOptions(verb.Add(GetCmdTab(), "kubeconfig", kubeconfig).
		CmdDescription("get kubeconfig for seed").
		CmdArgDescription("[<seed>]"))
}

func kubeconfig(opts *cmdint.Options) error {
	return util.Doit(opts, NewHandler(util.NewKubeconfigOutput()))
}
