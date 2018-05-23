package gube

import (
	. "github.com/afritzler/garden-examiner/pkg/data"

	v1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Profile interface {
	GetName() string
	GetManifest() *v1beta1.CloudProfile
	GetInfrastructure() string
	RuntimeObjectWrapper
	GardenObject
}

type profile struct {
	_GardenObject
	name     string
	manifest v1beta1.CloudProfile
}

func NewProfileFromProfileManifest(g Garden, m v1beta1.CloudProfile) Profile {
	return (&profile{}).new(g, m)
}

func (s *profile) new(g Garden, m v1beta1.CloudProfile) Profile {
	m.Kind = "CloudProfile"
	m.APIVersion = v1beta1.SchemeGroupVersion.String()

	s._GardenObject.new(g)
	s.name = m.GetName()
	s.manifest = m
	return s
}
func (s *profile) GetName() string {
	return s.name
}

func (s *profile) GetManifest() *v1beta1.CloudProfile {
	return &s.manifest
}

func (s *profile) GetRuntimeObject() runtime.Object {
	return &s.manifest
}

func (s *profile) GetInfrastructure() string {
	if s.manifest.Spec.AWS != nil {
		return "aws"
	}
	if s.manifest.Spec.Azure != nil {
		return "azure"
	}
	if s.manifest.Spec.OpenStack != nil {
		return "openstack"
	}
	if s.manifest.Spec.GCP != nil {
		return "gcp"
	}
	if s.manifest.Spec.Local != nil {
		return "local"
	}
	return "unknown"
}

//////////////////////////////////////////////////////////////////////////////
// cache

type ProfileCacher struct {
	garden Garden
}

func NewProfileCacher(g Garden) Cacher {
	return &ProfileCacher{g}
}

func (this *ProfileCacher) GetAll() (Iterator, error) {
	elems, err := this.garden.GetProfiles()
	if err != nil {
		return nil, err
	}
	a := []interface{}{}
	for _, v := range elems {
		a = append(a, v)
	}
	return NewSliceIterator(a), nil
}

func (this *ProfileCacher) Get(key interface{}) (interface{}, error) {
	name := key.(string)
	return this.garden.GetProfile(name)
}

func (this *ProfileCacher) Key(elem interface{}) interface{} {
	return elem.(Profile).GetName()
}

type ProfileCache interface {
	GetProfiles() (map[string]Profile, error)
	GetProfile(name string) (Profile, error)
	Reset()
}

type profile_cache struct {
	cache Cache
}

func NewProfileCache(g Garden) ProfileCache {
	return &profile_cache{NewCache(NewProfileCacher(g))}
}

func (this *profile_cache) Reset() {
	this.cache.Reset()
}

func (this *profile_cache) GetProfiles() (map[string]Profile, error) {
	m := map[string]Profile{}
	i, err := this.cache.GetAll()
	if err != nil {
		return nil, err
	}
	for i.HasNext() {
		e := i.Next().(Profile)
		m[e.GetName()] = e
	}
	return m, nil
}

func (this *profile_cache) GetProfile(name string) (Profile, error) {
	e, err := this.cache.Get(name)
	if err != nil {
		return nil, err
	}
	return e.(Profile), nil
}
