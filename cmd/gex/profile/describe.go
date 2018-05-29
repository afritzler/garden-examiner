package profile

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
		"describe profile(s)",
	).
		CmdArgDescription("[<profile>]"))
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
	out, err := NewOutput(ctx.Garden)

	if err != nil {
		return err
	}
	i := this.Elems.Iterator()
	for i.HasNext() {
		fmt.Printf("---\n")
		out.Describe(i.Next().(gube.Profile))
	}
	return nil
}

type Output struct {
	shoots map[gube.ShootName]gube.Shoot
	seeds  map[string]gube.Seed
	*util.AttributeSet
}

func NewOutput(g gube.Garden) (*Output, error) {
	var err error
	o := &Output{}
	o.shoots, err = g.GetShoots()
	if err != nil {
		return nil, err
	}
	o.seeds, err = g.GetSeeds()
	if err != nil {
		return nil, err
	}
	o.AttributeSet = util.NewAttributeSet()
	return o, nil
}

func (this *Output) CountProfileShoots(name string) int {
	cnt := 0
	for _, s := range this.shoots {
		if s.GetProfileName() == name {
			cnt++
		}
	}
	return cnt
}

func (this *Output) CountSeedProfileShoots(seed string, name string) int {
	cnt := 0
	for _, s := range this.shoots {
		if s.GetSeedName() == seed && s.GetProfileName() == name {
			cnt++
		}
	}
	return cnt
}

func (this *Output) Describe(p gube.Profile) error {
	this.ResetAttributes()
	this.Attribute("Profile", p.GetName())
	this.Attribute("Infrastructure", p.GetInfrastructure())
	this.Attribute("Total Number of Shoots", strconv.Itoa(this.CountProfileShoots(p.GetName())))
	table := [][]string{[]string{"Seed", "Infra", "Region", "-Shoots"}}
	used := 0
	for _, s := range this.seeds {
		cnt := this.CountSeedProfileShoots(s.GetName(), p.GetName())
		if cnt > 0 || s.GetProfileName() == p.GetName() {
			flag := ""
			if s.GetProfileName() != p.GetName() {
				flag = "*"
			} else {
				used++
			}
			table = append(table, []string{flag + s.GetName(), s.GetInfrastructure(), s.GetRegion(), strconv.Itoa(cnt)})
		}
	}
	this.Attribute("Number of assigned Seeds", strconv.Itoa(used))
	this.PrintAttributes()
	if len(table) > 1 {
		util.FormatTable("  ", table)
	}
	return nil
}
