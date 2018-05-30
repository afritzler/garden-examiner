package cmdint

type ConfigurableCommandOptionHelper struct {
	cmd    ConfigurableCommand
	helper *OptionConfigHelper
}

var _ ConfigurableCommand = &ConfigurableCommandOptionHelper{}

func (this *ConfigurableCommandOptionHelper) AsCommand() ConfigurableCommand {
	return this.cmd.AsCommand()
}

func (this *ConfigurableCommandOptionHelper) Execute(ctx *Options, args []string) error {
	return this.cmd.Execute(ctx, args)
}
func (this *ConfigurableCommandOptionHelper) GetCmdDescription() string {
	return this.cmd.GetCmdDescription()
}
func (this *ConfigurableCommandOptionHelper) GetCmdArgDescription() string {
	return this.cmd.GetCmdArgDescription()
}
func (this *ConfigurableCommandOptionHelper) Help(name string, args []string) {
	this.cmd.Help(name, args)
}

func (this *ConfigurableCommandOptionHelper) Command() ConfigurableCommand {
	return this.cmd
}

//
// Options
//
func (this *ConfigurableCommandOptionHelper) ArgOption(key string) *ConfigurableCommandOptionHelper {
	return this.cmd.ArgOption(key)
}
func (this *ConfigurableCommandOptionHelper) FlagOption(key string) *ConfigurableCommandOptionHelper {
	return this.cmd.FlagOption(key)
}

//
// Option attributes
//
func (this *ConfigurableCommandOptionHelper) Context(ctx string) *ConfigurableCommandOptionHelper {
	this.helper.Context(ctx)
	return this
}

func (this *ConfigurableCommandOptionHelper) Default(def interface{}) *ConfigurableCommandOptionHelper {
	this.helper.Default(def)
	return this
}

func (this *ConfigurableCommandOptionHelper) Short(short rune) *ConfigurableCommandOptionHelper {
	this.helper.Short(short)
	return this
}

func (this *ConfigurableCommandOptionHelper) Long(long string) *ConfigurableCommandOptionHelper {
	this.helper.Long(long)
	return this
}

func (this *ConfigurableCommandOptionHelper) ArgDescription(desc string) *ConfigurableCommandOptionHelper {
	this.helper.ArgDescription(desc)
	return this
}

func (this *ConfigurableCommandOptionHelper) Description(desc ...string) *ConfigurableCommandOptionHelper {
	this.helper.Description(desc...)
	return this
}

func (this *ConfigurableCommandOptionHelper) Args(n int) *ConfigurableCommandOptionHelper {
	this.helper.Args(n)
	return this
}

func (this *ConfigurableCommandOptionHelper) Array() *ConfigurableCommandOptionHelper {
	this.helper.Array()
	return this
}

func (this *ConfigurableCommandOptionHelper) Env(name string) *ConfigurableCommandOptionHelper {
	this.helper.Env(name)
	return this
}
