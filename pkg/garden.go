package gube

import (
	"fmt"

	restclient "k8s.io/client-go/rest"

	_ "github.com/afritzler/garden-examiner/pkg/data"
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
	Cluster
}

type garden struct {
	cluster
	access    *garden_access
	effective Garden
}

var _ Garden = &garden{}

func NewGarden(config *restclient.Config) (Garden, error) {
	access, err := newGardenAccess(config)
	if err != nil {
		return nil, err
	}
	return (&garden{}).new(access, nil), nil
}

func NewGardenFromConfigfile(configfile string) (Garden, error) {
	access, err := newGardenAccessFromConfigfile(configfile)
	if err != nil {
		return nil, err
	}
	return (&garden{}).new(access, nil), nil

}

func (g *garden) new(access *garden_access, eff Garden) *garden {
	g.cluster.new(g)
	if eff == nil {
		eff = g
	}
	g.effective = eff
	g.access = access
	return g
}

func (this *garden) NewWrapper(g Garden) Garden {
	return (&garden{}).new(this.access, g)
}

func (this *garden) GetKubeconfig() ([]byte, error) {
	cfg := this.access.GetKubeconfig()
	if cfg == nil {
		return nil, fmt.Errorf("no kubeconfig available for garden")
	}
	return cfg, nil
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
