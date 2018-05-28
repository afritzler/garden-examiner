package shoot

import (
	"fmt"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/iaas"
	"github.com/afritzler/garden-examiner/cmd/gex/output"
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
	fmt.Printf("Shoot: %s\n", s.GetName().GetName())
	fmt.Printf("Project: %s\n", s.GetName().GetProjectName())
	fmt.Printf("Profile: %s (%s)\n", s.GetProfileName(), s.GetInfrastructure())
	seed, _ := s.GetNamespaceInSeed()
	fmt.Printf("Seed Namespace: %s\n", seed)
	cnt := "unknown"
	c, err := s.GetNodeCount()
	if err == nil {
		cnt = fmt.Sprintf("%d", c)
	}
	fmt.Printf("Size: %s\n", cnt)
	fmt.Printf("State: %s\n", s.GetState())
	//iaas, err := s.GetIaaSInfo()
	iaas.Describe(s)

	fmt.Printf("Error: %s\n", s.GetError())
	return nil
}
