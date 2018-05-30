package cmdint

import (
	"bufio"
	"fmt"
	"os"
	"unicode/utf8"
)

type ScriptSource interface {
	IsInteractive() bool
	NextLine() *string
}

type CmdInt struct {
	cmdtab CmdTab
	ctx    *Options
}

func NewCmdInt(cmdtab CmdTab, args []string) *CmdInt {
	cmdint := &CmdInt{cmdtab, nil}
	if args != nil && len(args) > 0 {
		//		cmdint.ctx = cmdtab.ParseOptions()
	}
	return cmdint
}

func (this CmdInt) Execute(input ScriptSource) error {
	line := ""
	no := 0
	for l := input.NextLine(); l != nil; l = input.NextLine() {
		no++
		r, s := utf8.DecodeLastRuneInString(*l)
		if r == '\\' {
			line = line + (*l)[0:len(*l)-s]
		} else {
			line = line + *l
			args, err := Split(line)
			if err == nil {
				err = this.cmdtab.Execute(this.ctx, args)
			}
			if err != nil {
				return fmt.Errorf("line %d: %s", no, err)
			}
		}
	}

	return nil
}

func Split(line string) ([]string, error) {
	var arg *string = nil
	empty := ""
	args := []string{}
	mask := false
	quote := false

	for _, c := range line {
		if mask {
			switch c {
			case '"' | '\\' | ' ' | '\t':
			default:
				return nil, fmt.Errorf("character '%s' cannot be masked", string(c))
			}
		} else {
			switch c {
			case '"':
				quote = !quote
				if quote && arg == nil {
					arg = &empty
				}
				continue
			case ' ' | '\t':
				if arg != nil {
					args = append(args, *arg)
					arg = nil
				}
				continue
			case '\\':
				mask = true
				continue
			}
		}
		t := string(c)
		if arg != nil {
			t = *arg + t
		}
		arg = &t
	}
	if arg != nil {
		args = append(args, *arg)
	}
	return args, nil
}

type FileSource struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewFileSource(f *os.File) ScriptSource {
	scanner := bufio.NewScanner(f)
	return &FileSource{f, scanner}
}

func (this *FileSource) IsInteractive() bool {
	return false
}

func (this *FileSource) NextLine() *string {
	if this.scanner.Scan() {
		t := this.scanner.Text()
		return &t
	}
	return nil
}
