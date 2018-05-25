package util

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/afritzler/garden-examiner/pkg"

	"github.com/afritzler/garden-examiner/cmd/gex/context"
)

type ElementMapper func(*context.Context, interface{}) (interface{}, []string, error)

type KubectlOutput struct {
	*SingleElementOutput
	kubecfg []byte
	mapper  func(*context.Context, interface{}) (interface{}, []string, error)
	args    []string
}

var _ Output = &KubectlOutput{}

func NewKubectlOutput(args []string, mapper ElementMapper) *KubectlOutput {
	return &KubectlOutput{NewSingleElementOutput(), nil, mapper, args}
}

func (this *KubectlOutput) Close(ctx *context.Context) error {
	return nil
}

func (this *KubectlOutput) Add(ctx *context.Context, e interface{}) error {
	if this.mapper != nil {
		m, args, err := this.mapper(ctx, e)
		if err != nil {
			return err
		}
		if args != nil {
			this.args = append(this.args, args...)
		}
		e = m
	}
	s := e.(gube.KubeconfigProvider)
	cfg, err := s.GetKubeconfig()
	if err != nil {
		return err
	}
	this.kubecfg = cfg
	return this.SingleElementOutput.Add(ctx, e)
}

func (this *KubectlOutput) Out(ctx *context.Context) error {
	return this.Kubectl(nil, this.GetArgs()...)
}

func (this *KubectlOutput) Kubectl(input []byte, args ...string) error {
	return Kubectl(this.kubecfg, input, args...)
}

func (this *KubectlOutput) GetArgs() []string {
	return this.args
}

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
