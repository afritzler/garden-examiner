package garden

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/output"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "kubectl", kubectl).Raw().
		CmdDescription("run kubectl for garden").
		CmdArgDescription("[<garden>] {<kubectl args/options>}").
		ArgOption("garden"))
}

func kubectl(opts *cmdint.Options) error {
	return cmdline.ExecuteOutputRaw("garden", opts, output.NewKubectlOutput(append([]string{"-n", "garden"}, opts.Arguments...), nil), TypeHandler)
}
