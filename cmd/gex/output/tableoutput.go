package output

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
	. "github.com/afritzler/garden-examiner/pkg/data"
)

type TableProcessingOutput struct {
	ElementOutput
	header []string
	opts   *cmdint.Options
}

var _ Output = &TableProcessingOutput{}

func NewProcessingTableOutput(opts *cmdint.Options, chain ProcessChain, header ...string) *TableProcessingOutput {
	return (&TableProcessingOutput{}).new(opts, chain, header)
}

func (this *TableProcessingOutput) new(opts *cmdint.Options, chain ProcessChain, header []string) *TableProcessingOutput {
	this.header = header
	this.ElementOutput.new(chain)
	this.opts = opts
	return this
}

func (this *TableProcessingOutput) Out(*context.Context) error {
	lines := [][]string{this.header}

	sort := this.opts.GetArrayOptionValue(constants.O_SORT)
	slice := Slice(this.Elems)
	if sort != nil {
		cols := make([]string, len(this.header))
		idxs := map[string]int{}
		for i, n := range this.header {
			cols[i] = strings.ToLower(n)
			if strings.HasPrefix(cols[i], "-") {
				cols[i] = cols[i][1:]
			}
			idxs[cols[i]] = i
		}
		for _, k := range sort {
			key, _ := cmdint.SelectBest(strings.ToLower(k), cols...)
			if key == "" {
				return fmt.Errorf("unknown field '%s'", k)
			}
			slice.Sort(compare_column(idxs[key]))
		}
	}

	util.FormatTable("", append(lines, StringArraySlice(slice)...))
	return nil
}

func compare_column(c int) CompareFunction {
	return func(a interface{}, b interface{}) int {
		aa := a.([]string)
		ab := b.([]string)
		if len(aa) > c && len(ab) > c {
			return strings.Compare(aa[c], ab[c])
		}
		return len(aa) - len(ab)
	}

}
