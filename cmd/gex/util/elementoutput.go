package util

import (
	"github.com/afritzler/garden-examiner/cmd/gex/context"
)

type ElementOutput struct {
	Elems []interface{}
}

func NewElementOutput() *ElementOutput {
	return &ElementOutput{[]interface{}{}}
}

func (this *ElementOutput) Add(ctx *context.Context, e interface{}) error {
	this.Elems = append(this.Elems, e)
	return nil
}

func (this *ElementOutput) Out(ctx *context.Context) {
}
