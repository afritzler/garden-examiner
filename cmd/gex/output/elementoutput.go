package output

import (
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	. "github.com/afritzler/garden-examiner/pkg/data"
)

type ElementOutput struct {
	source ProcessingSource
	Elems  Iterable
}

func NewElementOutput(chain ProcessChain) *ElementOutput {
	return (&ElementOutput{}).new(chain)
}

func (this *ElementOutput) new(chain ProcessChain) *ElementOutput {
	this.source = NewIncrementalProcessingSource()
	if chain == nil {
		this.Elems = this.source
	} else {
		this.Elems = Process(this.source).Asynchronously().Apply(chain)
	}
	return this
}

func (this *ElementOutput) Add(ctx *context.Context, e interface{}) error {
	this.source.Add(e)
	return nil
}

func (this *ElementOutput) Close(ctx *context.Context) error {
	this.source.Close()
	return nil
}

func (this *ElementOutput) Out(ctx *context.Context) {
}
