package gube

type CachedGarden interface {
	Reset()
	Garden
}

type cached_garden struct {
	Garden
	profiles ProfileCache
	shoots   ShootCache
}

var _ Garden = &cached_garden{}

func NewCachedGarden(g Garden) CachedGarden {
	n := &cached_garden{nil, NewProfileCache(g), NewShootCache(g)}
	n.Garden = g.NewWrapper(g)
	return n
}

func (this *cached_garden) Reset() {
	this.profiles.Reset()
	this.shoots.Reset()
}

func (this *cached_garden) GetProfile(name string) (Profile, error) {
	return this.profiles.GetProfile(name)
}

func (this *cached_garden) GetProfiles() (map[string]Profile, error) {
	return this.profiles.GetProfiles()
}

func (this *cached_garden) GetShoot(name *ShootName) (Shoot, error) {
	return this.shoots.GetShoot(name)
}

func (this *cached_garden) GetShoots() (map[ShootName]Shoot, error) {
	return this.shoots.GetShoots()
}
