package iaas

import (
	"fmt"
	"strings"

	"github.com/afritzler/garden-examiner/cmd/gex/env"
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
	err = util.ExecCmd("openstack "+strings.Join(args, " "), nil, "OS_IDENTITY_API_VERSION=3", "OS_AUTH_VERSION=3", "OS_AUTH_STRATEGY=keystone", "OS_AUTH_URL="+authURL, "OS_TENANT_NAME="+string(config["tenantName"]),
		"OS_PROJECT_DOMAIN_NAME="+string(config["domainName"]), "OS_USER_DOMAIN_NAME="+string(config["domainName"]), "OS_USERNAME="+string(config["username"]), "OS_PASSWORD="+string(config["password"]), "OS_REGION_NAME="+string(config["region"]))
	if err != nil {
		return fmt.Errorf("cannot execute 'openstack': %s", err)
	}
	return nil
}

func (this *openstack) Export(shoot gube.Shoot, config map[string]string, cachedir string) error {
	info, err := shoot.GetIaaSInfo()
	if err != nil {
		return err
	}
	authURL := info.(*gube.OpenstackInfo).GetAuthURL()
	if authURL == "" {
		fmt.Println("Fetching authURL was not successful")
		return nil
	}
	fmt.Printf("exporting Openstack CLI config for %s\n", shoot.GetName())
	env.Set("OS_IDENTITY_API_VERSION", "3")
	env.Set("OS_AUTH_VERSION", "3")
	env.Set("OS_AUTH_STRATEGY", "keystone")
	env.Set("OS_AUTH_URL", authURL)
	env.Set("OS_TENANT_NAME", string(config["tenantName"]))
	env.Set("OS_PROJECT_DOMAIN_NAME", string(config["domainName"]))
	env.Set("OS_USER_DOMAIN_NAME", string(config["domainName"]))
	env.Set("OS_USERNAME", string(config["username"]))
	env.Set("OS_PASSWORD,", string(config["password"]))
	env.Set("OS_REGION_NAME", string(config["region"]))
	return nil
}

func (this *openstack) Describe(shoot gube.Shoot, attrs *util.AttributeSet) error {
	info, err := shoot.GetIaaSInfo()
	if err == nil {
		iaas := info.(*gube.OpenstackInfo)
		attrs.Attribute("Openstack Information", "")
		attrs.Attribute("Keystone URL", iaas.GetAuthURL())
		attrs.Attribute("Domain Name", iaas.GetDomainName())
		attrs.Attribute("Tenant Name", iaas.GetTenantName())
		attrs.Attribute("Username", iaas.GetUserName())
		attrs.Attribute("Password", iaas.GetPassword())
		attrs.Attribute("Region", iaas.GetRegion())
		attrs.Attribute("Router Id", iaas.GetRouterId())
		attrs.Attribute("Network Id", iaas.GetNetworkId())
		attrs.Attribute("Subnet Id", iaas.GetSubnetId())
		attrs.Attribute("Security Group", iaas.GetSecurityGroupName())
	}
	return nil
}
