package output

import (
	"github.com/afritzler/garden-examiner/pkg"

	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
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
	return util.Kubectl(this.kubecfg, input, args...)
}

func (this *KubectlOutput) GetArgs() []string {
	return this.args
}
