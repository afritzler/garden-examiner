package iaas

import (
	"fmt"
	"strings"

	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/pkg"
)

func init() {
	RegisterIaasHandler(&aws{}, "aws")
}

type aws struct {
}

func (this *aws) Execute(shoot gube.Shoot, config map[string]string, args ...string) error {
	region := shoot.GetRegion()
	err := util.ExecCmd("aws "+strings.Join(args, " "), "AWS_ACCESS_KEY_ID="+string(config["accessKeyID"]), "AWS_SECRET_ACCESS_KEY="+string(config["secretAccessKey"]), "AWS_DEFAULT_REGION="+region, "AWS_DEFAULT_OUTPUT=text")
	if err != nil {
		return fmt.Errorf("cannot execute 'aws': %s", err)
	}
	return nil
}

func (this *aws) Describe(shoot gube.Shoot) error {
	info, err := shoot.GetIaaSInfo()
	if err == nil {
		iaas := info.(*gube.AWSInfo)
		fmt.Printf("AWS Information:\n")
		fmt.Printf("VPC Id: %s\n", iaas.GetVpcId())
		fmt.Printf("Security Group: %s\n", iaas.GetNodesSecurityGroupId())
	}
	return nil
}
