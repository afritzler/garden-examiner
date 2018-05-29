package iaas

import (
	"fmt"
	"strings"

	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/pkg"
)

func init() {
	RegisterIaasHandler(&openstack{}, "openstack")
}

type openstack struct {
}

func (this *openstack) Execute(shoot gube.Shoot, config map[string]string, args ...string) error {
	info, err := shoot.GetIaaSInfo()
	if err != nil {
		return err
	}
	authURL := info.(*gube.OpenstackInfo).GetAuthURL()
	if authURL == "" {
		fmt.Println("Fetching authURL was not successful")
		return nil
	}
	err = util.ExecCmd("openstack "+strings.Join(args, " "), "OS_IDENTITY_API_VERSION=3", "OS_AUTH_VERSION=3", "OS_AUTH_STRATEGY=keystone", "OS_AUTH_URL="+authURL, "OS_TENANT_NAME="+string(config["tenantName"]),
		"OS_PROJECT_DOMAIN_NAME="+string(config["domainName"]), "OS_USER_DOMAIN_NAME="+string(config["domainName"]), "OS_USERNAME="+string(config["username"]), "OS_PASSWORD="+string(config["password"]), "OS_REGION_NAME="+string(config["region"]))
	if err != nil {
		return fmt.Errorf("cannot execute 'openstack': %s", err)
	}
	return nil
}

func (this *openstack) Describe(shoot gube.Shoot) error {
	info, err := shoot.GetIaaSInfo()
	if err == nil {
		iaas := info.(*gube.OpenstackInfo)
		fmt.Printf("Openstack Information:\n")
		fmt.Printf("Keystone URL: %s\n", iaas.GetAuthURL())
		fmt.Printf("Region: %s\n", iaas.GetRegion())
		fmt.Printf("Router Id: %s\n", iaas.GetRouterId())
		fmt.Printf("Network Id: %s\n", iaas.GetNetworkId())
		fmt.Printf("Subnet Id: %s\n", iaas.GetSubnetId())
		fmt.Printf("Security Group: %s\n", iaas.GetSecurityGroupName())
	}
	return nil
}
