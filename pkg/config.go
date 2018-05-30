package gube

import (
	"fmt"
	"io/ioutil"

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
	GetKubeconfig() ([]byte, error)
	GetName() string
	GetGarden() (Garden, error)
}

func NewGardenSetConfig(path string) (GardenSetConfig, error) {
	return ReadGardenSetConfig(path)
}

type GardenSetConfigImpl struct {
	GithubURL string              `yaml:"githubURL,omitempty" json:"githubURL,omitempty"`
	Gardens   []*GardenConfigImpl `yaml:"gardens,omitempty" json:"gardens,omitempty"`
	Default   string              `yaml:"default,omitempty" json:"default,omitempty"`
}
type GardenConfigImpl struct {
	Name        string `yaml:"name,omitempty" json:"name,omitempty"`
	KubeConfig  string `yaml:"kubeconfig,omitempty" json:"kubeconfig,omitempty"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
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

func (this *GardenSetConfigImpl) GetGithubURL() string {
	return this.GithubURL
}

func (this *GardenSetConfigImpl) GetDefault() string {
	return this.Default
}

func (this *GardenConfigImpl) GetKubeconfig() ([]byte, error) {
	return ioutil.ReadFile(this.KubeConfig)
}

func (this *GardenConfigImpl) GetName() string {
	return this.Name
}

func (this *GardenConfigImpl) GetGarden() (Garden, error) {
	return NewGardenFromConfigfile(this.KubeConfig)
}
