package gube

import (
	corev1 "k8s.io/api/core/v1"
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
