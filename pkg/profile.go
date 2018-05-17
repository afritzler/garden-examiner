package gube

import (
	"fmt"

	v1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Profile interface {
	GetName() string
	GetManifest() *v1beta1.CloudProfile
	GetInfrastructure() string
	RuntimeObjectWrapper
}

type profile struct {
	garden   Garden
	name     string
	manifest v1beta1.CloudProfile
}

func NewProfileFromProfileManifest(g Garden, m v1beta1.CloudProfile) Profile {
	m.Kind = "CloudProfile"
	m.APIVersion = v1beta1.SchemeGroupVersion.String()
	return &profile{garden: g, name: m.GetName(), manifest: m}
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

type ProfileCache interface {
	GetProfiles() (map[string]Profile, error)
	GetProfile(name string) (Profile, error)
	Reset()
}

type profile_cache struct {
	garden   Garden
	profiles map[string]Profile
	complete bool
}

func NewProfileCache(g Garden) ProfileCache {
	return &profile_cache{g, nil, false}
}

func (this *profile_cache) Reset() {
	this.profiles = nil
	this.complete = false
}

func (this *profile_cache) GetProfiles() (map[string]Profile, error) {
	if this.profiles == nil || !this.complete {
		elems, err := this.garden.GetProfiles()
		if err != nil {
			return nil, err
		}
		this.profiles = elems
		this.complete = true
	}
	return this.profiles, nil
}

func (this *profile_cache) GetProfile(name string) (Profile, error) {
	var p Profile = nil
	if this.profiles != nil {
		p = this.profiles[name]
	}
	if p == nil && !this.complete {
		elem, err := this.garden.GetProfile(name)
		if err != nil {
			return nil, err
		}
		if this.profiles == nil {
			this.profiles = map[string]Profile{}
		}
		this.profiles[name] = elem
		p = elem
	}
	if p == nil {
		return nil, fmt.Errorf("profile '%s' not found", name)
	}
	return p, nil
}
