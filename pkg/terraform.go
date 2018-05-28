package gube

import (
	"encoding/json"
)

type terraformModule struct {
	Outputs map[string]*terraformOutput `json:"outputs"`
	Path    []string                    `json:"path"`
}

type terraformOutput struct {
	Sensitive bool        `json:"sensitive"`
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
}

type TerraformState struct {
	Modules []*terraformModule `json:"modules"`
}

func NewTerraformStateFromConfig(data map[string]string) (*TerraformState, error) {
	return NewTerraformState([]byte(data["terraform.tfstate"]))
}

func NewTerraformState(data []byte) (*TerraformState, error) {
	state := &TerraformState{}

	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return state, nil
}

func (this *TerraformState) GetOutput(name string) interface{} {
	o, ok := this.Modules[0].Outputs[name]
	if ok {
		return o.Value
	}
	return nil
}

func (this *TerraformState) GetOutputs() map[string]interface{} {
	out := map[string]interface{}{}
	for k, o := range this.Modules[0].Outputs {
		out[k] = o
	}
	return out
}
