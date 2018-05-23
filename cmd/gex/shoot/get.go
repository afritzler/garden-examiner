package shoot

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
	filters.AddOptions(verb.Add(GetCmdTab(), "get", get).CmdDescription("get shoot(s)").CmdArgDescription("[<shoot>]")).
		ArgOption(constants.O_OUTPUT).Short('o')
}

func get(opts *cmdint.Options) error {
	h, err := NewGetHandler(opts)
	if err != nil {
		return err
	}
	return util.Doit(opts, h)
}

/////////////////////////////////////////////////////////////////////////////

var get_outputs = util.NewOutputs(get_regular, util.Outputs{
	"wide":       get_wide,
	"kubeconfig": util.KubeconfigOutputFactory,
}).AddManifestOutputs()

func get_regular(opts *cmdint.Options) util.Output {
	return util.NewProcessingTableOutput(data.Chain().Map(map_get_regular_output),
		"SHOOT", "PROJECT", "INFRA", "SEED", "STATE", "ERROR")
}
func get_wide(opts *cmdint.Options) util.Output {
	return util.NewProcessingTableOutput(data.Chain().Parallel(20).Map(map_get_wide_output),
		"SHOOT", "PROJECT", "INFRA", "SEED", "NODES", "STATE", "ERROR")
}

/////////////////////////////////////////////////////////////////////////////

func map_get_regular_output(e interface{}) interface{} {
	s := e.(gube.Shoot)
	return []string{s.GetName().GetName(), s.GetName().GetProjectName(),
		s.GetInfrastructure(), s.GetSeedName(), s.GetState(), util.Oneline(s.GetError(), 90)}
}

func map_get_wide_output(e interface{}) interface{} {
	s := e.(gube.Shoot)
	cnt := "unknown"
	c, err := s.GetNodeCount()
	if err == nil {
		cnt = fmt.Sprintf("%d", c)
	}
	return []string{s.GetName().GetName(), s.GetName().GetProjectName(),
		s.GetInfrastructure(), s.GetSeedName(), cnt, s.GetState(), util.Oneline(s.GetError(), 90)}
}

func NewGetHandler(opts *cmdint.Options) (util.Handler, error) {
	return NewOptsHandler(opts, get_outputs)
}
