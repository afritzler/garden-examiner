package seed

import (
	"fmt"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/cmd/gex/verb"
	"github.com/afritzler/garden-examiner/pkg"
	"github.com/afritzler/garden-examiner/pkg/data"
)

func init() {
	filters.AddOptions(verb.Add(GetCmdTab(), "get", get).CmdDescription("get seed(s)").CmdArgDescription("[<seed>]")).
		ArgOption(constants.O_OUTPUT).Short('o').
		ArgOption(constants.O_SORT).Array()
}

func get(opts *cmdint.Options) error {
	return util.ExecuteMode(opts, get_outputs, TypeHandler)
}

/////////////////////////////////////////////////////////////////////////////

var get_outputs = util.NewOutputs(get_regular, util.Outputs{
	"kubeconfig": util.KubeconfigOutputFactory,
}).AddManifestOutputs()

func get_regular(opts *cmdint.Options) util.Output {
	return util.NewProcessingTableOutput(opts, data.Chain().Map(map_get_regular_output),
		"SEED", "INFRA", "REGION", "PROFILE", "SHOOT", "STATE", "ERROR")
}

func map_get_regular_output(e interface{}) interface{} {
	s := e.(gube.Seed)
	c := s.GetCloud()
	p, err := s.Garden().GetProfile(c.Profile)
	i := "unknown"
	if err == nil {
		i = p.GetInfrastructure()
	}
	shoot := ""
	state := ""
	msg := ""
	sn := s.GetShoot()
	if sn != nil {
		shoot = sn.GetName()
		sh, err := s.Garden().GetShoot(sn)
		if err != nil {
			state = fmt.Sprintf("%s", err)
		} else {
			state = sh.GetState()
			msg = sh.GetError()

		}
	}
	return []string{s.GetName(), i, c.Region, c.Profile, shoot, state, util.Oneline(msg, 90)}
}
