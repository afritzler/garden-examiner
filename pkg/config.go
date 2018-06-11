package gube

import (
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/ghodss/yaml"
	"github.com/mandelsoft/filepath/pkg/filepath"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type GardenSetConfig interface {
	GetConfig(name string) (GardenConfig, error)
	GetNames() []string
	GetConfigs() map[string]GardenConfig
	GetGithubURL() string
	GetDefault() string
}

type GardenConfig interface {
	GetName() string
	GetDescription() string
	GetGarden() (Garden, error)
	GetRuntimeObject() runtime.Object
	KubeconfigProvider
}

func NewDefaultGardenSetConfig(g Garden) GardenSetConfig {
	cfg := &GardenConfigImpl{
		Name:        "default",
		Description: "default garden",
		garden:      g,
	}

	return &GardenSetConfigImpl{
		Default: "default",
		Gardens: []*GardenConfigImpl{cfg},
	}
}

func NewGardenSetConfig(path string) (GardenSetConfig, error) {
	return ReadGardenSetConfig(path)
}

func ReadGardenSetConfig(path string) (*GardenSetConfigImpl, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := &GardenSetConfigImpl{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	config.SetPath(path)
	return config, nil
}

/////////////////////////////////////////////////////////////////////////////

type GardenSetConfigImpl struct {
	GithubURL string              `yaml:"githubURL,omitempty" json:"githubURL,omitempty"`
	Gardens   []*GardenConfigImpl `yaml:"gardens,omitempty" json:"gardens,omitempty"`
	Default   string              `yaml:"default,omitempty" json:"default,omitempty"`
	path      string
}

func (this *GardenSetConfigImpl) SetPath(path string) {
	this.path = path
	dir := filepath.Dir2(path)
	for _, g := range this.Gardens {
		g.makeAbsolute(dir)
	}
}

func (this *GardenSetConfigImpl) GetDefault() string {
	return this.Default
}

func (this *GardenSetConfigImpl) GetGithubURL() string {
	return this.GithubURL
}

func (this *GardenSetConfigImpl) GetConfig(name string) (GardenConfig, error) {
	if name == "" {
		name = this.Default
	}
	if name == "" {
		return nil, fmt.Errorf("No garden name given")
	}
	for _, conf := range this.Gardens {
		if conf.Name == name {
			return conf, nil
		}
	}
	return nil, fmt.Errorf("Garden '%s' not found", name)
}

func (this *GardenSetConfigImpl) GetNames() []string {
	result := []string{}
	for _, conf := range this.Gardens {
		result = append(result, conf.Name)
	}
	return result
}

func (this *GardenSetConfigImpl) GetConfigs() map[string]GardenConfig {
	result := map[string]GardenConfig{}
	for _, conf := range this.Gardens {
		result[conf.Name] = conf
	}
	return result
}

/////////////////////////////////////////////////////////////////////////////

type GardenConfigImpl struct {
	Name           string `yaml:"name,omitempty" json:"name,omitempty"`
	KubeConfigPath string `yaml:"kubeconfig,omitempty" json:"kubeconfig,omitempty"`
	Description    string `yaml:"description,omitempty" json:"description,omitempty"`
	lock           sync.Mutex
	kubeconfig     []byte
	garden         Garden
	effectivePath  string
}

func (this *GardenConfigImpl) makeAbsolute(dir string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if !filepath.IsAbs(this.KubeConfigPath) {
		this.effectivePath = filepath.Join(dir, this.KubeConfigPath)
	}
}

func (this *GardenConfigImpl) GetName() string {
	return this.Name
}

func (this *GardenConfigImpl) GetDescription() string {
	return this.Description
}

func (this *GardenConfigImpl) GetKubeconfig() ([]byte, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.kubeconfig == nil {
		var cfg []byte

		g, err := this._getGarden()
		if err != nil {
			return nil, err
		}
		cfg, err = g.GetKubeconfig()
		if cfg == nil && this.KubeConfigPath != "" {
			path := this.effectivePath
			if path == "" {
				path = this.KubeConfigPath
			}
			cfg, err = ioutil.ReadFile(path)
		}
		if err != nil {
			return nil, err
		}
		this.kubeconfig = cfg
	}
	return this.kubeconfig, nil
}

func (this *GardenConfigImpl) GetGarden() (Garden, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this._getGarden()
}

func (this *GardenConfigImpl) _getGarden() (Garden, error) {
	if this.garden == nil {
		path := this.effectivePath
		if path == "" {
			path = this.KubeConfigPath
		}
		g, err := NewGardenFromConfigfile(path)
		if err != nil {
			return nil, fmt.Errorf("cannot create garden object for %s(%s)", this.Name, path)
		}
		this.garden = g
	}
	return this.garden, nil
}

type gardenObject struct {
	*GardenConfigImpl `json:",inline"`
}

func (this *gardenObject) GetObjectKind() schema.ObjectKind {
	return nil
}
func (this *gardenObject) DeepCopyObject() runtime.Object {
	return nil
}

func (g *GardenConfigImpl) GetRuntimeObject() runtime.Object {
	return &gardenObject{g}
}
