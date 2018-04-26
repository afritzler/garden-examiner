package gube

import (
	corev1 "k8s.io/api/core/v1"
)

type Project interface {
	GetName() string
	GetNamespace() string
}

type project struct {
	name      string
	namespace string
}

func NewProjectFromNamespaceManifest(n *corev1.Namespace) Project {
	return &project{name: GetProjectNameFromNamespaceManifest(n), namespace: n.GetName()}
}

func (p *project) GetName() string {
	return p.name
}

func (p *project) GetNamespace() string {
	return p.namespace
}
