package util

import (
	"fmt"
)

type AttributeSet struct {
	attrs [][]string
}

func NewAttributeSet() *AttributeSet {
	a := &AttributeSet{}
	a.ResetAttributes()
	return a
}

func (this *AttributeSet) ResetAttributes() {
	this.attrs = [][]string{[]string{}}
}

func (this *AttributeSet) Attribute(name, value string) {
	this.attrs = append(this.attrs, []string{name + ":", value})
}

func (this *AttributeSet) Attributef(name, f string, args ...interface{}) {
	this.attrs = append(this.attrs, []string{name + ":", fmt.Sprintf(f, args...)})
}

func (this *AttributeSet) PrintAttributes() {
	FormatTable("", this.attrs)
}
