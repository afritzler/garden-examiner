package cmdint

import (
	"encoding/json"
	"fmt"
	"os"
)

type Options struct {
	Raw                        bool
	Parent                     *Options
	Command                    string
	Arguments                  []string
	Flags                      map[string]bool
	SingleArgumentOptions      map[string]string
	MultiArgumentOptions       map[string][]string
	SingleArgumentArrayOptions map[string][]string
	MultiArgumentArrayOptions  map[string][][]string
	Context                    interface{}
}

func NewOptions(parent *Options) *Options {
	return &Options{
		Raw:       false,
		Parent:    parent,
		Arguments: []string{},
		Flags:     map[string]bool{},
		SingleArgumentOptions:      map[string]string{},
		MultiArgumentOptions:       map[string][]string{},
		SingleArgumentArrayOptions: map[string][]string{},
		MultiArgumentArrayOptions:  map[string][][]string{},
	}
}

func (o *Options) IsFlag(key string) bool {
	return o.Flags[key]
}
func (o *Options) GetOptionValue(key string) *string {
	v, ok := o.SingleArgumentOptions[key]
	if ok {
		return &v
	}
	return nil
}
func (o *Options) GetOptionValues(key string) []string {
	return o.MultiArgumentOptions[key]
}
func (o *Options) GetArrayOptionValue(key string) []string {
	return o.SingleArgumentArrayOptions[key]
}
func (o *Options) GetArrayOptionValues(key string) [][]string {
	return o.MultiArgumentArrayOptions[key]
}
func (o *Options) GetPositionalArguments(key string) []string {
	return o.Arguments
}

func (o *Options) GetContextOptionValue(key string) *string {
	v := o.GetOptionValue(key)
	if v != nil {
		return v
	}
	if o.Parent != nil {
		return o.Parent.GetContextOptionValue(key)
	}
	return nil
}
func (o *Options) GetContextOptionValues(key string) []string {
	v := o.GetOptionValues(key)
	if v != nil {
		return v
	}
	if o.Parent != nil {
		return o.Parent.GetContextOptionValues(key)
	}
	return nil
}
func (o *Options) GetContextArrayOptionValue(key string) []string {
	v := o.GetArrayOptionValue(key)
	if v != nil {
		return v
	}
	if o.Parent != nil {
		return o.Parent.GetContextArrayOptionValue(key)
	}
	return nil
}
func (o *Options) GetContextArrayOptionValues(key string) [][]string {
	v := o.GetArrayOptionValues(key)
	if v != nil {
		return v
	}
	if o.Parent != nil {
		return o.Parent.GetContextArrayOptionValues(key)
	}
	return nil
}

func (o *Options) Errorf(args ...interface{}) {
	p := o
	t := ""
	for p != nil {
		if p.Command != "" {
			t = fmt.Sprintf("%s: %s", p.Command, t)
		}
		p = p.Parent
	}
	if len(args) == 1 {
		fmt.Fprintf(os.Stderr, "ERROR: %s: %v\n", t, args[0])
	} else {
		fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", t, fmt.Sprintf(args[0].(string), args[1:]...))
	}
	os.Exit(1)
}

func (this *Options) AsJson() string {
	b, _ := json.MarshalIndent(this, "", "  ")
	return string(b)
}

func (this *Options) Defaulted(spec OptionSpec, defs *Options) *Options {
	for n, b := range defs.Flags {
		if _, ok := this.Flags[n]; b || !ok {
			this.Flags[n] = b
		}
	}
	for n, v := range defs.SingleArgumentOptions {
		if _, ok := this.SingleArgumentOptions[n]; !ok {
			this.SingleArgumentOptions[n] = v
		}
	}
	for n, v := range defs.MultiArgumentOptions {
		if _, ok := this.MultiArgumentOptions[n]; !ok {
			this.MultiArgumentOptions[n] = v
		}
	}
	for n, v := range defs.SingleArgumentArrayOptions {
		if _, ok := this.SingleArgumentArrayOptions[n]; !ok {
			this.SingleArgumentArrayOptions[n] = v
		}
	}
	for n, v := range defs.MultiArgumentArrayOptions {
		if _, ok := this.MultiArgumentArrayOptions[n]; !ok {
			this.MultiArgumentArrayOptions[n] = v
		}
	}

	for _, o := range spec.GetOptions() {
		o.PropagateDefault(this, defs)
	}
	return this
}
