package cmdint

import (
	"fmt"
)

//////////////////////////////////////////////////////////////////////////
// Option specification
//////////////////////////////////////////////////////////////////////////

type Option struct {
	Key     string
	Long    string
	Short   string
	Args    int
	List    bool
	Env     string
	Syn     string
	Desc    string
	Context string
	Default interface{}
}

func (o *Option) GetArgDescription() string {
	if o.Syn == "" {
		if o.Args > 0 {
			for i := 1; i <= o.Args; i++ {
				o.Syn = fmt.Sprintf(" <arg%d>", i)
			}
			o.Syn = o.Syn[1:]
		}
	}
	return o.Syn
}

func (o *Option) GetDescription() string {
	return o.Desc
}

func (o *Option) PropagateDefault(opts *Options, ctx *Options) {
	if o.Context != "" {
		switch {
		case o.Args == 0:
			if _, ok := opts.Flags[o.Key]; !ok {
				if v := ctx.IsFlag(o.Context); v {
					opts.Flags[o.Key] = v
				}
			}
		case o.Args == 1 && !o.List:
			if opts.GetOptionValue(o.Key) == nil {
				if v := ctx.GetContextOptionValue(o.Context); v != nil {
					opts.SingleArgumentOptions[o.Key] = *v
				}
			}
		case o.Args > 1 && !o.List:
			if opts.GetOptionValues(o.Key) == nil {
				if v := ctx.GetContextOptionValues(o.Context); v != nil {
					if len(v) == o.Args {
						opts.MultiArgumentOptions[o.Key] = v
					}
				}
			}
		case o.Args == 1 && o.List:
			if opts.GetArrayOptionValue(o.Key) == nil {
				if v := ctx.GetContextArrayOptionValue(o.Context); v != nil {
					opts.SingleArgumentArrayOptions[o.Key] = v
				} else {
					if v := ctx.GetContextOptionValue(o.Context); v != nil {
						opts.SingleArgumentArrayOptions[o.Key] = []string{*v}
					}
				}
			}
		case o.Args > 1 && o.List:
			if opts.GetArrayOptionValues(o.Key) == nil {
				if v := ctx.GetContextArrayOptionValues(o.Context); v != nil {
					if len(v) > 0 && len(v[0]) == o.Args {
						opts.MultiArgumentArrayOptions[o.Key] = v
					}
				} else {
					if v := ctx.GetContextArrayOptionValue(o.Context); v != nil {
						if len(v) == o.Args {
							opts.MultiArgumentArrayOptions[o.Key] = [][]string{v}
						}
					}
				}
			}
		default:
			panic(fmt.Errorf("unhandled option mode %#v", o))
		}
	}

	if o.Default != nil {
		switch {
		case o.Args == 0:
			if _, ok := opts.Flags[o.Key]; !ok {
				opts.Flags[o.Key] = o.Default.(bool)
			}
		case o.Args == 1 && !o.List:
			if opts.GetOptionValue(o.Key) == nil {
				opts.SingleArgumentOptions[o.Key] = o.Default.(string)
			}
		case o.Args > 1 && !o.List:
			if opts.GetOptionValues(o.Key) == nil {
				opts.MultiArgumentOptions[o.Key] = o.Default.([]string)
			}
		case o.Args == 1 && o.List:
			if opts.GetArrayOptionValue(o.Key) == nil {
				opts.SingleArgumentArrayOptions[o.Key] = o.Default.([]string)
			}
		case o.Args > 1 && o.List:
			if opts.GetArrayOptionValues(o.Key) == nil {
				opts.MultiArgumentArrayOptions[o.Key] = o.Default.([][]string)
			}
		default:
			panic(fmt.Errorf("unhandled option mode %#v", o))
		}
	}
}

func (o *Option) longOption() string {
	if o.Long != "" {
		return "--" + o.Long
	}
	return ""
}
func (o *Option) shortOption() string {
	if o.Short != "" {
		return "-" + o.Short
	}
	return ""
}

/////////////////////////////////////////////////////////////////////////////
func (this *Option) SetContext(ctx string) *Option {
	this.Context = ctx
	return this
}

func (this *Option) SetShort(short rune) *Option {
	this.Short = string(short)
	return this
}

func (this *Option) SetLong(long string) *Option {
	this.Long = long
	return this
}

func (this *Option) SetArgDescription(desc string) *Option {
	this.Syn = desc
	return this
}

func (this *Option) SetDescription(desc ...string) *Option {
	this.Desc = compact(desc)
	return this
}

func (this *Option) SetArgs(n int) *Option {
	if this.Args == 0 {
		panic("args can only be set for non-flag options")
	}
	if n <= 0 {
		panic("argument count must be greater than zero")
	}
	if this.Args != n && this.Default != nil {
		panic("argument count be set before default value")
	}
	this.Args = n
	return this
}

func (this *Option) SetArray() *Option {
	if this.Args == 0 {
		panic("argument array can only be set for non-flag options")
	}
	if this.List != true && this.Default != nil {
		panic("argument array be set before default value")
	}
	this.List = true
	return this
}

func (this *Option) SetEnv(name string) *Option {
	this.Env = name
	return this
}

func (this *Option) SetDefault(def interface{}) *Option {
	if def != nil {
		switch {
		case this.Args == 0:
			def = def.(bool)
		case this.Args == 1 && !this.List:
			def = def.(string)
		case this.Args == 1 && this.List:
			s, ok := def.(string)
			if ok {
				def = []string{s}
			} else {
				def = def.([]string)
			}
		case this.Args > 1 && !this.List:
			a := def.([]string)
			if len(a) != this.Args {
				panic(fmt.Errorf("argument size mismatch (%d != %d)", len(a), this.Args))
			}
		case this.Args > 1 && this.List:
			a, ok := def.([]string)
			if ok {
				if len(a) != this.Args {
					panic(fmt.Errorf("argument size mismatch (%d != %d)", len(a), this.Args))
				}
				def = [][]string{a}
			} else {
				a := def.([][]string)
				for i, e := range a {
					if len(e) != this.Args {
						panic(fmt.Errorf("argument size mismatch for entry %d (%d != %d)", i, len(a), this.Args))
					}
				}
			}
		default:
			panic(fmt.Errorf("oops: missing case %#v", this))
		}
	}
	this.Default = def
	return this
}
