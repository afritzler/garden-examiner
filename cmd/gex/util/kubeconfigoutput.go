package util

import (
	"fmt"
	"strings"

	"github.com/afritzler/garden-examiner/pkg"

	"github.com/afritzler/garden-examiner/cmd/gex/context"
)

type kubeconfig_output struct {
	data []string
}

func NewKubeconfigOutput() Output {
	return &kubeconfig_output{[]string{}}
}

func (this *kubeconfig_output) Add(ctx *context.Context, e interface{}) error {
	s := e.(gube.KubeconfigProvider)
	cfg, err := s.GetKubeconfig()
	if err != nil {
		return err
	}
	this.data = append(this.data, string(cfg))
	return nil
}

func (this *kubeconfig_output) Out(ctx *context.Context) {
	for _, cfg := range this.data {
		if !strings.HasPrefix(cfg, "---\n") {
			fmt.Println("---")
		}
		fmt.Println(cfg)
	}
}
