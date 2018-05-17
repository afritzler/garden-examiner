package seed

import (
	"fmt"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/cmd/gex/verb"
)

func init() {
	filters.AddOptions(verb.Add(GetCmdTab(), "kubectl", kubectl).Raw().
		CmdDescription("run kubectl for seed").
		CmdArgDescription("[<seed>] {<kubectl args/options>}").
		ArgOption("seed"))
}

func kubectl(opts *cmdint.Options) error {
	fmt.Printf("INITIAL: %v\n", opts.Arguments)
	return util.DoitRaw("seed", opts, NewHandler(util.NewKubectlOutput(opts.Arguments, nil)))
}
