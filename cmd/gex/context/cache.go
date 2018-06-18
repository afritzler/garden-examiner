package context

import (
	"github.com/afritzler/garden-examiner/pkg"
	"github.com/mandelsoft/filepath/pkg/filepath"
)

func (this *Context) CacheDir() string {
	if this.Gexdir == "" {
		return ""
	}
	return filepath.Join(this.Gexdir, "cache", this.Name)
}

func (this *Context) CacheDirForShoot(s gube.Shoot) string {
	if this.Gexdir == "" {
		return ""
	}
	return filepath.Join(this.CacheDir(), "projects", s.GetName().GetProjectName(), s.GetName().GetName())
}

func (this *Context) CacheDirForSeed(s gube.Seed) string {
	if this.Gexdir == "" {
		return ""
	}
	return filepath.Join(this.CacheDir(), "seeds", s.GetName())
}

func (this *Context) CacheDirForGarden() string {
	if this.Gexdir == "" {
		return ""
	}
	return filepath.Join(this.CacheDir(), "garden")
}

func (this *Context) CacheDirFor(e interface{}) string {
	if this.Gexdir == "" {
		return ""
	}
	switch r := e.(type) {
	case gube.Shoot:
		return this.CacheDirForShoot(r)
	case gube.Seed:
		return this.CacheDirForSeed(r)
	}
	return ""
}
