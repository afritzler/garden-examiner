package shoot

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/cmd/gex/verb"
	"github.com/afritzler/garden-examiner/pkg"
)

func init() {
	filters.AddOptions(verb.Add(GetCmdTab(), "kubectl", kubectl).Raw().
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
	var mapper util.ElementMapper = nil
	if opts.IsFlag("cp") {
		mapper = seed_kubectl_mapper
	}
	return util.ExecuteOutputRaw("shoot", opts, util.NewKubectlOutput(opts.Arguments, mapper), TypeHandler)
}
