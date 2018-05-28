package verb

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"
)

func init() {
	NewVerb("iaas", cmdint.MainTab()).CmdArgDescription("<type> ...").
		CmdDescription("general iaas command",
			"The first argument is the element type followed by",
			"element name option and/or iaas arguments/options.",
			"If no element option is given, it must be defaulted by the",
			"selection command or the appropriate selection option.",
		).
		CatchUnknownCommand(catch_cluster).Raw()
}
