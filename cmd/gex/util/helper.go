package util

import (
	"strings"
	"unicode/utf8"
)

func Oneline(str string, max int) string {
	str = strings.Replace(str, "\n", "\\n", -1)
	l := utf8.RuneCountInString(str)

	if max > 0 && l > max {
		limit := max
		if max > 6 {
			limit = max - 3
		}
		limited := ""
		for i, r := range str {
			if i > limit {
				break
			}
			limited += string(r)
		}
		if limit != max {
			return limited + "..."
		}
		return limited
	}

	return str
}
