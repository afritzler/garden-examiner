package garden

import (
	"fmt"
	"strconv"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/output"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/pkg"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "describe", describe).CmdDescription(
		"describe garden(s)",
	).
		CmdArgDescription("[<garden>]"))
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
		out, err := NewOutput(i.Next().(gube.GardenConfig))
		if err != nil {
			return err
		}
		out.Describe()
	}
	return nil
}

type Output struct {
	config   gube.GardenConfig
	garden   gube.Garden
	infokube *util.InfoKube
	*util.AttributeSet
}

func NewOutput(cfg gube.GardenConfig) (*Output, error) {
	var err error
	o := &Output{}
	o.config = cfg
	g, err := cfg.GetGarden()
	if err != nil {
		return nil, err
	}
	o.garden = g
	dimensions := []string{"Infra", "Seed", "Profile", "Region"}
	o.infokube = util.NewInfoKube(dimensions)
	shoots, err := g.GetShoots()
	if err != nil {
		return nil, err
	}
	for _, s := range shoots {
		sn := s.GetSeedName()
		seed, err := s.GetSeed()
		if err == nil {
			sn = fmt.Sprintf("%s (%s)", s.GetSeedName(), seed.GetInfrastructure())
		}
		o.infokube.AddElement(nil, s.GetInfrastructure(), sn,
			s.GetProfileName(), s.GetRegion())
	}

	seeds, err := g.GetSeeds()
	if err != nil {
		return nil, err
	}
	for _, s := range seeds {
		o.infokube.AddKey("Seed", s.GetName())
	}
	o.AttributeSet = util.NewAttributeSet()
	return o, nil
}

func (this *Output) Describe() error {
	this.ResetAttributes()
	this.Attributef("Garden", "%s (%s)", this.config.GetName(), this.config.GetDescription())
	this.Attribute("Total Number of Shoots", strconv.Itoa(this.infokube.GetCount()))
	this.Attribute("Total Number of Seeds", strconv.Itoa(len(this.infokube.GetKeys("Seed"))))
	this.PrintAttributes()
	fmt.Printf("Infrastructure Overview:\n")
	this.infokube.Table("", []string{"Infra", "Region", "Seed"}, util.Coord{})
	return nil
}
