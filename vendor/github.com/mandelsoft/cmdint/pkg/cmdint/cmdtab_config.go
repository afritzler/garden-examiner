package cmdint

import (
	"fmt"
)

type ConfigurableCmdTab interface {
	CmdTab
	GetOptions() OptionSpec
	ConfigurableCommandFlavor

	AsCmdTab() ConfigurableCmdTab

	SetupFunction(f CommandFunction) ConfigurableCmdTab
	DefaultFunction(f CommandFunction) ConfigurableCmdTab
	CatchUnknownCommand(f CatchFunction) ConfigurableCmdTab
	Raw() ConfigurableCmdTab
	CmdArgDescription(desc string) ConfigurableCmdTab
	CmdDescription(desc ...string) ConfigurableCmdTab
	Option(*Option) ConfigurableCmdTab
	Command(name string, c Command) ConfigurableCmdTab
	SimpleCommand(name string, cmd CommandFunction) *CmdTabCommandHelper
	ArgOption(key string) *CmdTabOptionHelper
	FlagOption(key string) *CmdTabOptionHelper
}

type ConfigurableCmdTabCommand interface {
	AsCmdTab() ConfigurableCmdTab
	AsSubCommand() *SimpleCommand
	Execute(ctx *Options, args []string) error
	GetCmdDescription() string
	GetCmdArgDescription() string
	GetName() string

	SetupFunction(f CommandFunction) ConfigurableCmdTab
	DefaultFunction(f CommandFunction) ConfigurableCmdTab
	CatchUnknownCommand(f CatchFunction) ConfigurableCmdTab
	CmdDescription(desc ...string) *CmdTabCommandHelper
	CmdArgDescription(desc string) *CmdTabCommandHelper
	Command(name string, c Command) ConfigurableCmdTab
	SimpleCommand(name string, cmd CommandFunction) *CmdTabCommandHelper
	ArgOption(key string) *CmdTabCommandOptionHelper
	FlagOption(key string) *CmdTabCommandOptionHelper
}

//
// Configuration
//

func (this *_CmdTab) AsCmdTab() ConfigurableCmdTab {
	return this
}

func (this *_CmdTab) AsCommand() ConfigurableCommand {
	return &CmdTabConfigurableCommandHelper{this}
}

func (this *_CmdTab) CmdArgDescription(desc string) ConfigurableCmdTab {
	this.argdesc = desc
	return this
}

func (this *_CmdTab) CmdDescription(desc ...string) ConfigurableCmdTab {
	this.desc = compact(desc)
	return this
}

func (this *_CmdTab) SetupFunction(setup CommandFunction) ConfigurableCmdTab {
	this.setup = setup
	return this
}
func (this *_CmdTab) DefaultFunction(f CommandFunction) ConfigurableCmdTab {
	this.deffunc = f
	return this
}
func (this *_CmdTab) CatchUnknownCommand(f CatchFunction) ConfigurableCmdTab {
	this.catchfunc = f
	return this
}
func (this *_CmdTab) Raw() ConfigurableCmdTab {
	this.optionspec.Raw()
	return this
}
func (this *_CmdTab) Option(o *Option) ConfigurableCmdTab {
	if this.optionspec.Get(o.Key) != nil {
		panic(fmt.Errorf("option '%s' already defined", o.Key))
	}
	return this
}
func (this *_CmdTab) Command(name string, c Command) ConfigurableCmdTab {
	if _, ok := this.commands[name]; ok {
		panic(fmt.Errorf("command '%s' already declared", name))
	}
	this.commands[name] = &CommandEntry{command: c}
	return this
}

func (this *_CmdTab) SimpleCommand(name string, cmd CommandFunction) *CmdTabCommandHelper {
	if _, ok := this.commands[name]; ok {
		panic(fmt.Errorf("command '%s' already declared", name))
	}
	c := NewCommand(cmd)
	this.commands[name] = &CommandEntry{command: c}
	return &CmdTabCommandHelper{this, c}
}

func (this *_CmdTab) ArgOption(key string) *CmdTabOptionHelper {
	return &CmdTabOptionHelper{this, this.optionspec.ArgOption(key)}
}
func (this *_CmdTab) FlagOption(key string) *CmdTabOptionHelper {
	return &CmdTabOptionHelper{this, this.optionspec.FlagOption(key)}
}

