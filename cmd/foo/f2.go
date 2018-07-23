package main

import (
	"fmt"
	"reflect"
)

type Container struct {
	Value           string
	ArrayValue      []string
	MultiValue      []string
	MultiArrayValue [][]string
}

func f2() {
	t := reflect.TypeOf(Container{})
	v := reflect.New(t)
	c := v.Interface().(*Container)

	v = v.Elem()
	v_f := v.FieldByName("Value")
	v_f.SetString("test")

	v_f = v.FieldByName("ArrayValue")
	v_f.Set(reflect.ValueOf([]string{"foo", "bar"}))

	v_f = v.FieldByName("MultiArrayValue")
	reflect.Append(v_f, reflect.ValueOf([]string{"foo", "bar"}))

	v_f = v.FieldByName("MultiValue")
	Append(v_f, "first")

	v_f = v.FieldByName("MultiArrayValue")
	Append(v_f, []string{"alice", "bob"})

	fmt.Printf("%#v\n", c)

}

func Append(a reflect.Value, v interface{}) {
	if a.IsNil() {
		na := reflect.MakeSlice(a.Type(), 0, 1)
		a.Set(na)
	}
	a.Set(reflect.Append(a, reflect.ValueOf(v)))
}
