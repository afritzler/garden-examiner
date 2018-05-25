package output

import (
	"fmt"

	"github.com/afritzler/garden-examiner/cmd/gex/context"
)

type SingleElementOutput struct {
	Elem interface{}
}

var _ Output = &SingleElementOutput{}

func NewSingleElementOutput() *SingleElementOutput {
	return &SingleElementOutput{}
}

func (this *SingleElementOutput) Add(ctx *context.Context, e interface{}) error {
	if this.Elem == nil {
		this.Elem = e
		return nil
	}
	return fmt.Errorf("only one element can be selected, but multiple elements selected/found")
}

func (this *SingleElementOutput) Close(ctx *context.Context) error {
	return nil
}

func (this *SingleElementOutput) Out(ctx *context.Context) error {
	return nil
}
