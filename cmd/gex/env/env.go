package env

import (
	"fmt"
	"os"
	"strings"
)

var env = os.NewFile(3, "env-setting")
var err = fmt.Errorf("Environment not modifiable, please use the command alias")
var done = false

func IsModifiable() bool {
	return env != nil
}

func Warning() bool {
	if !IsModifiable() {
		if !done {
			fmt.Printf("%s\n", err)
			done = true
		}
		return true
	}
	return false
}

func Printf(format string, args ...interface{}) error {
	if !IsModifiable() {
		return err
	}
	fmt.Fprintf(env, format, args...)
	return nil
}

func Set(name, value string) error {
	return Printf("export %s='%s'\n", name, strings.Replace(value, "'", "\\'", -1))
}

func UnSet(name string) error {
	return Printf("unset %s\n", name)
}
