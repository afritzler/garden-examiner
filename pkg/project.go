package gube

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	. "github.com/afritzler/garden-examiner/pkg/data"
)

type Project interface {
	GetName() string
	GetNamespace() string
	GardenObject
}

type project struct {
	_GardenObject
	name      string
	namespace string
}

func NewProjectFromNamespaceManifest(g Garden, n *corev1.Namespace) Project {
	return (&project{}).new(g, GetProjectNameFromNamespaceManifest(n), n.GetName())
}

func (p *project) new(g Garden, n string, ns string) Project {
	p._GardenObject.new(g)
	p.name = n
	p.namespace = ns
	return p
}
func (p *project) GetName() string {
	return p.name
}

func (p *project) GetNamespace() string {
	return p.namespace
}

//////////////////////////////////////////////////////////////////////////////
// cache

type ProjectCacher struct {
	garden Garden
}

func NewProjectCacher(g Garden) Cacher {
	return &ProjectCacher{g}
}

func (this *ProjectCacher) GetAll() (Iterator, error) {
	elems, err := this.garden.GetProjects()
	if err != nil {
		return nil, err
	}
	a := []interface{}{}
	for _, v := range elems {
		a = append(a, v)
	}
	return NewSliceIterator(a), nil
}

func (this *ProjectCacher) Get(key interface{}) (interface{}, error) {
	name := key.(string)
	return this.garden.GetProject(name)
}

func (this *ProjectCacher) Key(elem interface{}) interface{} {
	return elem.(Project).GetName()
}

type ProjectCache interface {
	GetProjects() (map[string]Project, error)
	GetProject(name string) (Project, error)
	GetProjectByNamespace(namespace string) (Project, error)
	Reset()
}

type project_cache struct {
	cache       UnsyncedCache
	byNamespace map[string]Project
}

func NewProjectCache(g Garden) ProjectCache {
	return &project_cache{NewCache(NewProjectCacher(g)).(UnsyncedCache), map[string]Project{}}
}

func (this *project_cache) Reset() {
	this.cache.Reset()
	this.byNamespace = map[string]Project{}
}

func (this *project_cache) GetProjectByNamespace(namespace string) (Project, error) {
	this.cache.Lock()
	defer this.cache.Unlock()
	if len(this.byNamespace) == 0 {
		it, err := this.cache.NotSynced_GetAll()
		if err != nil {
			return nil, err
		}
		for it.HasNext() {
			n := it.Next().(Project)
			this.byNamespace[n.GetNamespace()] = n
		}
	}
	p, ok := this.byNamespace[namespace]
	if !ok {
		return nil, fmt.Errorf("no project found for namespace '%s'", namespace)
	}
	return p, nil
}

func (this *project_cache) GetProjects() (map[string]Project, error) {
	m := map[string]Project{}
	i, err := this.cache.GetAll()
	if err != nil {
		return nil, err
	}
	for i.HasNext() {
		e := i.Next().(Project)
		m[e.GetName()] = e
	}
	return m, nil
}

func (this *project_cache) GetProject(name string) (Project, error) {
	e, err := this.cache.Get(name)
	if err != nil {
		return nil, err
	}
	return e.(Project), nil
}
