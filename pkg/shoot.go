package gube

import (
	"fmt"

	. "github.com/afritzler/garden-examiner/pkg/data"
	v1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Shoot interface {
	GetName() *ShootName
	GetNamespaceInSeed() string
	GetManifest() *v1beta1.Shoot
	GetDomainName() string
	GetSeedName() string
	GetSeed() (Seed, error)
	GetProject() (Project, error)
	GetSecretRef() (*corev1.SecretReference, error)
	GetSecretContentFromSeed(name string) (map[string]string, error)
	GetBasicAuth() (user, password string, err error)
	GetCloudProviderConfig() (map[string]string, error)
	GetConfigMapEntriesFromSeed(name string) (map[string]string, error)
	GetTerraformJobData(job string, data string) (string, error)
	GetIngressFromSeed(name string) (*extv1beta1.Ingress, error)
	GetIngressHostFromSeed(name string) (string, error)
	GetInfrastructure() string
	GetInfrastructureConfig() interface{}
	GetRegion() string
	GetIaaSInfo() (IaaSInfo, error)
	GetProfileName() string
	GetProfile() (Profile, error)
	GetReconcilationState() string
	GetReconcilationError() string
	GetReconcilationProgress() int
	GetState() string
	GetError() string
	GetConditionErrors() map[string]string
	Cluster
	RuntimeObjectWrapper
	GardenObject
}

type shoot struct {
	_GardenObject
	cluster
	name          *ShootName
	namespace     string
	seednamespace string
	manifest      v1beta1.Shoot
}

var _ Shoot = &shoot{}

func NewShootFromShootManifest(g Garden, m v1beta1.Shoot) (Shoot, error) {
	n, err := NewShootNameFromShootManifest(g, m)
	if err != nil {
		return nil, err
	}
	s := (&shoot{}).new(g, n, m)
	return s, nil
}

func (s *shoot) new(g Garden, n *ShootName, m v1beta1.Shoot) Shoot {
	m.Kind = "Shoot"
	m.APIVersion = v1beta1.SchemeGroupVersion.String()

	s._GardenObject.new(g)
	s.cluster.new("shoot "+n.String(), s)
	s.name = n
	s.manifest = m
	s.namespace = m.GetObjectMeta().GetNamespace()
	return s
}

func (s *shoot) GetName() *ShootName {
	return s.name
}

func (s *shoot) AsShoot() (Shoot, error) {
	return s, nil
}
func (s *shoot) GetShootName() *ShootName {
	return s.GetName()
}

func (s *shoot) GetNamespaceInSeed() string {
	return s.manifest.Status.TechnicalID
}

func (s *shoot) GetNamespace() (string, error) {
	if s.namespace == "" {
		p, err := s.GetProject()
		if err != nil {
			return "", fmt.Errorf("cannot get namespace for shoot '%s': %s", s.name, err)
		}
		s.namespace = p.GetNamespace()
	}
	return s.namespace, nil
}

func (s *shoot) GetManifest() *v1beta1.Shoot {
	return &s.manifest
}

func (s *shoot) GetRuntimeObject() runtime.Object {
	return &s.manifest
}

func (s *shoot) GetSeedName() string {
	return *s.manifest.Spec.Cloud.Seed
}

func (s *shoot) GetDomainName() string {
	return *s.manifest.Spec.DNS.Domain
}

func (s *shoot) GetSeed() (Seed, error) {
	// should never fail with panic :-P
	return s.garden.GetSeed(s.GetSeedName())
}

func (s *shoot) GetProject() (Project, error) {
	return s.garden.GetProject(s.name.GetProjectName())
}

func (s *shoot) GetState() string {
	state := s.GetReconcilationState()
	if state == "Succeeded" {
		if s.GetConditionErrors() != nil {
			return "Problem"
		}
	}
	return state
}

func (s *shoot) GetReconcilationState() string {
	if s.manifest.Status.LastOperation == nil {
		return "unknown"
	}
	return string(s.manifest.Status.LastOperation.State)
}

func (s *shoot) GetReconcilationProgress() int {
	if s.manifest.Status.LastOperation == nil {
		return 0
	}
	return s.manifest.Status.LastOperation.Progress
}

func (s *shoot) GetError() string {
	cond := s.GetConditionErrors()
	e := s.GetReconcilationError()
	if cond != nil {
		for n, m := range cond {
			if e != "" {
				e = e + "\n"

			}
			e = e + n + ": " + m
		}
	}
	return e
}

func (s *shoot) GetReconcilationError() string {
	if s.manifest.Status.LastOperation == nil {
		return ""
	}
	if s.manifest.Status.LastOperation.State != v1beta1.ShootLastOperationStateSucceeded {
		if s.manifest.Status.LastError != nil {
			return s.manifest.Status.LastError.Description
		}
	}
	return ""
}

func (s *shoot) GetConditionErrors() map[string]string {
	errors := map[string]string{}
	for _, c := range s.manifest.Status.Conditions {
		if c.Status == "False" {
			errors[string(c.Type)] = c.Message
		}
	}
	if len(errors) == 0 {
		return nil
	}
	return errors
}

func (s *shoot) GetKubeconfig() ([]byte, error) {
	ref, err := s.GetSecretRef()
	if err != nil {
		return nil, fmt.Errorf("could not get secret ref for shoot '%s': %s", s.name, err)
	}
	secret, err := s.garden.GetSecretByRef(*ref)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret for shoot '%s': %s", s.name, err)
	}
	return secret.Data[secretkubeconfig], nil
}