///////////////////////////////////////////////////////////////////////////
// Command Table Command Configuration

type CmdTabCommandHelper struct {
	cmds *_CmdTab
	cmd  *SimpleCommand
}

var _ CmdTab = &CmdTabCommandHelper{}

func (this *CmdTabCommandHelper) AsCmdTab() ConfigurableCmdTab {
	return this.cmds
}

func (this *CmdTabCommandHelper) AsSubCommand() *SimpleCommand {
	return this.cmd
}

func (this *CmdTabCommandHelper) Execute(ctx *Options, args []string) error {
	return this.cmds.Execute(ctx, args)
}
func (this *CmdTabCommandHelper) GetCmdDescription() string {
	return this.cmds.GetCmdDescription()
}
func (this *CmdTabCommandHelper) GetCmdArgDescription() string {
	return this.cmds.GetCmdArgDescription()
}
func (this *CmdTabCommandHelper) GetName() string {
	return this.cmds.GetName()
}
func (this *CmdTabCommandHelper) Help(cmd string, args []string) {
	this.cmds.Help(cmd, args)
}
func (this *CmdTabCommandHelper) GetOptions() OptionSpec {
	return this.cmds.GetOptions()
}
func (this *CmdTabCommandHelper) GetDefaultFunction() CommandFunction {
	return this.cmds.GetDefaultFunction()
}
func (this *CmdTabCommandHelper) GetCommand(name string) Command {
	return this.cmds.GetCommand(name)
}

func (this *CmdTabCommandHelper) SetupFunction(setup CommandFunction) ConfigurableCmdTab {
	return this.cmds.SetupFunction(setup)
}
func (this *CmdTabCommandHelper) DefaultFunction(f CommandFunction) ConfigurableCmdTab {
	return this.cmds.DefaultFunction(f)
}
func (this *CmdTabCommandHelper) CatchUnknownCommand(f CatchFunction) ConfigurableCmdTab {
	return this.cmds.CatchUnknownCommand(f)
}

func (this *CmdTabCommandHelper) Mixed() *CmdTabCommandHelper {
	this.cmd.Mixed()
	return this
}
func (this *CmdTabCommandHelper) Raw() *CmdTabCommandHelper {
	this.cmd.Raw()
	return this
}

func (this *CmdTabCommandHelper) CmdDescription(desc ...string) *CmdTabCommandHelper {
	this.cmd.desc = compact(desc)
	return this
}
func (this *CmdTabCommandHelper) CmdArgDescription(desc string) *CmdTabCommandHelper {
	this.cmd.argdesc = desc
	return this
}

func (this *CmdTabCommandHelper) Command(name string, c Command) ConfigurableCmdTab {
	return this.cmds.Command(name, c)
}
func (this *CmdTabCommandHelper) SimpleCommand(name string, cmd CommandFunction) *CmdTabCommandHelper {
	return this.cmds.SimpleCommand(name, cmd)
}

func (this *CmdTabCommandHelper) ArgOption(key string) *CmdTabCommandOptionHelper {
	return &CmdTabCommandOptionHelper{this, this.cmd.ArgOption(key)}
}
func (this *CmdTabCommandHelper) FlagOption(key string) *CmdTabCommandOptionHelper {
	return &CmdTabCommandOptionHelper{this, this.cmd.FlagOption(key)}
}

///////////////////////////////////////////////////////////////////////////
// Command Table as Command Configuration

type CmdTabOptionHelper struct {
	cmds   *_CmdTab
	helper *OptionConfigHelper
}

var _ CmdTab = &CmdTabOptionHelper{}

func (this *CmdTabOptionHelper) AsCmdTab() ConfigurableCmdTab {
	return this.cmds
}
func (this *CmdTabOptionHelper) AsCommand() ConfigurableCommand {
	return this.cmds.AsCommand()
}

