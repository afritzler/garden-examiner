package cmd_select

import (
	"fmt"
	"os"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"

	"github.com/afritzler/garden-examiner/pkg"
)

var cmdtab cmdint.ConfigurableCmdTab = cmdint.NewCmdTab("select").
	DefaultFunction(cmd_select).
	CmdDescription("set default shoot/seed/project\n" +
		"works only with the gex alias feeding the appropriate\n" +
		"environment variables.").
	CmdArgDescription("<command> <element name>")

func init() {
	cmdint.MainTab().
		Command("select", cmdtab)
}

func GetCmdTab() cmdint.ConfigurableCmdTab {
	return cmdtab
}

func cmd_select(opts *cmdint.Options) error {
	found := 0
	if v := opts.GetOptionValue(constants.O_SEL_SHOOT); v != nil {
		fmt.Printf("SHOOT = %s\n", *v)
		found++
	}
	if v := opts.GetOptionValue(constants.O_SEL_PROJECT); v != nil {
		fmt.Printf("PROJECT = %s\n", *v)
		found++
	}
	if v := opts.GetOptionValue(constants.O_SEL_SEED); v != nil {
		fmt.Printf("SEED = %s\n", *v)
		found++
	}
	if found == 0 {
		return fmt.Errorf("no selection found")
	}
	return nil
}

type select_output struct {
	*util.SingleElementOutput
}

var _ util.Output = &select_output{}

func NewSelectOutput() *select_output {
	return &select_output{util.NewSingleElementOutput()}
}

func (this *select_output) Out(ctx *context.Context) {
	shoot := ""
	seed := ""
	project := ""
	switch e := this.Elem.(type) {
	case gube.Shoot:
		shoot = e.GetName().String()
		project = e.GetName().GetProjectName()
		seed = e.GetSeedName()
	case gube.Seed:
		seed = e.GetName()
	case gube.Project:
		project = e.GetName()
	default:
		panic(fmt.Errorf("invalid elem type for select: %T\n", this.Elem))
	}

	this.Write(&shoot, &project, &seed)
}

func (this *select_output) Write(shoot, project, seed *string) {
	f := os.NewFile(3, "env-setting")
	if f == nil {
		fmt.Printf("Please use the gex alias to make the selection effective.\n")
	}
	envout(f, shoot, "SHOOT")
	envout(f, project, "PROJECT")
	envout(f, seed, "SEED")
}

func envout(f *os.File, value *string, key string) {
	line := ""
	if value == nil || *value == "" {
		line = fmt.Sprintf("unset GEX_%s", key)
	} else {
		line = fmt.Sprintf("export GEX_%s=\"%v\"", key, *value)
	}
	if f != nil {
		fmt.Fprintf(f, "%s\n", line)
	}
	fmt.Println(line)

}
