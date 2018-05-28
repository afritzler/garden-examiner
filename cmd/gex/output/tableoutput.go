package output

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
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

	FormatTable(append(lines, StringArraySlice(slice)...))
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

///////////////////////////////////////////////////////////////////////////

func FormatTable(data [][]string) {
	columns := []int{}
	max := 0

	for _, row := range data {
		for i, col := range row {
			if i >= len(columns) {
				columns = append(columns, len(col))
			} else {
				if columns[i] < len(col) {
					columns[i] = len(col)
				}
			}
			if len(col) > max {
				max = len(col)
			}
		}
	}

	if len(columns) <= 3 && max > 100 {
		first := []string{}
		for i, row := range data {
			if i == 0 {
				first = row
			} else {
				for c, col := range row {
					fmt.Printf("%s: %s\n", first[c], col)
				}
				fmt.Printf("---\n")
			}
		}
	} else {
		format := ""
		for _, col := range columns {
			format = fmt.Sprintf("%s%%-%ds ", format, col)
		}
		format = format[:len(format)-1] + "\n"
		for _, row := range data {
			r := []interface{}{}
			for i := 0; i < len(columns); i++ {
				if i < len(row) {
					r = append(r, row[i])
				} else {
					r = append(r, "")
				}
			}
			fmt.Printf(format, r...)
		}
	}
}
