package gube

import (
	corev1 "k8s.io/api/core/v1"
	restclient "k8s.io/client-go/rest"
)

type Garden interface {
	NewWrapper(g Garden) Garden
	GetShoots() (map[ShootName]Shoot, error)
	GetShoot(*ShootName) (Shoot, error)
	GetSeeds() (map[string]Seed, error)
	GetSeed(name string) (Seed, error)
	GetProjects() (map[string]Project, error)
	GetProject(name string) (Project, error)
	GetProjectByNamespace(namespace string) (Project, error)
	GetProfiles() (map[string]Profile, error)
	GetProfile(name string) (Profile, error)
	GetSecretByRef(secretref corev1.SecretReference) (*corev1.Secret, error)
}

type garden struct {
	access    *garden_access
	effective Garden
}

func NewGarden(config *restclient.Config) (Garden, error) {
	access, err := newGardenAccess(config)
	if err != nil {
		return nil, err
	}
	g := &garden{access, nil}
	g.effective = g
	return g, nil
}

func (this *garden) NewWrapper(g Garden) Garden {
	return &garden{this.access, g}
}

func (this *garden) GetShoots() (map[ShootName]Shoot, error) {
	return this.access.GetShoots(this.effective)
}

func (this *garden) GetShoot(name *ShootName) (Shoot, error) {
	return this.access.GetShoot(this.effective, name)
}

func (this *garden) GetSeeds() (map[string]Seed, error) {
	return this.access.GetSeeds(this.effective)
}

func (this *garden) GetSeed(name string) (Seed, error) {
	return this.access.GetSeed(this.effective, name)
}

func (this *garden) GetProjects() (map[string]Project, error) {
	return this.access.GetProjects(this.effective)
}

func (this *garden) GetProject(name string) (Project, error) {
	return this.access.GetProject(this.effective, name)
}

func (this *garden) GetProjectByNamespace(namespace string) (Project, error) {
	return this.access.GetProjectByNamespace(this.effective, namespace)
}

func (this *garden) GetProfiles() (map[string]Profile, error) {
	return this.access.GetProfiles(this.effective)
}

func (this *garden) GetProfile(name string) (Profile, error) {
	return this.access.GetProfile(this.effective, name)
}

func (this *garden) GetSecretByRef(secretref corev1.SecretReference) (*corev1.Secret, error) {
	return this.access.GetSecretByRef(this.effective, secretref)
}
