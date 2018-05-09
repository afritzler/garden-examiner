package util

import (
	"fmt"
	"os"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/pkg"
	"github.com/mandelsoft/cmdint/pkg/cmdint"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

type Output interface {
	Add(ctx *context.Context, e interface{}) error
	Out(*context.Context)
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

func (this *ManifestOutput) Add(ctx *context.Context, e interface{}) error {
	this.data = append(this.data, e.(gube.RuntimeObjectWrapper).GetRuntimeObject())
	return nil
}

func (this *YAMLOutput) Out(ctx *context.Context) {
	ser := json.NewYAMLSerializer(json.DefaultMetaFactory, nil, nil)
	for _, m := range this.data {
		fmt.Println("---")
		err := ser.Encode(m, os.Stdout)
		if err != nil {

		}
	}
}

func (this *JSONOutput) Out(*context.Context) {
	ser := json.NewSerializer(json.DefaultMetaFactory, nil, nil, this.pretty)
	for _, m := range this.data {
		err := ser.Encode(m, os.Stdout)
		if err != nil {

		}
		if this.pretty {
			fmt.Println()
		}
	}
}