func (s *shoot) GetSecretRef() (*corev1.SecretReference, error) {
	ns, err := s.GetNamespace()
	if err != nil {
		return nil, err
	}
	return &corev1.SecretReference{Name: fmt.Sprintf("%s.kubeconfig", s.name.GetName()), Namespace: ns}, nil
}

func (s *shoot) GetCloudProviderConfig() (map[string]string, error) {
	return s.GetSecretContentFromSeed("cloudprovider")
}

func (s *shoot) GetBasicAuth() (user, password string, err error) {
	content, err := s.GetSecretContentFromSeed("kubecfg")
	if err != nil {
		return "", "", err
	}
	user, ok := content["username"]
	if !ok {
		return "", "", fmt.Errorf("no user configured for shoot '%s'", s.GetName())
	}
	pass, ok := content["password"]
	if !ok {
		return "", "", fmt.Errorf("no password configured for shoot '%s'", s.GetName())
	}
	return user, pass, nil
}

func (s *shoot) GetSecretContentFromSeed(name string) (map[string]string, error) {
	ns := s.GetNamespaceInSeed()
	seed, err := s.GetSeed()
	if err != nil {
		return nil, err
	}
	secret, err := seed.GetSecretByRef(corev1.SecretReference{Name: name, Namespace: ns})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret '%s' for shoot '%s': %s", name, s.name, err)
	}
	config := map[string]string{}
	for k, v := range secret.Data {
		config[k] = string(v)
	}
	return config, nil
}

func (s *shoot) GetTerraformJobData(job string, data string) (string, error) {
	switch job {
	case "infra":
	case "external-dns":
	case "internal-dns":
	case "ingress":
	default:
		return "", fmt.Errorf("invalid job '%s', select one of infra, external-dns, internal-dns, ingress", job)
	}
	cm := ""
	field := ""
	switch data {
	case "state":
		cm = "state"
		field = "terraform.tfstate"
	case "config":
		cm = "config"
		field = "variables.tf"
	case "script":
		cm = "config"
		field = "main.tf"
	default:
		return "", fmt.Errorf("invalid data '%s', select one of: state, config, script")
	}
	entries, err := s.GetConfigMapEntriesFromSeed(fmt.Sprintf("%s.%s.tf-%s", s.GetName().GetName(), job, cm))
	if err != nil {
		return "", err
	}
	return entries[field], nil
}

func (s *shoot) GetConfigMapEntriesFromSeed(name string) (map[string]string, error) {
	ns := s.GetNamespaceInSeed()
	seed, err := s.GetSeed()
	if err != nil {
		return nil, err
	}
	return seed.GetConfigMapEntries(name, ns)
}

func (s *shoot) GetIngressFromSeed(name string) (*extv1beta1.Ingress, error) {
	ns := s.GetNamespaceInSeed()
	seed, err := s.GetSeed()
	if err != nil {
		return nil, err
	}
	return seed.GetIngress(name, ns)
}

