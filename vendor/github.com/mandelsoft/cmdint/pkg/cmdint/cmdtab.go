package cmdint

import (
	"fmt"
	"sort"
	"unicode/utf8"
)

type CatchFunction func(cmdtab CmdTab, opts *Options) error

type CmdTab interface {
	GetName() string
	GetDefaultFunction() CommandFunction
	GetCommand(name string) Command
	Command
}
type _CmdTab struct {
	name       string
	optionspec OptionSpec
	commands   map[string]*CommandEntry
	argdesc    string
	desc       string
	setup      CommandFunction
	deffunc    CommandFunction
	catchfunc  CatchFunction
}

type CommandEntry struct {
	command Command
}

var _ CmdTab = &_CmdTab{}
var _ Command = &_CmdTab{}

func NewCmdTab(name string, spec ...OptionSpec) *_CmdTab {
	if len(spec) > 1 {
		panic(fmt.Errorf("only one option spec argument allowed"))
	}
	s := NewOptionSpec()
	if len(spec) > 0 && spec[0] != nil {
		s = spec[0]
	}
	return &_CmdTab{name, s, map[string]*CommandEntry{}, "", "", nil, nil, nil}
}

func (this *_CmdTab) GetName() string {
	return this.name
}
func (this *_CmdTab) GetCmdDescription() string {
	return this.desc
}
func (this *_CmdTab) GetCmdArgDescription() string {
	return this.argdesc
}
func (this *_CmdTab) GetOptions() OptionSpec {
	return this.optionspec
}
func (this *_CmdTab) GetDefaultFunction() CommandFunction {
	return this.deffunc
}
func (this *_CmdTab) GetCommand(name string) Command {
	e, ok := this.commands[name]
	if ok {
		return e.command
	}
	return nil
}

func (this *_CmdTab) Execute(ctx *Options, args []string) error {
	if ctx == nil {
		ctx = NewOptions(ctx)
		ctx.Command = this.GetName()
		ctx.Arguments = append([]string{this.GetName()}, args...)
	}
	if len(args) > 0 {
		if args[0] == "help" || args[0] == "--help" || args[0] == "-?" {
			this.help(ctx, args[1:])
			return nil
		}
	}
	opts, err := this.optionspec.Parse(ctx, args)
	if err != nil {
		return err
	}

	if this.setup != nil {
		err := this.setup(opts)
		if err != nil {
			return err
		}
	}
	if len(opts.Arguments) == 0 {
		if this.catchfunc != nil {
			return this.catchfunc(this, opts)
		}
		if this.deffunc != nil {
			return this.deffunc(opts)
		} else {
			return fmt.Errorf("no command specified")
		}
	}
	// fmt.Printf("lookup command %s: %v\n", opts.Arguments[0], this.commands)
	n, c := this.determineCommand(opts.Arguments[0])
	if c == nil {
		if this.catchfunc != nil {
			return this.catchfunc(this, opts)
		}
		return fmt.Errorf("unknown command '%s'", opts.Arguments[0])
	}
	// fmt.Printf("calling command %s with %v\n", n, opts.Arguments[1:])
	opts.Command = n
	err = c.command.Execute(opts, opts.Arguments[1:])
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %s", n, err)
}

func (this *_CmdTab) help(ctx *Options, args []string) {
	if ctx != nil {
		this.Help(ctx.Command, args)
	} else {
		this.Help(this.GetName(), args)
	}
}

func (this *_CmdTab) Help(name string, args []string) {
	d := DecodeDescription(this)
	a := this.GetCmdArgDescription()

	fmt.Printf("Synopsis: %s %s - %s\n", name, a, d[0])

	desc := this.optionspec.GetOptionHelp()
	if len(desc) > 0 {
		fmt.Printf("\nOptions:\n")
		fmt.Printf("%s", desc)
	}
	fmt.Printf("\nCommands:\n")
	keys := []string{}
	for k, _ := range this.commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	max_s := 0
	max_n := 0
	for n, c := range this.commands {
		l := utf8.RuneCountInString(n)
		if l > max_n {
			max_n = l
		}
		a := c.command.GetCmdArgDescription()
		l = utf8.RuneCountInString(a)
		if l > max_s {
			max_s = l
		}
	}
	for _, n := range keys {
		c := this.commands[n]
		d := DecodeDescription(c.command)
		a := c.command.GetCmdArgDescription()
		fmt.Printf("  %-*s %-*s %s\n", max_n, n, max_s, a, d[0])
	}

	if len(d) > 1 {
		fmt.Println()
		for _, t := range d[1:] {
			fmt.Printf("%s\n", t)
		}
	}

	fmt.Println()
	for _, n := range keys {
		c := this.commands[n]
		d := DecodeDescription(c.command)
		if len(d) > 1 {
			a := c.command.GetCmdArgDescription()
			fmt.Printf("\n%s %s %s\n", n, a, d[0])
			for _, l := range d[1:] {
				fmt.Printf("  %s\n", l)
			}
		}
	}

}

func (this *_CmdTab) determineCommand(name string) (string, *CommandEntry) {
	c, ok := this.commands[name]
	if ok {
		return name, c
	}
	min := -1
	f := ""
	for n, e := range this.commands {
		d := Levenshtein(n, name)
		if d < len(n)/2 && (min == -1 || min > d) {
			f, c, min = n, e, d
		} else {
			if d == min {
				c = nil
			}
		}
	}
	if c == nil {
		for n, e := range this.commands {
			if match_omit(n, name) {
				return n, e
			}
		}
	}
	return f, c
}
