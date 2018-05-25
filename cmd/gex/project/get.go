package project

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/output"
	"github.com/afritzler/garden-examiner/pkg"
	"github.com/afritzler/garden-examiner/pkg/data"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "get", get).CmdDescription("get projects(s)").
		CmdArgDescription("[<project>]")).
		ArgOption(constants.O_OUTPUT).Short('o').
		ArgOption(constants.O_SORT).Array()
}

func get(opts *cmdint.Options) error {
	return cmdline.ExecuteMode(opts, get_outputs, TypeHandler)
}

/////////////////////////////////////////////////////////////////////////////

var get_outputs = output.NewOutputs(get_regular)

func get_regular(opts *cmdint.Options) output.Output {
	return output.NewProcessingTableOutput(opts, data.Chain().Map(map_get_regular_output),
		"SEED", "INFRA", "REGION", "PROFILE", "SHOOT", "STATE", "ERROR")
}

func map_get_regular_output(e interface{}) interface{} {
	p := e.(gube.Project)
	return []string{p.GetName(), p.GetNamespace()}
}
