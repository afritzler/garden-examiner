package util

import (
	"fmt"

	"github.com/afritzler/garden-examiner/cmd/gex/context"
)

type TableOutput struct {
	data [][]string
}

var _ Output = &TableOutput{}

func (this *TableOutput) Add(ctx *context.Context, e interface{}) error {
	panic(fmt.Errorf("called abstract Add method"))
	return nil
}

func (this *TableOutput) Out(*context.Context) error {
	FormatTable(this.data)
	return nil
}

func (this *TableOutput) AddLine(line []string) *TableOutput {
	this.data = append(this.data, line)
	return this
}

func NewTableOutput(data [][]string) *TableOutput {
	return &TableOutput{data}
}

func FormatTable(data [][]string) {
	columns := []int{}

	for _, row := range data {
		for i, col := range row {
			if i >= len(columns) {
				columns = append(columns, len(col))
			} else {
				if columns[i] < len(col) {
					columns[i] = len(col)
				}
			}
		}
	}

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
