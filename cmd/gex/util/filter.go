package util

import (
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/mandelsoft/cmdint/pkg/cmdint"
)

type Filter interface {
	Match(ctx *context.Context, elem interface{}, opts *cmdint.Options) (bool, error)
	AddOptions(cmd cmdint.ConfigurableCmdTabCommand) cmdint.ConfigurableCmdTabCommand
}

type Filters struct {
	filters []Filter
}

func NewFilters() *Filters {
	return &Filters{[]Filter{}}
}

var _ Filter = &Filters{}

func (this *Filters) Add(f Filter) *Filters {
	this.filters = append(this.filters, f)
	return this
}

func (this *Filters) Match(ctx *context.Context, elem interface{}, opts *cmdint.Options) (bool, error) {
	for _, f := range this.filters {
		ok, err := f.Match(ctx, elem, opts)
		if !ok || err != nil {
			return ok, err
		}
	}
	return true, nil
}

func (this *Filters) AddOptions(cmd cmdint.ConfigurableCmdTabCommand) cmdint.ConfigurableCmdTabCommand {
	for _, f := range this.filters {
		cmd = f.AddOptions(cmd)
	}
	return cmd
}
