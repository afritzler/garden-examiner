package gube

import (
	"fmt"

	. "github.com/afritzler/garden-examiner/pkg/data"
	v1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Shoot interface {
	GetName() *ShootName
	GetNamespaceInSeed() (string, error)
	GetManifest() *v1beta1.Shoot
	GetSeedName() string
	GetSeed() (Seed, error)
	GetProject() (Project, error)
	GetSecretRef() (*corev1.SecretReference, error)
	GetInfrastructure() string
	GetState() string
	GetProgress() int
	GetError() string
	Cluster
	RuntimeObjectWrapper
	GardenObject
}

type shoot struct {
	_GardenObject
	cluster
	name          *ShootName
	namespace     string
	seednamespace string
	manifest      v1beta1.Shoot
}

var _ Shoot = &shoot{}

func NewShootFromShootManifest(g Garden, m v1beta1.Shoot) (Shoot, error) {
	n, err := NewShootNameFromShootManifest(g, m)
	if err != nil {
		return nil, err
	}
	s := (&shoot{}).new(g, n, m)
	return s, nil
}

func (s *shoot) new(g Garden, n *ShootName, m v1beta1.Shoot) Shoot {
	m.Kind = "Shoot"
	m.APIVersion = v1beta1.SchemeGroupVersion.String()

	s._GardenObject.new(g)
	s.cluster.new(s)
	s.name = n
	s.manifest = m
	s.namespace = m.GetObjectMeta().GetNamespace()
	return s
}

func (s *shoot) GetName() *ShootName {
	return s.name
}

func (s *shoot) GetNamespaceInSeed() (string, error) {
	if s.seednamespace == "" {
		p, err := s.GetProject()
		if err != nil {
			return "", fmt.Errorf("cannot get namespace for shoot '%s': %s", s.name, err)
		}
		s.seednamespace = fmt.Sprintf("shoot-%s-%s", p.GetName(), s.GetName().GetName())
	}
	return s.seednamespace, nil
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

func (s *shoot) GetProgress() int {
	return s.manifest.Status.LastOperation.Progress
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

//////////////////////////////////////////////////////////////////////////////
// cache

type ShootCacher struct {
	garden Garden
}

func NewShootCacher(g Garden) Cacher {
	return &ShootCacher{g}
}

func (this *ShootCacher) GetAll() (Iterator, error) {
	fmt.Printf("cacher get all shoots\n")
	elems, err := this.garden.GetShoots()
	if err != nil {
		fmt.Printf("cacher got error %s\n", err)
		return nil, err
	}
	fmt.Printf("cacher got %d shoots\n", len(elems))
	a := []interface{}{}
	for _, v := range elems {
		a = append(a, v)
	}
	return NewSliceIterator(a), nil
}

func (this *ShootCacher) Get(key interface{}) (interface{}, error) {
	name := key.(ShootName)
	return this.garden.GetShoot(&name)
}

func (this *ShootCacher) Key(elem interface{}) interface{} {
	return elem.(Shoot).GetName()
}

type ShootCache interface {
	GetShoots() (map[ShootName]Shoot, error)
	GetShoot(name *ShootName) (Shoot, error)
	Reset()
}

type shoot_cache struct {
	cache Cache
}

func NewShootCache(g Garden) ShootCache {
	return &shoot_cache{NewCache(NewShootCacher(g))}
}

func (this *shoot_cache) Reset() {
	this.cache.Reset()
}

func (this *shoot_cache) GetShoots() (map[ShootName]Shoot, error) {
	m := map[ShootName]Shoot{}
	i, err := this.cache.GetAll()
	if err != nil {
		return nil, err
	}
	for i.HasNext() {
		e := i.Next().(Shoot)
		m[*e.GetName()] = e
	}
	return m, nil
}

func (this *shoot_cache) GetShoot(name *ShootName) (Shoot, error) {
	e, err := this.cache.Get(*name)
	if err != nil {
		return nil, err
	}
	return e.(Shoot), nil
}
