package iaas

import (
	"fmt"

	"github.com/afritzler/garden-examiner/pkg"
)

func init() {
	RegisterIaasHandler(&openstack{}, "openstack")
}

type openstack struct {
}

func (this *openstack) Execute(shoot gube.Shoot, config map[string]string, args ...string) error {
	fmt.Println("Hello OPENSTACK: %v", args)
	return nil
}
