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
		ArgOption("shoot"))
}

func seed_mapper(ctx *context.Context, e interface{}) (interface{}, []string, error) {
	shoot := e.(gube.Shoot)

	seed, err := ctx.Garden.GetSeed(shoot.GetSeedName())
	if err != nil {
		return nil, nil, err
	}
	ns, err := shoot.GetNamespaceInSeed()
	if err != nil {
		return nil, nil, err
	}
	return seed, []string{"-n", ns}, nil

}

func kubectl(opts *cmdint.Options) error {

	if len(opts.Arguments) > 0 && opts.Arguments[0] == "cp" {
		return util.ExecuteOutputRaw("shoot", opts, util.NewKubectlOutput(opts.Arguments[1:], seed_mapper), TypeHandler)
	}
	return util.ExecuteOutputRaw("shoot", opts, util.NewKubectlOutput(opts.Arguments, nil), TypeHandler)
}
