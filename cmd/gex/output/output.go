package output

import (
	"fmt"
	"os"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/pkg"
	"github.com/gardener/gardener/pkg/client/garden/clientset/versioned/scheme"
	"github.com/mandelsoft/cmdint/pkg/cmdint"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

type Output interface {
	Add(ctx *context.Context, e interface{}) error
	Close(ctx *context.Context) error
	Out(*context.Context) error
}

type ManifestOutput struct {
	data []runtime.Object
}

type YAMLOutput struct {
	ManifestOutput
}

type JSONOutput struct {
	ManifestOutput
	pretty bool
}

func (this *ManifestOutput) Add(ctx *context.Context, e interface{}) error {
	this.data = append(this.data, e.(gube.RuntimeObjectWrapper).GetRuntimeObject())
	return nil
}

func (this *ManifestOutput) Close(ctx *context.Context) error {
	return nil
}

var typer = runtime.MultiObjectTyper{scheme.Scheme}

func (this *YAMLOutput) Out(ctx *context.Context) error {
	ser := json.NewYAMLSerializer(json.DefaultMetaFactory, nil, typer)
	for _, m := range this.data {
		fmt.Println("---")
		err := ser.Encode(m, os.Stdout)
		if err != nil {

		}
	}
	return nil
}

func (this *JSONOutput) Out(*context.Context) error {
	ser := json.NewSerializer(json.DefaultMetaFactory, nil, typer, this.pretty)
	for _, m := range this.data {
		err := ser.Encode(m, os.Stdout)
		if err != nil {

		}
		if this.pretty {
			fmt.Println()
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////

type OutputFactory func(*cmdint.Options) Output

type Outputs map[string]OutputFactory

func NewOutputs(def OutputFactory, others ...Outputs) Outputs {
	o := Outputs{"": def}
	for _, other := range others {
		for k, v := range other {
			o[k] = v
		}
	}
	return o
}

func (this Outputs) Select(name string) OutputFactory {
	c, ok := this[name]
	if !ok {
		keys := []string{}
		for k, _ := range this {
			keys = append(keys, k)
		}
		k, _ := cmdint.SelectBest(name, keys...)
		if k != "" {
			c = this[k]
		}
	}
	return c
}

func (this Outputs) Create(opts *cmdint.Options) (Output, error) {
	f := opts.GetOptionValue(constants.O_OUTPUT)
	if f == nil {
		return this[""](opts), nil
	}
	c := this.Select(*f)
	if c != nil {
		o := c(opts)
		if o != nil {
			return o, nil
		}
	}
	return nil, fmt.Errorf("invalid output format '%s'", *f)
}

func (this Outputs) AddManifestOutputs() Outputs {
	this["yaml"] = func(opts *cmdint.Options) Output {
		return &YAMLOutput{ManifestOutput{data: []runtime.Object{}}}
	}
	this["json"] = func(opts *cmdint.Options) Output {
		return &JSONOutput{ManifestOutput{data: []runtime.Object{}}, true}
	}
	this["JSON"] = func(opts *cmdint.Options) Output {
		return &JSONOutput{ManifestOutput{data: []runtime.Object{}}, false}
	}
	return this
}

func GetOutput(opts *cmdint.Options, def Output) (Output, error) {
	o := def
	f := opts.GetOptionValue(constants.O_OUTPUT)
	if f != nil {
		switch *f {
		case "yaml":
			o = &YAMLOutput{ManifestOutput{data: []runtime.Object{}}}
		case "json":
			o = &JSONOutput{ManifestOutput{data: []runtime.Object{}}, true}
		case "JSON":
			o = &JSONOutput{ManifestOutput{data: []runtime.Object{}}, false}
		default:
			return nil, fmt.Errorf("invalid output format '%s'", *f)
		}
	}
	return o, nil
}
