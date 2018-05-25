package util

import (
	"fmt"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/context"
	_ "github.com/afritzler/garden-examiner/pkg"
	"github.com/afritzler/garden-examiner/pkg/data"
)

type ElementTypeHandler interface {
	RequireScan(string) bool
	MatchName(interface{}, string) (bool, error)
	Get(*context.Context, string) (interface{}, error)
	GetAll(ctx *context.Context, opts *cmdint.Options) ([]interface{}, error)
	GetFilter() Filter
	GetDefault(opts *cmdint.Options) *string
}

type Handler interface {
	GetDefault(opts *cmdint.Options) *string
	RequireScan(string) bool
	MatchName(interface{}, string) (bool, error)
	Get(*context.Context, string) (interface{}, error)
	Iterator(ctx *context.Context, opts *cmdint.Options) (data.Iterator, error)
	Match(*context.Context, interface{}, *cmdint.Options) (bool, error)
	Add(*context.Context, interface{}) error
	Close(*context.Context) error
	Out(*context.Context) error
}

/////////////////////////////////////////////////////////////////////////////
// Basic handler

type StandardHandler struct {
	output Output
	elems  data.IndexedAccess
	impl   ElementTypeHandler
}

func NewStandardOutputHandler(o Output, impl ElementTypeHandler) *StandardHandler {
	return (&StandardHandler{}).new(o, impl)
}

func (this *StandardHandler) new(o Output, impl ElementTypeHandler) *StandardHandler {
	this.output = o
	this.elems = nil
	this.impl = impl
	return this
}

func ExecuteOutput(opts *cmdint.Options, o Output, impl ElementTypeHandler) error {
	return NewStandardOutputHandler(o, impl).Doit(opts)
}
func ExecuteOutputRaw(option string, opts *cmdint.Options, o Output, impl ElementTypeHandler) error {
	return NewStandardOutputHandler(o, impl).DoitRaw(option, opts)
}

func ExecuteMode(opts *cmdint.Options, outs Outputs, impl ElementTypeHandler) error {
	o, err := outs.Create(opts)
	if err != nil {
		return err
	}
	return NewStandardOutputHandler(o, impl).Doit(opts)
}

func (this *StandardHandler) Doit(opts *cmdint.Options) error {
	return Doit(opts, this)
}
func (this *StandardHandler) DoitRaw(option string, opts *cmdint.Options) error {
	return DoitRaw(option, opts, this)
}

func (this *StandardHandler) GetDefault(opts *cmdint.Options) *string {
	return this.impl.GetDefault(opts)
}

func (this *StandardHandler) Iterator(ctx *context.Context, opts *cmdint.Options) (data.Iterator, error) {
	if this.elems == nil {
		elems, err := this.impl.GetAll(ctx, opts)
		if err != nil {
			return nil, err
		}
		this.elems = data.IndexedSliceAccess(elems)
	}
	return data.NewIndexedIterator(this.elems), nil
}

func (this *StandardHandler) RequireScan(name string) bool {
	return this.impl.RequireScan(name)
}
func (this *StandardHandler) MatchName(e interface{}, name string) (bool, error) {
	return this.impl.MatchName(e, name)
}
func (this *StandardHandler) Get(ctx *context.Context, name string) (interface{}, error) {
	return this.impl.Get(ctx, name)
}
func (this *StandardHandler) Match(ctx *context.Context, e interface{}, opts *cmdint.Options) (bool, error) {
	return this.impl.GetFilter().Match(ctx, e, opts)
}
func (this *StandardHandler) Add(ctx *context.Context, e interface{}) error {
	return this.output.Add(ctx, e)
}
func (this *StandardHandler) Close(ctx *context.Context) error {
	return this.output.Close(ctx)
}
func (this *StandardHandler) Out(ctx *context.Context) error {
	return this.output.Out(ctx)
}

/////////////////////////////////////////////////////////////////////////////
// Standard Command Logic

func DoitRaw(name_option string, opts *cmdint.Options, h Handler) error {
	ctx := context.Get(opts)
	name := ""
	if v := opts.GetOptionValue(name_option); v != nil {
		name = *v
	}
	if name == "" {
		if def := h.GetDefault(opts); def != nil {
			name = *def
		}
	}
	if name == "" {
		return fmt.Errorf("no element selected")
	}
	opts.Arguments = []string{name}
	return doDedicated(ctx, opts, h)
}

func Doit(opts *cmdint.Options, h Handler) error {
	ctx := context.Get(opts)

	if len(opts.Arguments) == 0 {
		if def := h.GetDefault(opts); def != nil {
			//fmt.Printf("DEFAULT: %s\n", *def)
			opts.Arguments = []string{*def}
		}
	}
	all := len(opts.Arguments) == 1 && opts.Arguments[0] == "all"
	if len(opts.Arguments) > 0 && !all {
		return doDedicated(ctx, opts, h)
	} else {
		return doAll(ctx, opts, h, !all)
	}
}

func doAll(ctx *context.Context, opts *cmdint.Options, h Handler, filter bool) error {
	i, err := h.Iterator(ctx, opts)
	if err != nil {
		return err
	}
	for i.HasNext() {
		ok := true
		e := i.Next()
		if filter {
			ok, err = h.Match(ctx, e, opts)
			if err != nil {
				return err
			}
		}
		if ok {
			err := h.Add(ctx, e)
			if err != nil {
				return err
			}
		}
	}
	h.Close(ctx)
	return h.Out(ctx)
}

func doDedicated(ctx *context.Context, opts *cmdint.Options, h Handler) error {
	for _, n := range opts.Arguments {
		if h.RequireScan(n) {
			i, err := h.Iterator(ctx, opts)
			if err != nil {
				return err
			}
			for _, n := range opts.Arguments {
				if !h.RequireScan(n) {
					e, err := h.Get(ctx, n)
					if err != nil {
						return err
					}
					if e == nil {
						return fmt.Errorf("'%s' not found", n)
					}
					ok, err := h.Match(ctx, e, opts)
					if err != nil {
						return err
					}
					if ok {
						err := h.Add(ctx, e)
						if err != nil {
							return err
						}
					}
				} else {
					//fmt.Printf("LOOKUP %s\n", n)
					found := false
					i, err = h.Iterator(ctx, opts)
					if err != nil {
						return err
					}
					for i.HasNext() {
						e := i.Next()
						ok, err := h.Match(ctx, e, opts)
						if err != nil {
							return err
						}
						//fmt.Printf("  check %s: %s\n", e.(gube.Shoot).GetName(), ok)
						if ok {
							ok, err := h.MatchName(e, n)
							if err != nil {
								return err
							}
							if ok {
								err := h.Add(ctx, e)
								if err != nil {
									return err
								}
								found = true
							}
						}
					}
					if !found {
						return fmt.Errorf("'%s' not found", n)
					}
				}
			}
			h.Out(ctx)
			return nil
		}
	}

	for _, n := range opts.Arguments {
		e, err := h.Get(ctx, n)
		if err != nil {
			return err
		}
		err = h.Add(ctx, e)
		if err != nil {
			return err
		}
	}
	h.Close(ctx)
	return h.Out(ctx)
}
