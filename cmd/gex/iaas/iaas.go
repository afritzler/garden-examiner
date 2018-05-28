package iaas

import (
	"fmt"

	"github.com/afritzler/garden-examiner/pkg"
)

var IaasHandlers = map[string]IaasHandler{}

type IaasHandler interface {
	Execute(shoot gube.Shoot, config map[string]string, args ...string) error
	Describe(shoot gube.Shoot) error
}

func RegisterIaasHandler(h IaasHandler, name string) {
	IaasHandlers[name] = h
}

func Get(shoot gube.Shoot) IaasHandler {
	return IaasHandlers[shoot.GetInfrastructure()]
}

func Describe(shoot gube.Shoot) error {
	h := IaasHandlers[shoot.GetInfrastructure()]
	if h != nil {
		return h.Describe(shoot)
	}
	fmt.Printf("no handler for infrastructure '%s'\n", shoot.GetInfrastructure())
	return nil
}
