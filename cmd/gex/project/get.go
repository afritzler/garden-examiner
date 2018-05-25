package project

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/cmd/gex/verb"
	"github.com/afritzler/garden-examiner/pkg"
	"github.com/afritzler/garden-examiner/pkg/data"
)

func init() {
	filters.AddOptions(verb.Add(GetCmdTab(), "get", get).CmdDescription("get projects(s)").
		CmdArgDescription("[<project>]")).
		ArgOption(constants.O_OUTPUT).Short('o').
		ArgOption(constants.O_SORT).Array()
}

func get(opts *cmdint.Options) error {
	return util.ExecuteMode(opts, get_outputs, TypeHandler)
}

/////////////////////////////////////////////////////////////////////////////

var get_outputs = util.NewOutputs(get_regular)

func get_regular(opts *cmdint.Options) util.Output {
	return util.NewProcessingTableOutput(opts, data.Chain().Map(map_get_regular_output),
		"SEED", "INFRA", "REGION", "PROFILE", "SHOOT", "STATE", "ERROR")
}

func map_get_regular_output(e interface{}) interface{} {
	p := e.(gube.Project)
	return []string{p.GetName(), p.GetNamespace()}
}
