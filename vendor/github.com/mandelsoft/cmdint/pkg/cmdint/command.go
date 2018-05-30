package cmdint

type CommandFunction func(opts *Options) error

type Command interface {
	GetCmdArgDescription() string
	GetCmdDescription() string
	Execute(ctx *Options, args []string) error
	Help(name string, args []string)
}

type ConfigurableCommandFlavor interface {
	AsCommand() ConfigurableCommand
}

type ConfigurableCommand interface {
	Command
	ConfigurableCommandFlavor

	ArgOption(key string) *ConfigurableCommandOptionHelper
	FlagOption(key string) *ConfigurableCommandOptionHelper
}
