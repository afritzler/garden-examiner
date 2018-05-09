package gube

import (
	"fmt"

	v1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
	GetInfrastructure() string
	GetNodeCount() (int, error)
	GetState() string
	GetError() string
	RuntimeObjectWrapper
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

func (s *shoot) GetRuntimeObject() runtime.Object {
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

func (s *shoot) GetState() string {
	return string(s.manifest.Status.LastOperation.State)
}

func (s *shoot) GetError() string {
	if s.manifest.Status.LastOperation.State != v1beta1.ShootLastOperationStateSucceeded {
		if s.manifest.Status.LastError != nil {
			return s.manifest.Status.LastError.Description
		}
	}
	return ""
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

func (s *shoot) GetInfrastructure() string {
	if s.manifest.Spec.Cloud.AWS != nil {
		return "aws"
	}
	if s.manifest.Spec.Cloud.Azure != nil {
		return "azure"
	}
	if s.manifest.Spec.Cloud.OpenStack != nil {
		return "openstack"
	}
	if s.manifest.Spec.Cloud.GCP != nil {
		return "gcp"
	}
	if s.manifest.Spec.Cloud.Local != nil {
		return "local"
	}
	return "unknown"
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

//////////////////////////////////////////////////////////////////////////////
// cache

type ShootCache interface {
	GetShoots() (map[ShootName]Shoot, error)
	GetShoot(name *ShootName) (Shoot, error)
	Reset()
}

type shoot_cache struct {
	garden   Garden
	shoots   map[ShootName]Shoot
	complete bool
}

func NewShootCache(g Garden) ShootCache {
	return &shoot_cache{g, nil, false}
}

func (this *shoot_cache) Reset() {
	this.shoots = nil
	this.complete = false
}

func (this *shoot_cache) GetShoots() (map[ShootName]Shoot, error) {
	if this.shoots == nil || !this.complete {
		elems, err := this.garden.GetShoots()
		if err != nil {
			return nil, err
		}
		this.shoots = elems
		this.complete = true
	}
	return this.shoots, nil
}

func (this *shoot_cache) GetShoot(name *ShootName) (Shoot, error) {
	var p Shoot = nil
	if this.shoots != nil {
		p = this.shoots[*name]
	}
	if p == nil && !this.complete {
		elem, err := this.garden.GetShoot(name)
		if err != nil {
			return nil, err
		}
		if this.shoots == nil {
			this.shoots = map[ShootName]Shoot{}
		}
		this.shoots[*name] = elem
		p = elem
	}
	if p == nil {
		return nil, fmt.Errorf("shoot '%s' not found", name)
	}
	return p, nil
}
