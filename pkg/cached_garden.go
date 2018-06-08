package gube

import (
	"fmt"
)

var _ = fmt.Errorf

type CachedGarden interface {
	Reset()
	Garden
}

type cached_garden struct {
	Garden
	projects ProjectCache
	profiles ProfileCache
	shoots   ShootCache
}

var _ Garden = &cached_garden{}

func NewCachedGarden(g Garden) CachedGarden {
	return (&cached_garden{}).new(g)
}

func (this *cached_garden) new(g Garden) CachedGarden {
	this.Garden = g.NewWrapper(this)
	this.projects = NewProjectCache(this.Garden)
	this.profiles = NewProfileCache(this.Garden)
	this.shoots = NewShootCache(this.Garden)
	return this
}

func (this *cached_garden) Reset() {
	this.projects.Reset()
	this.profiles.Reset()
	this.shoots.Reset()
}

func (this *cached_garden) GetProject(name string) (Project, error) {
	//fmt.Printf("GET CACHED  %s\n", name)
	return this.projects.GetProject(name)
}

func (this *cached_garden) GetProjects() (map[string]Project, error) {
	return this.projects.GetProjects()
}

func (this *cached_garden) GetProjectByNamespace(namespace string) (Project, error) {
	return this.projects.GetProjectByNamespace(namespace)
}

func (this *cached_garden) GetProfile(name string) (Profile, error) {
	//fmt.Printf("GET CACHED PROFILE %s\n", name)
	return this.profiles.GetProfile(name)
}

func (this *cached_garden) GetProfiles() (map[string]Profile, error) {
	return this.profiles.GetProfiles()
}

func (this *cached_garden) GetShoot(name *ShootName) (Shoot, error) {
	//fmt.Printf("GET CACHED SHOOT %s\n", name)
	return this.shoots.GetShoot(name)
}

func (this *cached_garden) GetShoots() (map[ShootName]Shoot, error) {
	return this.shoots.GetShoots()
}