func (this *CmdTabOptionHelper) Execute(ctx *Options, args []string) error {
	return this.cmds.Execute(ctx, args)
}
func (this *CmdTabOptionHelper) GetCmdDescription() string {
	return this.cmds.GetCmdDescription()
}
func (this *CmdTabOptionHelper) GetCmdArgDescription() string {
	return this.cmds.GetCmdArgDescription()
}
func (this *CmdTabOptionHelper) GetName() string {
	return this.cmds.GetName()
}
func (this *CmdTabOptionHelper) Help(name string, args []string) {
	this.cmds.Help(name, args)
}
func (this *CmdTabOptionHelper) GetOptions() OptionSpec {
	return this.cmds.GetOptions()
}
func (this *CmdTabOptionHelper) GetDefaultFunction() CommandFunction {
	return this.cmds.GetDefaultFunction()
}
func (this *CmdTabOptionHelper) GetCommand(name string) Command {
	return this.cmds.GetCommand(name)
}

func (this *CmdTabOptionHelper) SetupFunction(setup CommandFunction) ConfigurableCmdTab {
	this.cmds.SetupFunction(setup)
	return this
}
func (this *CmdTabOptionHelper) DefaultFunction(f CommandFunction) ConfigurableCmdTab {
	this.cmds.DefaultFunction(f)
	return this
}
func (this *CmdTabOptionHelper) CatchUnknownCommand(f CatchFunction) ConfigurableCmdTab {
	this.cmds.CatchUnknownCommand(f)
	return this
}
func (this *CmdTabOptionHelper) Raw() ConfigurableCmdTab {
	this.cmds.Raw()
	return this
}
func (this *CmdTabOptionHelper) CmdArgDescription(desc string) ConfigurableCmdTab {
	this.cmds.CmdArgDescription(desc)
	return this.cmds
}
func (this *CmdTabOptionHelper) CmdDescription(desc ...string) ConfigurableCmdTab {
	this.cmds.CmdDescription(desc...)
	return this.cmds
}
func (this *CmdTabOptionHelper) Option(o *Option) ConfigurableCmdTab {
	return this.cmds.Option(o)
}
func (this *CmdTabOptionHelper) Command(name string, c Command) ConfigurableCmdTab {
	return this.cmds.Command(name, c)
}
func (this *CmdTabOptionHelper) SimpleCommand(name string, cmd CommandFunction) *CmdTabCommandHelper {
	return this.cmds.SimpleCommand(name, cmd)
}

func (this *CmdTabOptionHelper) ArgOption(key string) *CmdTabOptionHelper {
	return &CmdTabOptionHelper{this.cmds, this.helper.ArgOption(key)}
}
func (this *CmdTabOptionHelper) FlagOption(key string) *CmdTabOptionHelper {
	return &CmdTabOptionHelper{this.cmds, this.helper.FlagOption(key)}
}

func (this *CmdTabOptionHelper) Context(ctx string) *CmdTabOptionHelper {
	this.helper.Context(ctx)
	return this
}
func (this *CmdTabOptionHelper) Default(def interface{}) *CmdTabOptionHelper {
	this.helper.Default(def)
	return this
}
func (this *CmdTabOptionHelper) Short(short rune) *CmdTabOptionHelper {
	this.helper.Short(short)
	return this
}
func (this *CmdTabOptionHelper) Long(long string) *CmdTabOptionHelper {
	this.helper.Long(long)
	return this
}
func (this *CmdTabOptionHelper) ArgDescription(desc string) *CmdTabOptionHelper {
	this.helper.ArgDescription(desc)
	return this
}
func (this *CmdTabOptionHelper) Description(desc ...string) *CmdTabOptionHelper {
	this.helper.Description(desc...)
	return this
}
func (this *CmdTabOptionHelper) Args(n int) *CmdTabOptionHelper {
	this.helper.Args(n)
	return this
}
func (this *CmdTabOptionHelper) Array() *CmdTabOptionHelper {
	this.helper.Array()
	return this
}
func (this *CmdTabOptionHelper) Env(name string) *CmdTabOptionHelper {
	this.helper.Env(name)
	return this
}

///////////////////////////////////////////////////////////////////////////
// Command Table Command Option Configuration

type CmdTabCommandOptionHelper struct {
	*CmdTabCommandHelper
	helper *SimpleCommandOptionHelper
}

var _ CmdTab = &CmdTabCommandOptionHelper{}

