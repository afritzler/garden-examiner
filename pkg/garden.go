package gube

import (
	"fmt"

	gardenclientset "github.com/gardener/gardener/pkg/client/garden/clientset/versioned"
	"github.com/gardener/gardener/pkg/operation/common"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

type Garden interface {
	GetShoots() (map[ShootName]Shoot, error)
	GetSeeds() (map[string]Seed, error)
	GetSeed(name string) (Seed, error)
	GetProjects() (map[string]Project, error)
	GetProject(name string) (Project, error)
	GetProjectByNamespace(namespace string) (Project, error)
	GetSecretByRef(secretref corev1.SecretReference) (*corev1.Secret, error)
}

type garden struct {
	gardenset *gardenclientset.Clientset
	kubeset   *kubernetes.Clientset
}

func NewGarden(config *restclient.Config) (Garden, error) {
	gardenset, err := gardenclientset.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate garden client: %s", err)
	}
	kubeset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate kubernetes client: %s", err)
	}
	return &garden{gardenset: gardenset, kubeset: kubeset}, nil
}

func (g *garden) GetShoots() (map[ShootName]Shoot, error) {
	shoots, err := g.gardenset.GardenV1beta1().Shoots("").List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get shoots: %s", err)
	}
	result := map[ShootName]Shoot{}
	for _, s := range shoots.Items {
		shoot, err := NewShootFromShootManifest(g, s)
		if err != nil {
			return nil, err
		}
		result[*shoot.GetName()] = shoot
	}
	return result, nil
}

func (g *garden) GetSeeds() (map[string]Seed, error) {
	seeds, err := g.gardenset.GardenV1beta1().Seeds().List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get seeds: %s", err)
	}
	result := map[string]Seed{}
	for _, s := range seeds.Items {
		seed := NewSeedFromSeedManifest(g, s)
		result[seed.GetName()] = seed
	}
	return result, nil
}

func (g *garden) GetSeed(name string) (Seed, error) {
	m, err := g.gardenset.GardenV1beta1().Seeds().Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get seed %s: %s", name, err)
	}
	return NewSeedFromSeedManifest(g, *m), nil
}

func (g *garden) GetProjects() (map[string]Project, error) {
	namespaces, err := g.kubeset.CoreV1().Namespaces().List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", common.GardenRole, common.GardenRoleProject),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get project namespaces: %s", err)
	}
	result := map[string]Project{}
	for _, n := range namespaces.Items {
		project := NewProjectFromNamespaceManifest(&n)
		result[project.GetName()] = project
	}
	return result, nil
}

func (g *garden) GetProject(name string) (Project, error) {
	namespaces, err := g.kubeset.CoreV1().Namespaces().List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", common.ProjectName, name),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get project %s: %s", name, err)
	}
	if len(namespaces.Items) > 0 {
		return NewProjectFromNamespaceManifest(&namespaces.Items[0]), nil
	} else {
		return nil, fmt.Errorf("failed to get project: got an empty project list for %s", name)
	}
}

func (g *garden) GetProjectByNamespace(namespace string) (Project, error) {
	n, err := g.kubeset.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get project namespace %s: %s", namespace, err)
	}
	return NewProjectFromNamespaceManifest(n), nil
}

func (g *garden) GetSecretByRef(secretref corev1.SecretReference) (*corev1.Secret, error) {
	secret, err := g.kubeset.CoreV1().Secrets(secretref.Namespace).Get(secretref.Name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %s for namespace %s: %s", secretref.Name, secretref.Namespace, err)
	}
	return secret, nil
}
