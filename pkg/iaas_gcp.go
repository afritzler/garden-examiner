package gube

import (
	"fmt"
)

func init() {
	registerIaaSHandler("gcp", &GCPHandler{})
}

type GCPHandler struct {
}

type GCPInfo struct {
	*_IaaSInfo
}

var _ IaaSInfo = &GCPInfo{}

func (this *GCPHandler) GetIaaSInfo(shoot Shoot) (IaaSInfo, error) {
	info := &GCPInfo{_IaaSInfo: NewStandardIaaSInfo(shoot)}

	//fmt.Printf("GCP: %+v\n", info.GetInfraOutputs())
	return info, nil
}

func (this *GCPInfo) GetKeyInfo() string {
	id := this.GetServiceAccountEMail()
	if id != "" {
		return fmt.Sprintf("sa_user: %s", id)
	}
	return this._IaaSInfo.GetKeyInfo()
}

func (this *GCPInfo) GetVpcName() string {
	return this.getInfraStringOutput("vpc_name")
}

func (this *GCPInfo) GetServiceAccountEMail() string {
	return this.getInfraStringOutput("service_account_email")
}
