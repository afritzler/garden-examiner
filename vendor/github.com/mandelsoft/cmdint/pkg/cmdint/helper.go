package cmdint

import (
	"fmt"
	"strings"
)

func DecodeDescription(e interface{}) []string {
	var d string
	switch t := e.(type) {
	case Command:
		d = t.GetCmdDescription()
	case *Option:
		d = t.GetDescription()
	default:
		panic(fmt.Errorf("invalid type: %T", e))
	}
	return strings.Split(d, "\n")
}

func compact(desc []string) string {
	return strings.Join(desc, "\n")
}
