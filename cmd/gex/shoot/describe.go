package shoot

import (
	"fmt"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/iaas"
	"github.com/afritzler/garden-examiner/cmd/gex/output"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/pkg"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "describe", describe).CmdDescription(
		"describe shoot(s)",
	).
		CmdArgDescription("[<shoot>]"))
}

func describe(opts *cmdint.Options) error {
	return cmdline.ExecuteOutput(opts, NewDescribeOutput(), TypeHandler)
}

/////////////////////////////////////////////////////////////////////////////

type describe_output struct {
	*output.ElementOutput
}

func NewDescribeOutput() *describe_output {
	o := &describe_output{}
	o.ElementOutput = output.NewElementOutput(nil)
	return o
}

func (this *describe_output) Out(ctx *context.Context) error {
	i := this.Elems.Iterator()
	for i.HasNext() {
		fmt.Printf("---\n")
		Describe(i.Next().(gube.Shoot))
	}
	return nil
}

func Describe(s gube.Shoot) error {
	attrs := util.NewAttributeSet()
	attrs.Attribute("Shoot", s.GetName().GetName())
	attrs.Attribute("Project", s.GetName().GetProjectName())
	attrs.Attributef("Profile", "%s (%s)", s.GetProfileName(), s.GetInfrastructure())
	seed, _ := s.GetNamespaceInSeed()
	attrs.Attribute("Seed Namespace", seed)
	cnt := "unknown"
	c, err := s.GetNodeCount()
	if err == nil {
		cnt = fmt.Sprintf("%d", c)
	}
	attrs.Attribute("Number of Nodes", cnt)
	attrs.Attribute("State", s.GetState())
	attrs.PrintAttributes()
	//iaas, err := s.GetIaaSInfo()
	iaas.Describe(s)

	fmt.Printf("Error: %s\n", s.GetError())
	return nil
}
