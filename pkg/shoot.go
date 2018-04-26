package gube

import (
	"fmt"

	v1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

type Shoot interface {
	GetName() *ShootName
	GetManifest() *v1beta1.Shoot
	GetSeedName() string
	GetSeed() (Seed, error)
	GetProject() (Project, error)
	GetKubeconfig() ([]byte, error)
	GetSecretRef() (*corev1.SecretReference, error)
	GetClientset() (*kubernetes.Clientset, error)
	GetNodeCount() (int, error)
}

type shoot struct {
	garden    Garden
	name      *ShootName
	namespace string
	manifest  v1beta1.Shoot
	clientset *kubernetes.Clientset
}

var _ Shoot = &shoot{}

func NewShootFromShootManifest(g Garden, m v1beta1.Shoot) (Shoot, error) {
	n, err := NewShootNameFromShootManifest(g, m)
	if err != nil {
		return nil, err
	}
	return &shoot{garden: g, name: n, manifest: m, namespace: m.GetObjectMeta().GetNamespace()}, nil
}

func (s *shoot) GetName() *ShootName {
	return s.name
}

func (s *shoot) GetNamespace() (string, error) {
	if s.namespace == "" {
		p, err := s.GetProject()
		if err != nil {
			return "", fmt.Errorf("cannot get namespace for shoot '%s': %s", s.name, err)
		}
		s.namespace = p.GetNamespace()
	}
	return s.namespace, nil
}

func (s *shoot) GetManifest() *v1beta1.Shoot {
	return &s.manifest
}

func (s *shoot) GetSeedName() string {
	return *s.manifest.Spec.Cloud.Seed
}

func (s *shoot) GetSeed() (Seed, error) {
	// should never fail with panic :-P
	return s.garden.GetSeed(s.GetSeedName())
}

func (s *shoot) GetProject() (Project, error) {
	return s.garden.GetProject(s.name.GetProjectName())
}

func (s *shoot) GetKubeconfig() ([]byte, error) {
	ref, err := s.GetSecretRef()
	if err != nil {
		return nil, fmt.Errorf("could not get secret ref for shoot '%s': %s", s.name, err)
	}
	secret, err := s.garden.GetSecretByRef(*ref)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret for shoot '%s': %s", s.name, err)
	}
	return secret.Data[secretkubeconfig], nil
}

func (s *shoot) GetClientConfig() (*restclient.Config, error) {
	bytes, err := s.GetKubeconfig()
	if err != nil {
		return nil, err
	}
	return NewConfigFromBytes(bytes)
}

func (s *shoot) GetClientset() (*kubernetes.Clientset, error) {
	bytes, err := s.GetKubeconfig()
	if err != nil {
		return nil, err
	}
	return NewClientFromBytes(bytes)
}

func (s *shoot) GetSecretRef() (*corev1.SecretReference, error) {
	ns, err := s.GetNamespace()
	if err != nil {
		return nil, err
	}
	return &corev1.SecretReference{Name: fmt.Sprintf("%s.kubeconfig", s.name.GetName()), Namespace: ns}, nil
}

func (s *shoot) GetNodeCount() (int, error) {
	cs, err := s.GetClientset()
	if err != nil {
		return 0, err
	}
	list, err := cs.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to get node count for shoot %s: %s", s.name, err)
	}
	return len(list.Items), nil
}
