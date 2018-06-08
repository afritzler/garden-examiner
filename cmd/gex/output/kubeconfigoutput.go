package output

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/pkg"
)

func KubeconfigOutputFactory(opts *cmdint.Options) Output {
	return NewKubeconfigOutput()
}

func NewKubeconfigOutput() Output {
	return NewStringOutput(map_kubeconfig_output)
}

func map_kubeconfig_output(e interface{}) interface{} {
	s := e.(gube.KubeconfigProvider)
	cfg, err := s.GetKubeconfig()
	if err != nil {
		return err
	}
	return string(cfg)
}
