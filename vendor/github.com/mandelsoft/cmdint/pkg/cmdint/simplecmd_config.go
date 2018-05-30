package cmdint

///////////////////////////////////////////////////////////////////////////
// Configuration
///////////////////////////////////////////////////////////////////////////

func (this *SimpleCommand) AsCommand() ConfigurableCommand {
	return &SimpleConfigurableCommandHelper{this}
}

///////////////////////////////////////////////////////////////////////////
// Command  Config
//

func (this *SimpleCommand) CmdDescription(desc ...string) *SimpleCommand {
	this.desc = compact(desc)
	return this
}
func (this *SimpleCommand) CmdArgDescription(desc string) *SimpleCommand {
	this.argdesc = desc
	return this
}

func (this *SimpleCommand) Mixed() *SimpleCommand {
	this.optionspec.Mixed()
	return this
}
func (this *SimpleCommand) Raw() *SimpleCommand {
	this.optionspec.Raw()
	return this
}

func (this *SimpleCommand) ArgOption(key string) *SimpleCommandOptionHelper {
	return &SimpleCommandOptionHelper{this, this.optionspec.ArgOption(key)}
}
func (this *SimpleCommand) FlagOption(key string) *SimpleCommandOptionHelper {
	return &SimpleCommandOptionHelper{this, this.optionspec.FlagOption(key)}
}

///////////////////////////////////////////////////////////////////////////
// Command Option Config
//

type SimpleCommandOptionHelper struct {
	cmd    *SimpleCommand
	helper *OptionConfigHelper
}

var _ Command = &SimpleCommandOptionHelper{}

func (this *SimpleCommandOptionHelper) AsCommand() ConfigurableCommand {
	return this.cmd.AsCommand()
}

func (this *SimpleCommandOptionHelper) Execute(ctx *Options, args []string) error {
	return this.cmd.Execute(ctx, args)
}
func (this *SimpleCommandOptionHelper) GetCmdDescription() string {
	return this.cmd.GetCmdDescription()
}
func (this *SimpleCommandOptionHelper) GetCmdArgDescription() string {
	return this.cmd.GetCmdArgDescription()
}
func (this *SimpleCommandOptionHelper) Help(name string, args []string) {
	this.cmd.Help(name, args)
}
func (this *SimpleCommandOptionHelper) GetOptions() OptionSpec {
	return this.cmd.GetOptions()
}

//
// Options
//
func (this *SimpleCommandOptionHelper) ArgOption(key string) *SimpleCommandOptionHelper {
	return &SimpleCommandOptionHelper{this.cmd, this.helper.spec.ArgOption(key)}
}
func (this *SimpleCommandOptionHelper) FlagOption(key string) *SimpleCommandOptionHelper {
	return &SimpleCommandOptionHelper{this.cmd, this.helper.spec.FlagOption(key)}
}

//
// Option attributes
//
func (this *SimpleCommandOptionHelper) Context(ctx string) *SimpleCommandOptionHelper {
	this.helper.Context(ctx)
	return this
}

func (this *SimpleCommandOptionHelper) Default(def interface{}) *SimpleCommandOptionHelper {
	this.helper.Default(def)
	return this
}

func (this *SimpleCommandOptionHelper) Short(short rune) *SimpleCommandOptionHelper {
	this.helper.Short(short)
	return this
}

func (this *SimpleCommandOptionHelper) Long(long string) *SimpleCommandOptionHelper {
	this.helper.Long(long)
	return this
}

func (this *SimpleCommandOptionHelper) ArgDescription(desc string) *SimpleCommandOptionHelper {
	this.helper.ArgDescription(desc)
	return this
}

func (this *SimpleCommandOptionHelper) Description(desc ...string) *SimpleCommandOptionHelper {
	this.helper.Description(desc...)
	return this
}

func (this *SimpleCommandOptionHelper) Args(n int) *SimpleCommandOptionHelper {
	this.helper.Args(n)
	return this
}

func (this *SimpleCommandOptionHelper) Array() *SimpleCommandOptionHelper {
	this.helper.Array()
	return this
}

func (this *SimpleCommandOptionHelper) Env(name string) *SimpleCommandOptionHelper {
	this.helper.Env(name)
	return this
}

///////////////////////////////////////////////////////////////////////////
// General Configurable Command

type SimpleConfigurableCommandHelper struct {
	cmd *SimpleCommand
}

var _ ConfigurableCommand = &SimpleConfigurableCommandHelper{}

func (this *SimpleConfigurableCommandHelper) AsCommand() ConfigurableCommand {
	return this
}

func (this *SimpleConfigurableCommandHelper) Execute(ctx *Options, args []string) error {
	return this.cmd.Execute(ctx, args)
}
func (this *SimpleConfigurableCommandHelper) GetCmdDescription() string {
	return this.cmd.GetCmdDescription()
}
func (this *SimpleConfigurableCommandHelper) GetCmdArgDescription() string {
	return this.cmd.GetCmdArgDescription()
}
func (this *SimpleConfigurableCommandHelper) Help(name string, args []string) {
	this.cmd.Help(name, args)
}
func (this *SimpleConfigurableCommandHelper) GetOptions() OptionSpec {
	return this.cmd.GetOptions()
}

//
// Option Config
//

func (this *SimpleConfigurableCommandHelper) ArgOption(name string) *ConfigurableCommandOptionHelper {
	return &ConfigurableCommandOptionHelper{this, this.cmd.ArgOption(name).helper}
}

func (this *SimpleConfigurableCommandHelper) FlagOption(name string) *ConfigurableCommandOptionHelper {
	return &ConfigurableCommandOptionHelper{this, this.cmd.FlagOption(name).helper}
}
