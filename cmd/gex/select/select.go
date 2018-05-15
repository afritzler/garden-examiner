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
	CmdDescription("set default shoot/eed/project\n" +
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
	elem interface{}
}

func NewSelectOutput() util.Output {
	return &select_output{}
}

func (this *select_output) Add(ctx *context.Context, e interface{}) error {
	if this.elem == nil {
		this.elem = e
		return nil
	}
	return fmt.Errorf("only one element can be selected, but multiple elements selected/found")
}

func (this *select_output) Out(ctx *context.Context) {
	shoot := ""
	seed := ""
	project := ""
	switch e := this.elem.(type) {
	case gube.Shoot:
		shoot = e.GetName().String()
		project = e.GetName().GetProjectName()
		seed = e.GetSeedName()
	case gube.Seed:
		seed = e.GetName()
	default:
		panic(fmt.Errorf("invalid elem type for select: %T\n", this.elem))
	}

	f := os.NewFile(3, "env-setting")
	if f == nil {
		fmt.Printf("Please use the gex alias to make the selection effective.\n")
	}
	envout(f, shoot, "SHOOT")
	envout(f, project, "PROJECT")
	envout(f, seed, "SEED")
}

func envout(f *os.File, value, key string) {
	line := ""
	if value == "" {
		line = fmt.Sprintf("unset GEX_%s", key)
	} else {
		line = fmt.Sprintf("export GEX_%s=\"%v\"", key, value)
	}
	if f != nil {
		fmt.Fprintf(f, "%s\n", line)
	}
	fmt.Println(line)

}
