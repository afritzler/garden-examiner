package util

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Kubectl(config []byte, input []byte, args ...string) error {
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	defer r.Close()
	go func() {
		w.Write([]byte(config))
		w.Close()
	}()

	eff := append([]string{fmt.Sprintf("--kubeconfig=/dev/fd/%d", 3)}, args...)
	return ExecProcess(input, []*os.File{r}, "kubectl", eff...)
}

func ExecProcess(input []byte, extra []*os.File, c string, args ...string) error {
	var stdin = os.Stdin
	if input != nil {
		r, w, err := os.Pipe()
		if err != nil {
			return err
		}
		defer r.Close()
		go func() {
			w.Write([]byte(input))
			w.Close()
		}()
		stdin = r
	}
	cmd := exec.Command(c, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = stdin
	if extra != nil {
		cmd.ExtraFiles = extra
	}
	err := cmd.Run()
	if extra != nil {
		for _, f := range extra {
			f.Close()
		}
	}
	if err == nil && !cmd.ProcessState.Success() {
		return fmt.Errorf("command failed")
	}
	return err
}

func ExecCmd(cmd string, environment ...string) (err error) {
	var command *exec.Cmd
	parts := strings.Fields(cmd)
	head := parts[0]
	if len(parts) > 1 {
		parts = parts[1:len(parts)]
	} else {
		parts = nil
	}
	command = exec.Command(head, parts...)
	for index, env := range environment {
		if index == 0 {
			command.Env = append(os.Environ(),
				env,
			)
		} else {
			command.Env = append(command.Env,
				env,
			)
		}
	}
	val, err := command.Output()
	if err != nil {
		ee, ok := err.(*exec.ExitError)
		fmt.Println(string(ee.Stderr))
		if !ok {
			return err
		}
	}
	fmt.Println(string(val))
	return nil
}

func ExecCmdReturnOutput(cmd string, args ...string) (output string) {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return strings.TrimSpace(string(out[:]))
}