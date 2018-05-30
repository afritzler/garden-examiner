package cmdint

import (
	"strings"
	"unicode/utf8"
)

func (this *_OptionSpec) Mixed() OptionSpec {
	this.mode = OPTION_MODE_MIXED
	return this
}

func (this *_OptionSpec) Raw() OptionSpec {
	this.mode = OPTION_MODE_RAW
	return this
}

func (this *_OptionSpec) newOption(key string, flag bool) *Option {
	if len(key) == 0 {
		panic("option requires non empty key")
	}
	if strings.Index(key, " ") > 0 {
		panic("option key may not contain spaces")
	}

	long := ""
	short := ""

	r := '\u0000'
	t, size := utf8.DecodeRuneInString(key)
	if len(key) == size {
		short = key
		r = t
	} else {
		long = key
	}

	option := &Option{
		Key:   key,
		Long:  long,
		Short: short,
	}
	if !flag {
		option.Args = 1
	} else {
		option.Default = false
	}

	if short != "" {
		this.short[r] = option
	}
	if long != "" {
		option.Key = key
		this.long[long] = option
	}
	this.options[option.Key] = option
	return option
}

//////////////////////////////////////////////////////////////////////////
// Configuration
//////////////////////////////////////////////////////////////////////////

//////////////////////////////////////////////////////////////////////////
// OptionSpec Config

func (this *_OptionSpec) ArgOption(key string) *OptionConfigHelper {
	opt := this.newOption(key, false)
	return &OptionConfigHelper{opt, this}
}

func (this *_OptionSpec) FlagOption(key string) *OptionConfigHelper {
	opt := this.newOption(key, true)
	return &OptionConfigHelper{opt, this}
}

//////////////////////////////////////////////////////////////////////////
// Option Config

type OptionConfigHelper struct {
	option *Option
	spec   *_OptionSpec
}

var _ OptionSpec = &OptionConfigHelper{}

func (this *OptionConfigHelper) GetOptions() map[string]*Option {
	return this.spec.GetOptions()
}
func (this *OptionConfigHelper) Get(key string) *Option {
	return this.spec.Get(key)
}

func (this *OptionConfigHelper) Mixed() OptionSpec {
	return this.spec.Mixed()
}
func (this *OptionConfigHelper) Raw() OptionSpec {
	return this.spec.Raw()
}

func (this *OptionConfigHelper) Parse(ctx *Options, args []string) (*Options, error) {
	return this.spec.Parse(ctx, args)
}

func (this *OptionConfigHelper) GetArgDescription() string {
	return this.spec.GetArgDescription()
}

func (this *OptionConfigHelper) GetOptionHelp() string {
	return this.spec.GetOptionHelp()
}

//
// Add option
//
func (this *OptionConfigHelper) ArgOption(key string) *OptionConfigHelper {
	opt := this.spec.newOption(key, false)
	return &OptionConfigHelper{opt, this.spec}
}

func (this *OptionConfigHelper) FlagOption(key string) *OptionConfigHelper {
	opt := this.spec.newOption(key, true)
	return &OptionConfigHelper{opt, this.spec}
}

//
// Option attributes
//
func (this *OptionConfigHelper) Context(ctx string) *OptionConfigHelper {
	this.option.SetContext(ctx)
	return this
}

func (this *OptionConfigHelper) Default(def interface{}) *OptionConfigHelper {
	this.option.SetDefault(def)
	return this
}

func (this *OptionConfigHelper) Short(short rune) *OptionConfigHelper {
	if this.option.Short != "" {
		old, _ := utf8.DecodeRuneInString(this.option.Short)
		delete(this.spec.short, old)
	}
	this.spec.short[short] = this.option
	this.option.SetShort(short)
	return this
}

func (this *OptionConfigHelper) Long(long string) *OptionConfigHelper {
	if this.option.Long != "" {
		delete(this.spec.long, this.option.Long)
	}
	this.spec.long[long] = this.option
	this.option.SetLong(long)
	return this
}

func (this *OptionConfigHelper) ArgDescription(desc string) *OptionConfigHelper {
	this.option.SetArgDescription(desc)
	return this
}

func (this *OptionConfigHelper) Description(desc ...string) *OptionConfigHelper {
	this.option.SetDescription(desc...)
	return this
}

func (this *OptionConfigHelper) Args(n int) *OptionConfigHelper {
	this.option.SetArgs(n)
	return this
}

func (this *OptionConfigHelper) Array() *OptionConfigHelper {
	this.option.SetArray()
	return this
}

func (this *OptionConfigHelper) Env(name string) *OptionConfigHelper {
	this.option.SetEnv(name)
	return this
}
