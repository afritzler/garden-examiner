package util

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/afritzler/garden-examiner/pkg"

	"github.com/afritzler/garden-examiner/cmd/gex/context"
)

type kubectl_output struct {
	*SingleElementOutput
	mapper func(*context.Context, interface{}) (interface{}, []string, error)
	args   []string
}

var _ Output = &kubectl_output{}

func NewKubectlOutput(args []string, mapper func(*context.Context, interface{}) (interface{}, []string, error)) Output {
	return &kubectl_output{NewSingleElementOutput(), mapper, args}
}

func (this *kubectl_output) Add(ctx *context.Context, e interface{}) error {
	if this.mapper != nil {
		m, args, err := this.mapper(ctx, e)
		if err != nil {
			return err
		}
		this.args = append(this.args, args...)
		e = m
	}
	s := e.(gube.KubeconfigProvider)
	cfg, err := s.GetKubeconfig()
	if err != nil {
		return err
	}
	return this.SingleElementOutput.Add(ctx, cfg)
}

func (this *kubectl_output) Out(ctx *context.Context) error {
	return Kubectl(this.Elem.([]byte), this.args...)
}

func Kubectl(config []byte, args ...string) error {
	r, w, err := os.Pipe()
	defer r.Close()
	go func() {
		w.Write([]byte(config))
		w.Close()
	}()

	eff := append([]string{fmt.Sprintf("--kubeconfig=/dev/fd/%d", 3)}, args...)
	fmt.Printf("ARGS: %v\n", args)

	cmd := exec.Command("kubectl", eff...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = []*os.File{r}
	err = cmd.Run()
	r.Close()
	if err == nil && !cmd.ProcessState.Success() {
		return fmt.Errorf("command failed")
	}
	return err
}
