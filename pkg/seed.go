package gube

import (
	"fmt"

	v1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"github.com/gardener/gardener/pkg/operation/common"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

const secretkubeconfig = "kubeconfig"

type Seed interface {
	GetName() string
	GetManifest() *v1beta1.Seed
	GetCloud() v1beta1.SeedCloud
	GetClientset() (*kubernetes.Clientset, error)
	GetShoot() *ShootName
	GetInfrastructure() string
	KubeconfigProvider
	RuntimeObjectWrapper
}

type seed struct {
	garden   Garden
	name     string
	manifest v1beta1.Seed
}

func NewSeedFromSeedManifest(g Garden, m v1beta1.Seed) Seed {
	return &seed{garden: g, name: m.GetName(), manifest: m}
}

func (s *seed) GetName() string {
	return s.name
}

func (s *seed) GetManifest() *v1beta1.Seed {
	return &s.manifest
}

func (s *seed) GetCloud() v1beta1.SeedCloud {
	return s.manifest.Spec.Cloud
}

func (s *seed) GetRuntimeObject() runtime.Object {
	return &s.manifest
}

func (s *seed) GetKubeconfig() ([]byte, error) {
	secret, err := s.garden.GetSecretByRef(s.manifest.Spec.SecretRef)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret for seed %s: %s", s.name, err)
	}
	return secret.Data[secretkubeconfig], nil
}

func (s *seed) GetClientConfig() (*restclient.Config, error) {
	bytes, err := s.GetKubeconfig()
	if err != nil {
		return nil, err
	}
	return NewConfigFromBytes(bytes)
}

func (s *seed) GetClientset() (*kubernetes.Clientset, error) {
	bytes, err := s.GetKubeconfig()
	if err != nil {
		return nil, err
	}
	return NewClientFromBytes(bytes)
}

func (s *seed) GetInfrastructure() string {
	p, err := s.garden.GetProfile(s.manifest.Spec.Cloud.Profile)
	if err == nil {
		return p.GetInfrastructure()
	}
	return "unknown"
}

func (s *seed) GetShoot() *ShootName {
	for _, o := range s.manifest.ObjectMeta.GetOwnerReferences() {
		if o.APIVersion == v1beta1.SchemeGroupVersion.String() {
			if o.Kind == "Shoot" {
				return NewShootName(common.GardenNamespace, o.Name)
			}
		}
	}
	return nil
}
