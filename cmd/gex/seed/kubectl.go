package seed

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/output"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "kubectl", kubectl).Raw().
		CmdDescription("run kubectl for seed").
		CmdArgDescription("[<seed>] {<kubectl args/options>}").
		ArgOption("seed"))
}

func kubectl(opts *cmdint.Options) error {
	return cmdline.ExecuteOutputRaw("seed", opts, output.NewKubectlOutput(opts.Arguments, nil), TypeHandler)
}
