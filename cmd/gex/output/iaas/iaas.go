package iaas

import (
	"github.com/afritzler/garden-examiner/pkg"
)

var IaasHandlers = map[string]IaasHandler{}

type IaasHandler interface {
	Execute(shoot gube.Shoot, config map[string]string, args ...string) error
}

func RegisterIaasHandler(h IaasHandler, name string) {
	IaasHandlers[name] = h
}
