package gube

import (
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/ghodss/yaml"
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
	GetKubeconfig() ([]byte, error)
	GetGarden() (Garden, error)
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
	return config, nil
}

/////////////////////////////////////////////////////////////////////////////

type GardenSetConfigImpl struct {
	GithubURL string              `yaml:"githubURL,omitempty" json:"githubURL,omitempty"`
	Gardens   []*GardenConfigImpl `yaml:"gardens,omitempty" json:"gardens,omitempty"`
	Default   string              `yaml:"default,omitempty" json:"default,omitempty"`
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
		var err error

		if this.garden == nil {
			g, err := NewGardenFromConfigfile(this.KubeConfigPath)
			if err != nil {
				return nil, fmt.Errorf("cannot create garden object for %s", this.Name)
			}
			this.garden = g
		}
		cfg, err = this.garden.GetKubeconfig()
		if cfg == nil && this.KubeConfigPath != "" {
			cfg, err = ioutil.ReadFile(this.KubeConfigPath)
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
	if this.garden == nil {
		g, err := NewGardenFromConfigfile(this.KubeConfigPath)
		if err != nil {
			return nil, fmt.Errorf("cannot create garden object for %s", this.Name)
		}
		this.garden = g
	}
	return this.garden, nil
}
