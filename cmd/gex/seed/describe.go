package seed

import (
	"fmt"
	"strconv"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/output"
	"github.com/afritzler/garden-examiner/cmd/gex/shoot"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/pkg"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "describe", describe).CmdDescription(
		"describe seed(s)",
	).
		CmdArgDescription("[<seed>]"))
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

type ShootCount func(name string) int

func (this *describe_output) Out(ctx *context.Context) error {
	shoots, err := ctx.GetShoots()
	f := func(name string) int {
		cnt := 0
		for _, s := range shoots {
			if s.GetSeedName() == name {
				cnt++
			}
		}
		return cnt
	}
	if err != nil {
		return err
	}
	i := this.Elems.Iterator()
	for i.HasNext() {
		fmt.Printf("---\n")
		Describe(i.Next().(gube.Seed), f)
	}
	return nil
}

func Describe(s gube.Seed, shoot_count ShootCount) error {

	attrs := util.NewAttributeSet()
	attrs.Attribute("Seed", s.GetName())
	cfg, err := s.GetClientConfig()
	if err == nil {
		attrs.Attribute("API Server", cfg.Host)
	} else {
		attrs.Attribute("API Server", "unknown")
	}
	attrs.Attributef("Profile", "%s (%s)", s.GetProfileName(), s.GetInfrastructure())
	attrs.Attribute("Region", s.GetRegion())
	cnt := "unknown"
	c, err := s.GetNodeCount()
	if err == nil {
		cnt = fmt.Sprintf("%d", c)
	}
	attrs.Attribute("Number of Nodes", cnt)
	attrs.Attribute("Number of Shoots", strconv.Itoa(shoot_count(s.GetName())))
	if s.GetShootName() != nil {
		sh, err := s.AsShoot()
		if err != nil {
			fmt.Printf("Seed is shooted, but cannot get shoot: %s\n", err)
		} else {
			shoot.Describe(sh, attrs)
		}
	}
	attrs.PrintAttributes()
	return nil
}
