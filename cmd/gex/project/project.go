package shoot

import (
	_ "github.com/afritzler/garden-examiner/pkg"
	"github.com/mandelsoft/cmdint/pkg/cmdint"
)

var cmdtab cmdint.ConfigurableCmdTab = cmdint.NewCmdTab("project", nil).
	CmdDescription("garden projects\n" +
		"list one or more projects").
	CmdArgDescription("<command>")

func init() {
	cmdint.MainTab().
		Command("project", cmdtab)
}

func GetCmdTab() cmdint.ConfigurableCmdTab {
	return cmdtab
}
