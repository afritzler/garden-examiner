package gube

import (
	"fmt"
)

type GardenSetConfig interface {
	GetConfig(name string) (GardenConfig, error)
	GetNames() []string
	GetConfigs() map[string]GardenConfig
}

type GardenConfig interface {
	GetKubeconfig() ([]byte, error)
	GetName() string
	GetGarden() (Garden, error)
}

func NewGardenSetConfig(path string) (GardenSetConfig, error) {
	return nil, fmt.Errorf("To implement")
}
