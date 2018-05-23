package util

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/pkg"
	"github.com/afritzler/garden-examiner/pkg/data"

	"github.com/afritzler/garden-examiner/cmd/gex/context"
)

type kubeconfig_output struct {
	ElementOutput
}

var _Output = kubeconfig_output{}

func KubeconfigOutputFactory(opts *cmdint.Options) Output {
	return NewKubeconfigOutput()
}

func NewKubeconfigOutput() Output {
	return (&kubeconfig_output{}).new()
}

func (this *kubeconfig_output) new() Output {
	this.ElementOutput.new(data.Chain().Parallel(20).Map(map_kubeconfig_output))
	return this
}

func map_kubeconfig_output(e interface{}) interface{} {
	s := e.(gube.KubeconfigProvider)
	cfg, err := s.GetKubeconfig()
	if err != nil {
		return err
	}
	return string(cfg)
}

func (this *kubeconfig_output) Out(ctx *context.Context) error {
	i := this.Elems.Iterator()
	for i.HasNext() {
		switch cfg := i.Next().(type) {
		case error:
			return cfg
		case string:
			if !strings.HasPrefix(cfg, "---\n") {
				fmt.Println("---")
			}
			fmt.Println(cfg)
		}
	}
	return nil
}