func (this *CmdTabCommandOptionHelper) Execute(ctx *Options, args []string) error {
	return this.cmds.Execute(ctx, args)
}
func (this *CmdTabCommandOptionHelper) GetCmdDescription() string {
	return this.cmds.GetCmdDescription()
}
func (this *CmdTabCommandOptionHelper) GetCmdArgDescription() string {
	return this.cmds.GetCmdArgDescription()
}
func (this *CmdTabCommandOptionHelper) GetName() string {
	return this.cmds.GetName()
}

func (this *CmdTabCommandOptionHelper) ArgOption(key string) *CmdTabCommandOptionHelper {
	return &CmdTabCommandOptionHelper{this.CmdTabCommandHelper, this.helper.cmd.ArgOption(key)}
}
func (this *CmdTabCommandOptionHelper) FlagOption(key string) *CmdTabCommandOptionHelper {
	return &CmdTabCommandOptionHelper{this.CmdTabCommandHelper, this.helper.cmd.FlagOption(key)}
}

//
// Option attributes
//
func (this *CmdTabCommandOptionHelper) Context(ctx string) *CmdTabCommandOptionHelper {
	this.helper.Context(ctx)
	return this
}

func (this *CmdTabCommandOptionHelper) Default(def interface{}) *CmdTabCommandOptionHelper {
	this.helper.Default(def)
	return this
}

func (this *CmdTabCommandOptionHelper) Short(short rune) *CmdTabCommandOptionHelper {
	this.helper.Short(short)
	return this
}

func (this *CmdTabCommandOptionHelper) Long(long string) *CmdTabCommandOptionHelper {
	this.helper.Long(long)
	return this
}

func (this *CmdTabCommandOptionHelper) Description(desc ...string) *CmdTabCommandOptionHelper {
	this.helper.Description(desc...)
	return this
}

func (this *CmdTabCommandOptionHelper) ArgDescription(desc string) *CmdTabCommandOptionHelper {
	this.helper.ArgDescription(desc)
	return this
}

func (this *CmdTabCommandOptionHelper) Args(n int) *CmdTabCommandOptionHelper {
	this.helper.Args(n)
	return this
}

func (this *CmdTabCommandOptionHelper) Array() *CmdTabCommandOptionHelper {
	this.helper.Array()
	return this
}

func (this *CmdTabCommandOptionHelper) Env(name string) *CmdTabCommandOptionHelper {
	this.helper.Env(name)
	return this
}

///////////////////////////////////////////////////////////////////////////
// General Configurable Command

type CmdTabConfigurableCommandHelper struct {
	cmdtab ConfigurableCmdTab
}

var _ ConfigurableCommand = &CmdTabConfigurableCommandHelper{}

func (this *CmdTabConfigurableCommandHelper) AsCmdTab() ConfigurableCmdTab {
	return this.cmdtab
}
func (this *CmdTabConfigurableCommandHelper) AsCommand() ConfigurableCommand {
	return this
}

func (this *CmdTabConfigurableCommandHelper) Execute(ctx *Options, args []string) error {
	return this.cmdtab.Execute(ctx, args)
}
func (this *CmdTabConfigurableCommandHelper) GetCmdDescription() string {
	return this.cmdtab.GetCmdDescription()
}
func (this *CmdTabConfigurableCommandHelper) GetCmdArgDescription() string {
	return this.cmdtab.GetCmdArgDescription()
}
func (this *CmdTabConfigurableCommandHelper) GetName() string {
	return this.cmdtab.GetName()
}
func (this *CmdTabConfigurableCommandHelper) Help(name string, args []string) {
	this.cmdtab.Help(name, args)
}
func (this *CmdTabConfigurableCommandHelper) GetOptions() OptionSpec {
	return this.cmdtab.GetOptions()
}

//
// Option Config
//

func (this *CmdTabConfigurableCommandHelper) ArgOption(name string) *ConfigurableCommandOptionHelper {
	return &ConfigurableCommandOptionHelper{this, this.cmdtab.ArgOption(name).helper}
}

func (this *CmdTabConfigurableCommandHelper) FlagOption(name string) *ConfigurableCommandOptionHelper {
	return &ConfigurableCommandOptionHelper{this, this.cmdtab.FlagOption(name).helper}
}
