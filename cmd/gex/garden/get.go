package garden

import (
	"fmt"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/output"
	"github.com/afritzler/garden-examiner/pkg"
	"github.com/afritzler/garden-examiner/pkg/data"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "get", get).CmdDescription("get configured garden(s)").
		CmdArgDescription("[<garden>]")).
		ArgOption(constants.O_OUTPUT).Short('o')
}

func get(opts *cmdint.Options) error {
	return cmdline.ExecuteMode(opts, get_outputs, TypeHandler)
}

/////////////////////////////////////////////////////////////////////////////

var get_outputs = output.NewOutputs(get_regular, output.Outputs{
	"wide": get_wide,
}).AddManifestOutputs()

func get_regular(opts *cmdint.Options) output.Output {
	return output.NewProcessingTableOutput(opts, data.Chain().Map(map_get_regular_output),
		"GARDEN", "HOST", "DESCRIPTION")
}

func get_wide(opts *cmdint.Options) output.Output {
	return output.NewProcessingTableOutput(opts, data.Chain().Parallel(20).Map(map_get_wide_output),
		"GARDEN", "HOST", "-PROJECTS", "-SHOOTS", "-SEEDS", "DESCRIPTION")
}

func map_get_regular_output(e interface{}) interface{} {
	p := e.(gube.GardenConfig)
	g, err := p.GetGarden()
	host := "unknown"
	if err == nil {
		cfg, err := g.GetClientConfig()
		if err == nil {
			host = cfg.Host
		}
	}
	return []string{p.GetName(), host, p.GetDescription()}
}

func map_get_wide_output(e interface{}) interface{} {
	c := e.(gube.GardenConfig)
	projects := "unknown"
	shoots := "unknown"
	seeds := "unknown"
	host := "unknown"
	g, err := c.GetGarden()
	if err == nil {
		pr, err := g.GetProjects()
		if err == nil {
			projects = fmt.Sprintf("%d", len(pr))
		}
		sh, err := g.GetShoots()
		if err == nil {
			shoots = fmt.Sprintf("%d", len(sh))
		}
		se, err := g.GetSeeds()
		if err == nil {
			seeds = fmt.Sprintf("%d", len(se))
		}

		cfg, err := g.GetClientConfig()
		if err == nil {
			host = cfg.Host
		}
	}
	return []string{c.GetName(), host, projects, shoots, seeds, c.GetDescription()}
}
