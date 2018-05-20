package data

import (
	"sort"
)

type elements struct {
	data    []interface{}
	compare CompareFunction
}

func (a elements) Len() int           { return len(a.data) }
func (a elements) Swap(i, j int)      { a.data[i], a.data[j] = a.data[j], a.data[i] }
func (a elements) Less(i, j int) bool { return a.compare(a.data[i], a.data[j]) < 0 }

func Sort(data []interface{}, cmp CompareFunction) {
	sort.Sort(&elements{data, cmp})
}
