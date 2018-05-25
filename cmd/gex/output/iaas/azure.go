package iaas

import (
	"fmt"

	"github.com/afritzler/garden-examiner/pkg"
)

func init() {
	RegisterIaasHandler(&azure{}, "azure")
}

type azure struct {
}

func (this *azure) Execute(shoot gube.Shoot, config map[string]string, args ...string) error {
	fmt.Println("Hello AZURE: %v", args)
	return nil
}
