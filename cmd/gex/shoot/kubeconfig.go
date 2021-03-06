package shoot

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/output"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "kubeconfig", kubeconfig).
		CmdDescription("get kubeconfig for shoot").
		CmdArgDescription("[<shoot>]"))
}

func kubeconfig(opts *cmdint.Options) error {
	return cmdline.ExecuteOutput(opts, output.NewKubeconfigOutput(), TypeHandler)
}
