package util

import (
	"fmt"
	"strings"
)

func FormatTable(gap string, data [][]string) {
	columns := []int{}
	max := 0

	formats := []string{}
	if len(data) > 1 {
		for i, f := range data[0] {
			if strings.HasPrefix(f, "-") {
				formats = append(formats, "")
				data[0][i] = f[1:]
			} else {
				formats = append(formats, "-")
			}
		}
	}

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

	if len(columns) > 2 && max > 100 {
		first := []string{}
		setSep := false
		for i, row := range data {
			if i == 0 {
				first = row
			} else {
				for c, col := range row {
					if c < len(first) {
						fmt.Printf("%s%s: %s\n", gap, first[c], col)
					} else {
						fmt.Printf("%s%d: %s\n", gap, c, col)
					}
					setSep = true
				}
				if setSep == true {
					fmt.Printf("---\n")
					setSep = false
				}
			}
		}
	} else {
		format := gap
		for i, col := range columns {
			f := "-"
			if i < len(formats) {
				f = formats[i]
			}
			if i == len(columns)-1 && f == "-" {
				format = fmt.Sprintf("%s%%s ", format)
			} else {
				format = fmt.Sprintf("%s%%%s%ds ", format, f, col)
			}
		}
		format = format[:len(format)-1] + "\n"
		for _, row := range data {
			if len(row) > 0 {
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
}
