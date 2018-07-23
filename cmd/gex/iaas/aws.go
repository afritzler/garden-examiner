package iaas

import (
	"fmt"
	"strings"

	"github.com/afritzler/garden-examiner/cmd/gex/env"
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
	err := util.ExecCmd("aws "+strings.Join(args, " "), nil, "AWS_ACCESS_KEY_ID="+string(config["accessKeyID"]), "AWS_SECRET_ACCESS_KEY="+string(config["secretAccessKey"]), "AWS_DEFAULT_REGION="+region, "AWS_DEFAULT_OUTPUT=json")
	if err != nil {
		return fmt.Errorf("cannot execute 'aws': %s", err)
	}
	return nil
}

func (this *aws) Export(shoot gube.Shoot, config map[string]string, cachedir string) error {
	region := shoot.GetRegion()
	fmt.Printf("exporting AWS CLI config for %s\n", shoot.GetName())
	env.Set("AWS_ACCESS_KEY_ID", string(config["accessKeyID"]))
	env.Set("AWS_SECRET_ACCESS_KEY", string(config["secretAccessKey"]))
	env.Set("AWS_DEFAULT_REGION", region)
	env.Set("AWS_DEFAULT_OUTPUT", "json")
	return nil
}

func (this *aws) Describe(shoot gube.Shoot, attrs *util.AttributeSet) error {
	info, err := shoot.GetIaaSInfo()
	if err == nil {
		iaas := info.(*gube.AWSInfo)
		attrs.Attribute("AWS Information", "")
		attrs.Attribute("Region", iaas.GetRegion())
		attrs.Attribute("VPC Id", iaas.GetVpcId())
		attrs.Attribute("Security Group", iaas.GetNodesSecurityGroupId())
	}
	return nil
}
