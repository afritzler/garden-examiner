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
