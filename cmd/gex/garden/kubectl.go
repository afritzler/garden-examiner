package garden

import (
	"fmt"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "kubectl", kubectl).Raw().
		CmdDescription("run kubectl for garden").
		CmdArgDescription("[<garden>] {<kubectl args/options>}").
		ArgOption("garden"))
}

// func kubectl(opts *cmdint.Options) error {
// 	return cmdline.ExecuteOutputRaw("garden", opts, output.NewKubectlOutput(opts.Arguments, nil), TypeHandler)
// }

func kubectl(opts *cmdint.Options) error {
	fmt.Printf("using garden: %v\n", opts.Arguments)
	ctx := context.Get(opts)
	cfg, err := ctx.Garden.GetKubeconfig()
	if err != nil {
		return err
	}
	return util.Kubectl(cfg, nil, append([]string{"-n", "garden"}, opts.Arguments...)...)
}
