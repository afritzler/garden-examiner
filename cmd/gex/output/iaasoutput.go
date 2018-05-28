package output

import (
	"fmt"

	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/output/iaas"
	"github.com/afritzler/garden-examiner/pkg"
)

type IaasOutput struct {
	*SingleElementOutput
	mapper ElementMapper
	args   []string
}

var _ Output = &IaasOutput{}

func NewIaasOutput(args []string, mapper ElementMapper) *IaasOutput {
	return &IaasOutput{NewSingleElementOutput(), mapper, args}
}

func (this *IaasOutput) Close(ctx *context.Context) error {
	return nil
}

func (this *IaasOutput) Add(ctx *context.Context, e interface{}) error {
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
	return this.SingleElementOutput.Add(ctx, e)
}

func (this *IaasOutput) Out(ctx *context.Context) error {
	return this.Iaas(nil, this.GetArgs()...)
}

func (this *IaasOutput) Iaas(input []byte, args ...string) error {
	shoot, err := this.Elem.(gube.Shooted).AsShoot()
	if err != nil {
		return err
	}
	iaasType := shoot.GetInfrastructure()
	h, ok := iaas.IaasHandlers[iaasType]
	if !ok {
		return fmt.Errorf("No handler for infrastructure '%s'", iaasType)
	}
	config, err := shoot.GetCloudProviderConfig()
	if err != nil {
		return err
	}
	return h.Execute(shoot, config, args...)
}

func (this *IaasOutput) GetArgs() []string {
	return this.args
}
