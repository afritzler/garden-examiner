package context

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/pkg"
)

type Context struct {
	Name            string
	GardenSetConfig gube.GardenSetConfig
	GardenConfig    gube.GardenConfig
	Garden          gube.CachedGarden
}

func Get(opts *cmdint.Options) *Context {
	for opts != nil {
		if opts.Context != nil {
			return opts.Context.(*Context)
		}
		opts = opts.Parent
	}
	return nil
}

func (this *Context) GetProfile(name string) (gube.Profile, error) {

	return this.Garden.GetProfile(name)
}

func (this *Context) GetShoot(name *gube.ShootName) (gube.Shoot, error) {
	return this.Garden.GetShoot(name)
}

func (this *Context) GetShoots() (map[gube.ShootName]gube.Shoot, error) {
	return this.Garden.GetShoots()
}
