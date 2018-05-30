package gube

import (
	"fmt"
	"io/ioutil"

	gardenclientset "github.com/gardener/gardener/pkg/client/garden/clientset/versioned"
	"github.com/gardener/gardener/pkg/operation/common"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

type garden_access struct {
	gardenset  *gardenclientset.Clientset
	kubeset    *kubernetes.Clientset
	kubeconfig []byte
}

func newGardenAccess(config *restclient.Config) (*garden_access, error) {
	gardenset, err := gardenclientset.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate garden client: %s", err)
	}
	kubeset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate kubernetes client: %s", err)
	}
	return &garden_access{gardenset: gardenset, kubeset: kubeset}, nil
}

func newGardenAccessFromBytes(bytes []byte) (*garden_access, error) {
	config, err := NewConfigFromBytes(bytes)
	if err != nil {
		return nil, fmt.Errorf("cannot create kubeconfig: %s", err)
	}

	g, err := newGardenAccess(config)
	if err == nil {
		g.kubeconfig = bytes
	}
	return g, err
}

func newGardenAccessFromConfigfile(configfile string) (*garden_access, error) {
	bytes, err := ioutil.ReadFile(configfile)
	if err != nil {
		return nil, fmt.Errorf("cannot read kubeconfig '%s': %s", configfile, err)
	}
	a, err := newGardenAccessFromBytes(bytes)
	if err != nil {
		return nil, fmt.Errorf("cannot read kubeconfig '%s': %s", configfile, err)
	}
	return a, nil
}

func (this *garden_access) GetKubeconfig() []byte {
	return this.kubeconfig
}

func (this *garden_access) GetShoots(eff Garden) (map[ShootName]Shoot, error) {
	shoots, err := this.gardenset.GardenV1beta1().Shoots("").List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get shoots: %s", err)
	}
	result := map[ShootName]Shoot{}
	for _, s := range shoots.Items {
		shoot, err := NewShootFromShootManifest(eff, s)
		if err != nil {
			return nil, err
		}
		result[*shoot.GetName()] = shoot
	}
	return result, nil
}

func (this *garden_access) GetShoot(eff Garden, name *ShootName) (Shoot, error) {
	project, err := eff.GetProject(name.GetProjectName())
	if err != nil {
		return nil, err
	}
	m, err := this.gardenset.GardenV1beta1().Shoots(project.GetNamespace()).Get(name.GetName(), metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get shoot %s: %s", *name, err)
	}
	return NewShootFromShootManifest(eff, *m)
}

func (this *garden_access) GetSeeds(eff Garden) (map[string]Seed, error) {
	seeds, err := this.gardenset.GardenV1beta1().Seeds().List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get seeds: %s", err)
	}
	result := map[string]Seed{}
	for _, s := range seeds.Items {
		seed := NewSeedFromSeedManifest(eff, s)
		result[seed.GetName()] = seed
	}
	return result, nil
}

func (this *garden_access) GetSeed(eff Garden, name string) (Seed, error) {
	m, err := this.gardenset.GardenV1beta1().Seeds().Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get seed %s: %s", name, err)
	}
	return NewSeedFromSeedManifest(eff, *m), nil
}

func (this *garden_access) GetProjects(eff Garden) (map[string]Project, error) {
	namespaces, err := this.kubeset.CoreV1().Namespaces().List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", common.GardenRole, common.GardenRoleProject),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get project namespaces: %s", err)
	}
	result := map[string]Project{}
	for _, n := range namespaces.Items {
		project := NewProjectFromNamespaceManifest(eff, &n)
		result[project.GetName()] = project
	}
	return result, nil
}

func (this *garden_access) GetProject(eff Garden, name string) (Project, error) {
	namespaces, err := this.kubeset.CoreV1().Namespaces().List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", common.ProjectName, name),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get project %s: %s", name, err)
	}
	if len(namespaces.Items) > 0 {
		return NewProjectFromNamespaceManifest(eff, &namespaces.Items[0]), nil
	} else {
		return nil, fmt.Errorf("failed to get project: got an empty project list for %s", name)
	}
}

func (this *garden_access) GetProjectByNamespace(eff Garden, namespace string) (Project, error) {
	n, err := this.kubeset.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get project namespace %s: %s", namespace, err)
	}
	return NewProjectFromNamespaceManifest(eff, n), nil
}

func (this *garden_access) GetProfiles(eff Garden) (map[string]Profile, error) {
	elems, err := this.gardenset.GardenV1beta1().CloudProfiles().List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get cloud profiles: %s", err)
	}
	result := map[string]Profile{}
	for _, s := range elems.Items {
		elem := NewProfileFromProfileManifest(eff, s)
		result[elem.GetName()] = elem
	}
	return result, nil
}

func (this *garden_access) GetProfile(eff Garden, name string) (Profile, error) {
	//fmt.Printf("GET PROFILE %s\n", name)
	m, err := this.gardenset.GardenV1beta1().CloudProfiles().Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get cloud profile %s: %s", name, err)
	}
	return NewProfileFromProfileManifest(eff, *m), nil
}

func (this *garden_access) GetSecretByRef(eff Garden, secretref corev1.SecretReference) (*corev1.Secret, error) {
	secret, err := this.kubeset.CoreV1().Secrets(secretref.Namespace).Get(secretref.Name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %s for namespace %s: %s", secretref.Name, secretref.Namespace, err)
	}
	return secret, nil
}
