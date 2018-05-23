package main

import (
	"fmt"
)

type Node interface {
	Map(s string) Node
	Process(string) string
}

type _Node struct {
	parent Node
	data   string
}

var _ Node = &_Node{}

func NewNode() Node {
	return (&_Node{}).new(nil, "")
}

func (this *_Node) new(p *_Node, s string) *_Node {

	fmt.Printf("Parent: %+v\n", Node(p))
	if p == nil {
		this.parent = nil
	} else {
		this.parent = p
	}
	//this.parent = p
	this.data = s
	return this
}

func (this *_Node) Map(s string) Node {
	return (&_Node{}).new(this, s)
}

func (this *_Node) Process(data string) string {
	fmt.Printf("THIS: %+v\n", this)
	if this.parent == nil {
		fmt.Printf("parent :NIL\n")
		return data
	}
	fmt.Printf("recursion %p\n", this.parent)
	return this.parent.Process(data) + "->" + this.data
}

func Test() {
	n := NewNode().Map("a").Map("b")
	fmt.Printf("%s\n", n.Process("start"))
}
