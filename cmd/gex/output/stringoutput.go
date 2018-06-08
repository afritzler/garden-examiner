package output

import (
	"fmt"
	"strings"

	"github.com/afritzler/garden-examiner/pkg/data"

	"github.com/afritzler/garden-examiner/cmd/gex/context"
)

type StringOutput struct {
	ElementOutput
	linesep string
}

var _ Output = &StringOutput{}

func NewStringOutput(mapper data.MappingFunction, linesep string) *StringOutput {
	return (&StringOutput{}).new(mapper, linesep)
}

func (this *StringOutput) new(mapper data.MappingFunction, lineseperator string) *StringOutput {
	this.linesep = lineseperator
	this.ElementOutput.new(data.Chain().Parallel(20).Map(mapper))
	return this
}

func (this *StringOutput) Out(ctx *context.Context) error {
	var err error = nil
	i := this.Elems.Iterator()
	for i.HasNext() {
		switch cfg := i.Next().(type) {
		case error:
			err = cfg
			if this.linesep == "" {
				fmt.Printf("Error: %s\n", err)
			} else {
				fmt.Printf("%s\nError: %s\n", this.linesep, err)
			}
		case string:
			if cfg != "" {
				if this.linesep != "" {
					if !strings.HasPrefix(cfg, this.linesep+"\n") {
						fmt.Println(this.linesep)
					}
				}
				fmt.Println(cfg)
			}
		}
	}
	return err
}
