package env

import (
	"fmt"
	"os"
	"strings"
)

var env *os.File = nil

func init() {
	//env = os.NewFile(3, "env-setting")  // does not work anymore
	env, _ = os.OpenFile("/dev/fd/3", os.O_WRONLY, 0755)
}

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
	_, err := fmt.Fprintf(env, format, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot write script\n")
	}
	return nil
}

func Set(name, value string) error {
	return Printf("export %s='%s'\n", name, strings.Replace(value, "'", "\\'", -1))
}

func UnSet(name string) error {
	return Printf("unset %s\n", name)
}
