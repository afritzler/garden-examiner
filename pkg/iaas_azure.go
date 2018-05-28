package gube

import (
	"fmt"
)

func init() {
	registerIaaSHandler("azure", &AzureHandler{})
}

type AzureHandler struct {
}

type AzureInfo struct {
	*_IaaSInfo
}

var _ IaaSInfo = &AzureInfo{}

func (this *AzureHandler) GetIaaSInfo(shoot Shoot) (IaaSInfo, error) {
	info := &AzureInfo{_IaaSInfo: NewStandardIaaSInfo(shoot)}

	//fmt.Printf("Azure: %+v\n", info.GetInfraOutputs())
	return info, nil
}

func (this *AzureInfo) GetKeyInfo() string {
	id := this.GetResourceGroupName()
	if id != "" {
		return fmt.Sprintf("rsc_group: %s", id)
	}
	return this._IaaSInfo.GetKeyInfo()
}

func (this *AzureInfo) GetVNetName() string {
	return this.getInfraStringOutput("vnetName")
}

func (this *AzureInfo) GetSubnetName() string {
	return this.getInfraStringOutput("subnetName")
}

func (this *AzureInfo) GetResourceGroupName() string {
	return this.getInfraStringOutput("resourceGroupName")
}
