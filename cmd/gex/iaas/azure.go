package iaas

import (
	"fmt"
	"strings"

	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/pkg"
)

func init() {
	RegisterIaasHandler(&azure{}, "azure")
}

type azure struct {
}

func (this *azure) Execute(shoot gube.Shoot, config map[string]string, args ...string) error {
	err := util.ExecCmd("az login --service-principal -u " + string(config["clientID"]) + " -p " + string(config["clientSecret"]) + " --tenant " + string(config["tenantID"]))
	if err != nil {
		return fmt.Errorf("cannot login: %s", err)
	}
	err = util.ExecCmd("az " + strings.Join(args, " "))
	if err != nil {
		return fmt.Errorf("cannot execute 'az': %s", err)
	}
	return nil
}

func (this *azure) Export(shoot gube.Shoot, config map[string]string, cachedir string) error {
	return fmt.Errorf("no possible for Azure")
}

func (this *azure) Describe(shoot gube.Shoot) error {
	info, err := shoot.GetIaaSInfo()
	if err == nil {
		iaas := info.(*gube.AzureInfo)
		attrs := util.NewAttributeSet()
		fmt.Printf("Azure Information:\n")
		attrs.Attribute("Region", iaas.GetRegion())
		attrs.Attribute("Resource Group", iaas.GetResourceGroupName())
		attrs.Attribute("VNet", iaas.GetVNetName())
		attrs.Attribute("Subnet", iaas.GetSubnetName())
		attrs.PrintAttributes()
	}
	return nil
}