func (s *shoot) GetIngressHostFromSeed(name string) (string, error) {
	ingress, err := s.GetIngressFromSeed(name)
	if err != nil {
		return "", err
	}
	for _, r := range ingress.Spec.Rules {
		return r.Host, nil
	}
	return "", fmt.Errorf("no rule entry found for ingress '%s' of '%s' in seed %s",
		name, s.GetName(), s.GetName())
}

func (s *shoot) GetInfrastructure() string {
	if s.manifest.Spec.Cloud.AWS != nil {
		return "aws"
	}
	if s.manifest.Spec.Cloud.Azure != nil {
		return "azure"
	}
	if s.manifest.Spec.Cloud.OpenStack != nil {
		return "openstack"
	}
	if s.manifest.Spec.Cloud.GCP != nil {
		return "gcp"
	}
	if s.manifest.Spec.Cloud.Local != nil {
		return "local"
	}
	return "unknown"
}

func (s *shoot) GetInfrastructureConfig() interface{} {
	if s.manifest.Spec.Cloud.AWS != nil {
		return s.manifest.Spec.Cloud.AWS
	}
	if s.manifest.Spec.Cloud.Azure != nil {
		return s.manifest.Spec.Cloud.Azure
	}
	if s.manifest.Spec.Cloud.OpenStack != nil {
		return s.manifest.Spec.Cloud.OpenStack
	}
	if s.manifest.Spec.Cloud.GCP != nil {
		return s.manifest.Spec.Cloud.GCP
	}
	if s.manifest.Spec.Cloud.Local != nil {
		return s.manifest.Spec.Cloud.Local
	}
	return nil
}

func (s *shoot) GetIaaSInfo() (IaaSInfo, error) {
	k := s.GetInfrastructure()
	h := iaas[k]
	if h == nil {
		return nil, fmt.Errorf("no implementation for IaaS type '%s'", k)
	}
	return h.GetIaaSInfo(s)
}

func (s *shoot) GetRegion() string {
	return s.manifest.Spec.Cloud.Region
}

func (s *shoot) GetProfileName() string {
	return s.manifest.Spec.Cloud.Profile
}

func (s *shoot) GetProfile() (Profile, error) {
	name := s.GetProfileName()
	if name == "" {
		return nil, fmt.Errorf("no profile found for shoot %s", s.GetName())
	}
	return s.garden.GetProfile(name)
}

//////////////////////////////////////////////////////////////////////////////
// cache

type ShootCacher struct {
	garden Garden
}

func NewShootCacher(g Garden) Cacher {
	return &ShootCacher{g}
}

func (this *ShootCacher) GetAll() (Iterator, error) {
	fmt.Printf("cacher get all shoots\n")
	elems, err := this.garden.GetShoots()
	if err != nil {
		fmt.Printf("cacher got error %s\n", err)
		return nil, err
	}
	fmt.Printf("cacher got %d shoots\n", len(elems))
	a := []interface{}{}
	for _, v := range elems {
		a = append(a, v)
	}
	return NewSliceIterator(a), nil
}

func (this *ShootCacher) Get(key interface{}) (interface{}, error) {
	name := key.(ShootName)
	return this.garden.GetShoot(&name)
}

func (this *ShootCacher) Key(elem interface{}) interface{} {
	return *elem.(Shoot).GetName()
}

type ShootCache interface {
	GetShoots() (map[ShootName]Shoot, error)
	GetShoot(name *ShootName) (Shoot, error)
	Reset()
}

type shoot_cache struct {
	cache Cache
}

func NewShootCache(g Garden) ShootCache {
	return &shoot_cache{NewCache(NewShootCacher(g))}
}

func (this *shoot_cache) Reset() {
	this.cache.Reset()
}

func (this *shoot_cache) GetShoots() (map[ShootName]Shoot, error) {
	m := map[ShootName]Shoot{}
	i, err := this.cache.GetAll()
	if err != nil {
		return nil, err
	}
	for i.HasNext() {
		e := i.Next().(Shoot)
		m[*e.GetName()] = e
	}
	return m, nil
}

func (this *shoot_cache) GetShoot(name *ShootName) (Shoot, error) {
	e, err := this.cache.Get(*name)
	if err != nil {
		return nil, err
	}
	return e.(Shoot), nil
}
