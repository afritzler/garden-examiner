package verb

import (
	"github.com/afritzler/garden-examiner/cmd/gex/const"

	"github.com/mandelsoft/cmdint/pkg/cmdint"
)

func init() {
	get := NewVerb("get", cmdint.MainTab()).CmdArgDescription("<type> ...").
		CmdDescription("general get command",
			"The first argument is the element type followed by",
			"element names and/or options",
		).
		CatchUnknownCommand(catch_cluster).
		ArgOption(constants.O_OUTPUT).Short('o')

	NewVerb("kubeconfig", get).CmdArgDescription("<cluster type> ...").
		CmdDescription("general kubeconfig get command",
			"The first argument is the cluster type (shoot or seed) followed by",
			"element names and/or options",
		)

	NewVerb("iaas", cmdint.MainTab()).CmdArgDescription("<type> ...").Raw().
		CmdDescription("general iaas command",
			"The first argument is the element type followed by",
			"element name option and/or iaas arguments/options.",
			"If no element option is given, it must be defaulted by the",
			"selection command or the appropriate selection option.",
		).
		CatchUnknownCommand(catch_cluster)

	NewVerb("describe", cmdint.MainTab()).CmdArgDescription("<type> ...").
		CmdDescription("general describe command",
			"The first argument is the element type followed by",
			"element name option and/or iaas arguments/options.",
			"If no element option is given, it must be defaulted by the",
			"selection command or the appropriate selection option.",
		).
		CatchUnknownCommand(catch_cluster)
}

type Verb struct {
	cmdtab cmdint.ConfigurableCmdTab
}

var verbs map[string]*Verb = map[string]*Verb{}

func NewVerb(name string, tab cmdint.ConfigurableCmdTab) cmdint.ConfigurableCmdTab {
	verb := &Verb{cmdint.NewCmdTab(name)}
	verbs[name] = verb
	tab.Command(name, verb.cmdtab)
	return verb.cmdtab
}

func GetVerb(name string) cmdint.ConfigurableCmdTab {
	v := verbs[name]
	if v == nil {
		return nil
	}
	return v.cmdtab
}

func Add(tab cmdint.ConfigurableCmdTab, name string, cmd cmdint.CommandFunction) *cmdint.CmdTabCommandHelper {
	c := tab.SimpleCommand(name, cmd)
	v := GetVerb(name)
	if v != nil {
		sc := c.AsSubCommand()

		v.Command(tab.GetName(), sc)
	}
	return c
}
