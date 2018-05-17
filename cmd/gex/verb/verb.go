package verb

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"
)

func init() {
	get := NewVerb("get", cmdint.MainTab()).CmdArgDescription("<type> ...").
		CmdDescription("general get command",
			"The first argument is the element type followed by",
			"element names and/or options",
		)

	NewVerb("kubeconfig", get).CmdArgDescription("<cluster type> ...").
		CmdDescription("general kubeconfig get command",
			"The first argument is the cluster type (shoot or seed) followed by",
			"element names and/or options",
		)
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
		v.Command(tab.GetName(), c.AsSubCommand())
	}
	return c
}
