package verb

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"
)

func init() {
	NewVerb("kubectl", cmdint.MainTab()).CmdArgDescription("<type> ...").
		CmdDescription("general kubectl command",
			"The first argument is the element type followed by",
			"element name option and/or kubectl arguments/options.",
			"If no element option is given, it must be defaulted by the",
			"selection command or the appropriate selection option.",
			"If nothing is slected kubectl is run for the garden cluster.",
		).
		CatchUnknownCommand(catch_cluster).Raw().
		CmdArgDescription("{<kubectl opts/args>}").
		CmdDescription("run kubectl for garden cluster")
}
