package iaas

import (
	"fmt"

	"github.com/afritzler/garden-examiner/pkg"
)

func init() {
	RegisterIaasHandler(&gcp{}, "gcp")
}

type gcp struct {
}

func (this *gcp) Execute(shoot gube.Shoot, config map[string]string, args ...string) error {
	fmt.Println("Hello GCP: %v", args)
	return nil
}
