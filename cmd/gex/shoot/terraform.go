package shoot

import (
	"fmt"

	"github.com/afritzler/garden-examiner/pkg"
	"github.com/afritzler/garden-examiner/pkg/data"

	"github.com/afritzler/garden-examiner/cmd/gex/output"
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "terraform", terraform).
		CmdDescription("get terraform data for shoot",
			"Terraform data is available for the following jobs (--job):",
			"  - infra:       infrastructure setup (default)",
			"  - external-dns: setup of external DNS entries",
			"  - internal-dns: setup of DNS entries for internal IP addresses",
			"  - ingress:      setup of ingress resources on IaaS layer",
			"",
			"For every job the following information is available (--data):",
			"  - config: configuration values",
			"  - script: the terraform script",
			"  - state:  the terraform state (default)").
		CmdArgDescription("[<shoot>]")).
		ArgOption("job").Short('j').Default("infra").
		ArgOption("data").Short('d').Default("state")
}

var jobs []string = []string{"infra", "external-dns", "internal-dns", "ingress"}
var datas []string = []string{"config", "state", "script"}

func terraform(opts *cmdint.Options) error {
	job := opts.GetOptionValue("job")
	data := opts.GetOptionValue("data")
	j, _ := cmdint.SelectBest(*job, jobs...)
	d, _ := cmdint.SelectBest(*data, datas...)

	if j == "" {
		return fmt.Errorf("invalid job '%s', expected one of infra, external-dns, internal-dns, ingress", j)
	}
	if d == "" {
		return fmt.Errorf("invalid data type '%s', expected one of: state,config,script", d)
	}

	return cmdline.ExecuteOutput(opts, NewTerraformOutput(j, d), TypeHandler)
}

func NewTerraformOutput(cm string, field string) output.Output {
	return output.NewStringOutput(terraform_output_mapper(cm, field))
}

func terraform_output_mapper(job, data string) data.MappingFunction {
	return func(e interface{}) interface{} {
		s := e.(gube.Shoot)
		result, err := s.GetTerraformJobData(job, data)
		if err != nil {
			return err
		}
		return result
	}
}
