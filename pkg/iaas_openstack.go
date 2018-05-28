package gube

import (
	"fmt"
)

func init() {
	registerIaaSHandler("openstack", &OpenstackHandler{})

}

type OpenstackHandler struct {
}

type OpenstackInfo struct {
	*_IaaSInfo
	authURL string
}

var _ IaaSInfo = &OpenstackInfo{}

func (this *OpenstackHandler) GetIaaSInfo(shoot Shoot) (IaaSInfo, error) {
	info := &OpenstackInfo{_IaaSInfo: NewStandardIaaSInfo(shoot)}

	//fmt.Printf("OS: %+v\n", info.GetInfraOutputs())

	p, err := shoot.GetProfile()
	if err != nil {
		return nil, err
	}
	info.authURL = p.GetManifest().Spec.OpenStack.KeyStoneURL
	return info, nil
}

func (this *OpenstackInfo) GetKeyInfo() string {
	id := this.GetNetworkId()
	if id != "" {
		return fmt.Sprintf("network_id: %s", id)
	}
	return this._IaaSInfo.GetKeyInfo()
}

func (this *OpenstackInfo) GetAuthURL() string {
	return this.authURL
}

func (this *OpenstackInfo) GetRouterId() string {
	return this.getInfraStringOutput("router_id")
}

func (this *OpenstackInfo) GetNetworkId() string {
	return this.getInfraStringOutput("network_id")
}

func (this *OpenstackInfo) GetKeyName() string {
	return this.getInfraStringOutput("key_name")
}

func (this *OpenstackInfo) GetFloatingNetworkId() string {
	return this.getInfraStringOutput("floating_network_id")
}

func (this *OpenstackInfo) GetSubnetId() string {
	return this.getInfraStringOutput("subnet_id")
}

func (this *OpenstackInfo) GetSecurityGroupName() string {
	return this.getInfraStringOutput("security_group_name")
}
