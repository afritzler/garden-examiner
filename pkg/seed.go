package gube

import (
	"fmt"

	v1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"github.com/gardener/gardener/pkg/operation/common"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const secretkubeconfig = "kubeconfig"

type Seed interface {
	GetName() string
	GetManifest() *v1beta1.Seed
	GetCloud() v1beta1.SeedCloud
	GetRegion() string
	GetProfileName() string
	GetProfile() (Profile, error)
	GetInfrastructure() string
	Cluster
	RuntimeObjectWrapper
	GardenObject
}

type seed struct {
	_GardenObject
	cluster
	name     string
	manifest v1beta1.Seed
}

func NewSeedFromSeedManifest(g Garden, m v1beta1.Seed) Seed {
	return (&seed{}).new(g, m)
}

func (s *seed) new(g Garden, m v1beta1.Seed) Seed {
	m.Kind = "Seed"
	m.APIVersion = v1beta1.SchemeGroupVersion.String()
	s._GardenObject.new(g)
	s.cluster.new("seed "+m.GetName(), s)
	s.name = m.GetName()
	s.manifest = m
	return s
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

func (s *seed) GetRegion() string {
	return s.manifest.Spec.Cloud.Region
}

func (s *seed) GetProfileName() string {
	return s.manifest.Spec.Cloud.Profile
}

func (s *seed) GetProfile() (Profile, error) {
	return s.garden.GetProfile(s.GetProfileName())
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

func (s *seed) GetSecretByRef(secretref corev1.SecretReference) (*corev1.Secret, error) {
	kubeset, err := s.GetClientset()
	if err != nil {
		return nil, fmt.Errorf("failed to get secret for seed %s: %s", s.name, err)
	}
	secret, err := kubeset.CoreV1().Secrets(secretref.Namespace).Get(secretref.Name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %s for namespace %s: %s", secretref.Name, secretref.Namespace, err)
	}
	return secret, nil
}

func (s *seed) GetInfrastructure() string {
	p, err := s.garden.GetProfile(s.manifest.Spec.Cloud.Profile)
	if err == nil {
		return p.GetInfrastructure()
	}
	return "unknown"
}

func (s *seed) GetShootName() *ShootName {
	for _, o := range s.manifest.ObjectMeta.GetOwnerReferences() {
		if o.APIVersion == v1beta1.SchemeGroupVersion.String() {
			if o.Kind == "Shoot" {
				return NewShootName(common.GardenNamespace, o.Name)
			}
		}
	}
	return nil
}

func (s *seed) AsShoot() (Shoot, error) {
	n := s.GetShootName()
	if n == nil {
		return s.cluster.AsShoot()
	}
	return s.garden.GetShoot(n)
}
