package cmdint

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode/utf8"
)

//////////////////////////////////////////////////////////////////////////
// Option Spec Set
//////////////////////////////////////////////////////////////////////////

type OptionSpec interface {
	Parse(opts *Options, args []string) (*Options, error)

	Mixed() OptionSpec
	Raw() OptionSpec
	ArgOption(key string) *OptionConfigHelper
	FlagOption(key string) *OptionConfigHelper

	GetArgDescription() string
	GetOptionHelp() string

	GetOptions() map[string]*Option
	Get(key string) *Option
}

const (
	OPTION_MODE_PARSE = "parse"
	OPTION_MODE_MIXED = "mixed"
	OPTION_MODE_RAW   = "raw"
)

type _OptionSpec struct {
	options map[string]*Option
	short   map[rune]*Option
	long    map[string]*Option
	mode    string
}

var _ OptionSpec = &_OptionSpec{}

func NewOptionSpec() OptionSpec {
	return &_OptionSpec{map[string]*Option{}, map[rune]*Option{}, map[string]*Option{}, OPTION_MODE_PARSE}
}

func (this *_OptionSpec) Get(key string) *Option {
	return this.options[key]
}

func (this *_OptionSpec) GetOptions() map[string]*Option {
	return this.options
}

func (this *_OptionSpec) GetArgDescription() string {
	if len(this.options) > 0 {
		return "<options>"
	}
	return ""
}
func (this *_OptionSpec) GetOptionHelp() string {
	keys := []string{}
	desc := ""
	for k, _ := range this.options {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	max := 0
	for _, k := range keys {
		o := this.options[k]
		a := o.GetArgDescription()
		l := utf8.RuneCountInString(strings.TrimSpace(o.longOption()+" "+o.shortOption())) + 1 + utf8.RuneCountInString(a)
		if l > max {
			max = l
		}
	}
	for _, k := range keys {
		o := this.options[k]
		d := DecodeDescription(o)
		a := o.GetArgDescription()
		for i, t := range d {
			if i == 0 {
				desc += fmt.Sprintf("  %-*s %s\n", max, strings.TrimSpace(o.longOption()+" "+o.shortOption())+" "+a, t)
			} else {
				desc += fmt.Sprintf("  %-*s %s\n", max, "", t)
			}
		}
	}
	return desc
}

//////////////////////////////////////////////////////////////////////////
// Argument parsing
//////////////////////////////////////////////////////////////////////////

func (this *_OptionSpec) Parse(ctx *Options, args []string) (*Options, error) {
	var err error
	cur := 0
	options := NewOptions(nil)
	for cur < len(args) {
		arg := args[cur]
		if len(arg) > 1 && arg[0] == '-' {
			if arg[1] == '-' {
				if len(arg) == 2 {
					cur++
					break
				}
				cur, err = this.parseLongOption(arg[2:], cur+1, args, options)
				if err != nil {
					return options, err
				}
			} else {
				if this.mode == OPTION_MODE_RAW {
					break
				}
				// check for short option direct argument assignment
				i := strings.Index(arg, "=")
				if i > 0 {
					if i != 2 {
						return options, fmt.Errorf("invalid short argument assignment '%s'", arg)
					}
					cur, err = this.parseLongOption(arg[1:], cur+1, args, options)
					if err != nil {
						return options, err
					}
				} else {
					// parse short options
					cur++
					for _, o := range arg[1:] {
						cur, err = this.parseShortOption(o, cur, args, options)
						if err != nil {
							return options, err
						}
					}
				}
			}
		} else {
			if this.mode != OPTION_MODE_MIXED {
				break
			}
			options.Arguments = append(options.Arguments, args[cur])
			cur++
		}
	}
	if this.mode == OPTION_MODE_RAW {
		options.Raw = true
	}
	if cur < len(args) {
		options.Arguments = append(options.Arguments, args[cur:]...)
	}

	options.Parent = ctx
	if ctx != nil {
		options.Defaulted(this, ctx)
		options.Context = ctx.Context
	}

	env := os.Environ()

	for _, o := range this.options {
		if _, ok := options.Flags[o.Key]; !ok {
			options.Flags[o.Key] = false
		}
		if o.Env != "" {
			v := lookup(env, o.Env)
			if v != nil {
				switch {
				case o.Args == 0:
					if !options.Flags[o.Key] {
						switch strings.ToLower(*v) {
						case "true", "1":
							options.Flags[o.Key] = true
						}
					}
				case o.Args == 1:
					if _, ok := options.SingleArgumentOptions[o.Key]; !ok {
						options.SingleArgumentOptions[o.Key] = *v
					}
					if _, ok := options.SingleArgumentArrayOptions[o.Key]; !ok {
						options.SingleArgumentArrayOptions[o.Key] = []string{*v}
					}
				case o.Args > 1:
					args := strings.Split(*v, ",")
					if len(args) == o.Args {
						if _, ok := options.SingleArgumentArrayOptions[o.Key]; !ok {
							options.SingleArgumentArrayOptions[o.Key] = args
						}
						if _, ok := options.MultiArgumentArrayOptions[o.Key]; !ok {
							options.MultiArgumentArrayOptions[o.Key] = [][]string{args}
						}
					}
				}
			}
		}
	}
	return options, nil
}

func lookup(env []string, key string) *string {
	key = key + "="
	for _, e := range env {
		if strings.HasPrefix(e, key) {
			s := e[len(key):]
			return &s
		}
	}
	return nil
}

func (this *_OptionSpec) parseLongOption(name string, cur int, args []string, options *Options) (int, error) {
	optargs := []string{}

	i := strings.Index(name, "=")
	if i > 0 {
		optargs = append(optargs, name[i+1:])
		name = name[0:i]
	}

	option, ok := this.long[name]
	if !ok {
		r, size := utf8.DecodeRuneInString(name)
		if len(name) == size {
			option, ok = this.short[r]
		}
		if !ok {
			if this.mode == OPTION_MODE_RAW {
				options.Arguments = append(options.Arguments, args[cur])
				return cur + 1, nil
			}
			return cur, fmt.Errorf("unknown option '%s'", name)
		}
	}
	return this.parseOption(name, option, optargs, cur, args, options)
}

func (this *_OptionSpec) parseOption(name string, option *Option, optargs []string,
	cur int, args []string, options *Options) (int, error) {
	fmt.Printf("parse option %s: %v\n", option.Key, optargs)
	if option.Args > 0 {
		if len(optargs) == 0 {
			if len(args) < option.Args+cur {
				optargs = args[cur:]
				cur = len(args)
			} else {
				optargs = args[cur : cur+option.Args]
				cur += option.Args
			}
		}
		if option.Args != len(optargs) && len(optargs) == 1 {
			optargs = strings.Split(optargs[0], ",")
		}
		if option.Args != len(optargs) {
			return cur, fmt.Errorf("option '%s' requires %d argument(s) (have %d)", name, option.Args, len(optargs))
		}
	}
	switch {
	case option.List && option.Args == 1:
		result, ok := options.SingleArgumentArrayOptions[option.Key]
		if !ok {
			result = []string{}
		}
		result = append(result, optargs[0])
		options.SingleArgumentArrayOptions[option.Key] = result

	case option.List && option.Args > 1:
		result, ok := options.MultiArgumentArrayOptions[option.Key]
		if !ok {
			result = [][]string{}
		}
		result = append(result, optargs)
		options.MultiArgumentArrayOptions[option.Key] = result

	case !option.List && option.Args == 1:
		_, ok := options.SingleArgumentOptions[option.Key]
		if ok {
			return cur, fmt.Errorf("multiple option '%s' given", name)
		}
		options.SingleArgumentOptions[option.Key] = optargs[0]
	case !option.List && option.Args > 1:
		_, ok := options.SingleArgumentArrayOptions[option.Key]
		if ok {
			return cur, fmt.Errorf("multiple option '%s' given", name)
		}
		options.SingleArgumentArrayOptions[option.Key] = optargs
	case option.Args == 0:
		options.Flags[option.Key] = true
	default:
		return cur, fmt.Errorf("unknown option kind %+v", option)
	}
	return cur, nil
}

func (this *_OptionSpec) parseShortOption(name rune, cur int, args []string, options *Options) (int, error) {
	optargs := []string{}
	option, ok := this.short[name]
	if !ok {
		if this.mode == OPTION_MODE_RAW {
			return cur, nil
		}
		return cur, fmt.Errorf("unknown option '%s'", string(name))
	}
	return this.parseOption(string(name), option, optargs, cur, args, options)
}
