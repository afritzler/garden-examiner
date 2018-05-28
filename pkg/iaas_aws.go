package gube

import (
	"fmt"
)

func init() {
	registerIaaSHandler("aws", &AWSHandler{})

}

type AWSHandler struct {
}

type AWSInfo struct {
	*_IaaSInfo
}

var _ IaaSInfo = &AWSInfo{}

func (this *AWSHandler) GetIaaSInfo(shoot Shoot) (IaaSInfo, error) {
	info := &AWSInfo{_IaaSInfo: NewStandardIaaSInfo(shoot)}

	//fmt.Printf("AWS: %+v\n", info.GetInfraOutputs())
	return info, nil
}

func (this *AWSInfo) GetKeyInfo() string {
	id := this.GetVpcId()
	if id != "" {
		return fmt.Sprintf("vpc_id: %s", id)
	}
	return this._IaaSInfo.GetKeyInfo()
}

func (this *AWSInfo) GetKeyName() string {
	return this.getInfraStringOutput("keyName")
}

func (this *AWSInfo) GetVpcId() string {
	return this.getInfraStringOutput("vpc_id")
}

func (this *AWSInfo) GetNodesSecurityGroupId() string {
	return this.getInfraStringOutput("security_group_nodes")
}
