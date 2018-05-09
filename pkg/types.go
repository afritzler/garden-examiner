package gube

import (
	"fmt"

	v1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type RuntimeObjectWrapper interface {
	GetRuntimeObject() runtime.Object
}

type ShootName struct {
	name        string
	projectname string
}

func (n *ShootName) GetName() string {
	return n.name
}

func (n *ShootName) GetProjectName() string {
	return n.projectname
}

func (n *ShootName) String() string {
	return fmt.Sprintf("%s/%s", n.projectname, n.name)
}

func NewShootName(project, name string) *ShootName {
	return &ShootName{name, project}
}

func NewShootNameFromShootManifest(garden Garden, shoot v1beta1.Shoot) (*ShootName, error) {
	p, err := garden.GetProjectByNamespace(shoot.GetNamespace())
	if err != nil {
		return nil, err
	}
	return &ShootName{name: shoot.GetName(), projectname: p.GetName()}, nil
}

func NewConfigFromBytes(kubeconfig []byte) (*restclient.Config, error) {
	configObj, err := clientcmd.Load(kubeconfig)
	if err != nil {
		return nil, err
	}
	clientConfig := clientcmd.NewDefaultClientConfig(*configObj, &clientcmd.ConfigOverrides{})
	return clientConfig.ClientConfig()
}

func NewClientFromBytes(kubeconfig []byte) (*kubernetes.Clientset, error) {
	config, err := NewConfigFromBytes(kubeconfig)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
