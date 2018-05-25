package iaas

import (
	"fmt"

	"github.com/afritzler/garden-examiner/pkg"
)

func init() {
	RegisterIaasHandler(&aws{}, "aws")
}

type aws struct {
}

func (this *aws) Execute(shoot gube.Shoot, config map[string]string, args ...string) error {
	fmt.Println("Hello AWS: \n AWS_ACCESS_KEY_ID="+string(config["accessKeyID"])+"AWS_SECRET_ACCESS_KEY="+string(config["secretAccessKey"])+"AWS_DEFAULT_OUTPUT=text %v", args)
	// region = shoot.Spec.Cloud.Region
	return nil
}
