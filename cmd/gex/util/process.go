package util

import (
	"fmt"
	"os"
	"os/exec"
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
