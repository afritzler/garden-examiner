package shoot

import (
	"fmt"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/output"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/pkg"
	"github.com/afritzler/garden-examiner/pkg/data"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "get", get).CmdDescription(
		"get shoot(s)",
		"supported output modes are:",
		"- yaml|json|JSON  print manifest",
		"- wide            additional info",
		"- kubeconfig      print kube config",
		"- error           show complete error message",
	).
		CmdArgDescription("[<shoot>]").Mixed()).
		ArgOption(constants.O_OUTPUT).Short('o').
		ArgOption(constants.O_SORT).Array()
}

func get(opts *cmdint.Options) error {
	return cmdline.ExecuteMode(opts, get_outputs, TypeHandler)
}

/////////////////////////////////////////////////////////////////////////////

var get_outputs = output.NewOutputs(get_regular, output.Outputs{
	"wide":       get_wide,
	"kubeconfig": output.KubeconfigOutputFactory,
	"error":      get_error,
}).AddManifestOutputs()

func get_regular(opts *cmdint.Options) output.Output {
	return output.NewProcessingTableOutput(opts, data.Chain().Map(map_get_regular_output),
		"SHOOT", "PROJECT", "INFRA", "PROFLE", "SEED", "STATE", "ERROR")
}
func get_wide(opts *cmdint.Options) output.Output {
	return output.NewProcessingTableOutput(opts, data.Chain().Parallel(20).Map(map_get_wide_output),
		"SHOOT", "PROJECT", "INFRA", "PROFILE", "SEED", "-NODES", "IAAS", "STATE", "ERROR")
}
func get_error(opts *cmdint.Options) output.Output {
	return output.NewProcessingTableOutput(opts, data.Chain().Parallel(20).Map(map_get_error_output),
		"SHOOT", "ERROR")
}

/////////////////////////////////////////////////////////////////////////////

func map_get_regular_output(e interface{}) interface{} {
	s := e.(gube.Shoot)
	return []string{s.GetName().GetName(), s.GetName().GetProjectName(),
		s.GetInfrastructure(), s.GetProfileName(), s.GetSeedName(), s.GetState(), util.Oneline(s.GetError(), 90)}
}

func map_get_wide_output(e interface{}) interface{} {
	s := e.(gube.Shoot)
	cnt := "unknown"
	c, err := s.GetNodeCount()
	if err == nil {
		cnt = fmt.Sprintf("%d", c)
	}
	iaas, err := s.GetIaaSInfo()
	info := "unknown"
	if err == nil {
		info = iaas.GetKeyInfo()
	}
	return []string{s.GetName().GetName(), s.GetName().GetProjectName(),
		s.GetInfrastructure(), s.GetProfileName(), s.GetSeedName(), cnt, info, s.GetState(), util.Oneline(s.GetError(), 90)}
}

func map_get_error_output(e interface{}) interface{} {
	s := e.(gube.Shoot)
	return []string{s.GetName().GetName(), s.GetError()}
}
