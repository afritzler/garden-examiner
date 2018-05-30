package cmdint

import (
	"fmt"
)

type SimpleCommand struct {
	cmd        CommandFunction
	optionspec OptionSpec
	argdesc    string
	desc       string
}

var _ Command = &SimpleCommand{}

func NewCommand(cmd CommandFunction) *SimpleCommand {
	return &SimpleCommand{cmd, NewOptionSpec(), "", ""}
}

func (p *SimpleCommand) GetCmdArgDescription() string {
	return p.argdesc
}
func (p *SimpleCommand) GetCmdDescription() string {
	return p.desc
}
func (p *SimpleCommand) GetOptions() OptionSpec {
	return p.optionspec
}

func (this *SimpleCommand) Execute(ctx *Options, args []string) error {
	if len(args) > 0 {
		if args[0] == "help" || args[0] == "--help" || args[0] == "-?" {
			this.Help(ctx.Command, args[1:])
			return nil
		}
	}
	opts, err := this.optionspec.Parse(ctx, args)
	if err != nil {
		return err
	}
	return this.cmd(opts)
}

func (this *SimpleCommand) Help(name string, args []string) {
	d := DecodeDescription(this)
	a := this.GetCmdArgDescription()
	if a == "" {
		a = this.optionspec.GetArgDescription()
	}
	fmt.Printf("Synopsis: %s %s - %s\n", name, a, d[0])
	desc := this.optionspec.GetOptionHelp()
	if len(desc) > 0 {
		fmt.Printf("\nOptions:\n")
		fmt.Printf("%s", desc)
	}
	if len(d) > 1 {
		fmt.Println()
		for _, t := range d[1:] {
			fmt.Printf("%s\n", t)
		}
	}
}
