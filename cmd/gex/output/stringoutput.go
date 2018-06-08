package output

import (
	"fmt"
	"strings"

	"github.com/afritzler/garden-examiner/pkg/data"

	"github.com/afritzler/garden-examiner/cmd/gex/context"
)

type StringOutput struct {
	ElementOutput
}

var _ Output = &StringOutput{}

func NewStringOutput(mapper data.MappingFunction) *StringOutput {
	return (&StringOutput{}).new(mapper)
}

func (this *StringOutput) new(mapper data.MappingFunction) *StringOutput {
	this.ElementOutput.new(data.Chain().Parallel(20).Map(mapper))
	return this
}

func (this *StringOutput) Out(ctx *context.Context) error {
	i := this.Elems.Iterator()
	for i.HasNext() {
		switch cfg := i.Next().(type) {
		case error:
			return cfg
		case string:
			if !strings.HasPrefix(cfg, "---\n") {
				fmt.Println("---")
			}
			fmt.Println(cfg)
		}
	}
	return nil
}
