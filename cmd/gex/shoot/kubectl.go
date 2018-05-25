package shoot

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/output"
	"github.com/afritzler/garden-examiner/pkg"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "kubectl", kubectl).Raw().
		CmdDescription("run kubectl for shoot or control plane in seed").
		CmdArgDescription("[--shoot <shoot>] [cp] {<kubectl args/options>}").
		FlagOption("cp").Short('c').ArgDescription("switch to control plane").
		ArgOption("shoot"))
}

func seed_kubectl_mapper(ctx *context.Context, e interface{}) (interface{}, []string, error) {
	ns, err := e.(gube.Shoot).GetNamespaceInSeed()
	if err != nil {
		return nil, nil, err
	}
	seed, _, err := seed_mapper(ctx, e)
	return seed, []string{"-n", ns}, err

}

func kubectl(opts *cmdint.Options) error {
	var mapper output.ElementMapper = nil
	if opts.IsFlag("cp") {
		mapper = seed_kubectl_mapper
	}
	return cmdline.ExecuteOutputRaw("shoot", opts, output.NewKubectlOutput(opts.Arguments, mapper), TypeHandler)
}
