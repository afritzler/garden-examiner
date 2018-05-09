package util

import (
	"fmt"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/context"
	_ "github.com/afritzler/garden-examiner/pkg"
)

type Handler interface {
	RequireScan(string) bool
	MatchName(interface{}, string) (bool, error)
	Get(*context.Context, string) (interface{}, error)
	Iterator(ctx *context.Context, opts *cmdint.Options) (Iterator, error)
	Match(*context.Context, interface{}, *cmdint.Options) (bool, error)
	Add(*context.Context, interface{}) error
	Out(*context.Context)
}

func Doit(opts *cmdint.Options, h Handler) error {
	ctx := context.Get(opts)

	if len(opts.Arguments) > 0 && (len(opts.Arguments) != 1 || opts.Arguments[0] != "all") {
		return doDedicated(ctx, opts, h)
	} else {
		return doAll(ctx, opts, h)
	}
}

func doAll(ctx *context.Context, opts *cmdint.Options, h Handler) error {
	i, err := h.Iterator(ctx, opts)
	if err != nil {
		return err
	}
	for i.HasNext() {
		e := i.Next()
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
	}
	h.Out(ctx)
	return nil
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
					i.Reset()
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
	h.Out(ctx)
	return nil
}
