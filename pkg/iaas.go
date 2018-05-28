package gube

import (
	"fmt"
)

type IaaSInfo interface {
	GetKind() string
	GetRegion() string
	GetInfraOutputs() map[string]interface{}
	GetKeyInfo() string
}

type _IaaSInfo struct {
	kind            string
	region          string
	terraform_infra *TerraformState
}

var _ IaaSInfo = &_IaaSInfo{}

func NewStandardIaaSInfo(shoot Shoot) *_IaaSInfo {
	state := &TerraformState{}
	c, err := shoot.GetConfigMapEntriesFromSeed(shoot.GetName().GetName() + ".infra.tf-state")
	if err == nil {
		state, err = NewTerraformStateFromConfig(c)
		if err != nil {
			fmt.Printf("cannot unmarshal terraform state for shoot '%s': %s\n", shoot.GetName(), err)

		}
	} else {
		fmt.Printf("cannot get infrastructue terraform state for shoot '%s': %s\n", shoot.GetName(), err)
	}
	return &_IaaSInfo{shoot.GetInfrastructure(), shoot.GetRegion(), state}
}

func (this *_IaaSInfo) GetKind() string {
	return this.kind
}

func (this *_IaaSInfo) GetRegion() string {
	return this.region
}

func (this *_IaaSInfo) GetKeyInfo() string {
	return "unknown"
}

func (this *_IaaSInfo) GetInfraOutput(name string) interface{} {
	if this.terraform_infra == nil {
		return nil
	}
	return this.terraform_infra.GetOutput(name)
}

func (this *_IaaSInfo) getInfraStringOutput(name string) string {
	if this.terraform_infra == nil {
		return ""
	}
	o := this.terraform_infra.GetOutput(name)
	if o == nil {
		return ""
	}
	return o.(string)
}

func (this *_IaaSInfo) GetInfraOutputs() map[string]interface{} {
	if this.terraform_infra == nil {
		return nil
	}
	return this.terraform_infra.GetOutputs()
}

type IaaSHandler interface {
	GetIaaSInfo(Shoot) (IaaSInfo, error)
}

var iaas = map[string]IaaSHandler{}

func registerIaaSHandler(name string, handler IaaSHandler) {
	iaas[name] = handler
}
