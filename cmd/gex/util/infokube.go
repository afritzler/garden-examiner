package util

import (
	"fmt"
	"sort"
	"strconv"
)

/////////////////////////////////////////////////////////////////////////////
// helper

type compare_function func(a, b string) int

type dimension_names struct {
	names   []string
	indices map[string]int
}

func NewNames(names []string) *dimension_names {
	d := &dimension_names{[]string{}, map[string]int{}}
	for i, n := range names {
		d.indices[n] = i
	}
	return d
}

func (a *dimension_names) Len() int {
	return len(a.names)
}
func (a *dimension_names) Swap(i, j int) {
	a.names[i], a.names[j] = a.names[j], a.names[i]
}
func (a *dimension_names) Less(i, j int) bool {
	return a.indices[a.names[i]] < a.indices[a.names[j]]
}

func (a *dimension_names) append(d ...string) *dimension_names {
	a.names = append(a.names, d...)
	return a
}

const (
	NONE = "NONE"
	COL  = "COL"
	ROW  = "ROW"
)

/////////////////////////////////////////////////////////////////////////////
// InfoKube

type Coord map[string]string

func (this Coord) Clone() Coord {
	clone := Coord{}
	for k, v := range this {
		clone[k] = v
	}
	return clone
}

type InfoKube struct {
	values     map[string]map[string]*InfoKube
	count      int
	dimensions []string
	elems      []interface{}
}

func NewInfoKube(dim []string) *InfoKube {
	return &InfoKube{map[string]map[string]*InfoKube{}, 0, dim, []interface{}{}}
}

func (this *InfoKube) GetCount() int {
	return this.count
}

func (this *InfoKube) GetKeys(dim string) map[string]*InfoKube {
	kube, ok := this.values[dim]
	if ok {
		return kube
	}
	return map[string]*InfoKube{}
}

func (this *InfoKube) GetKey(dim string, key string) *InfoKube {
	keys, ok := this.values[dim]
	if ok {
		kube, ok := keys[key]
		if ok {
			return kube
		}
	}
	if stringIndex(dim, this.dimensions) >= 0 {
		return NewInfoKube(this.dimensions[stringIndex(dim, this.dimensions)+1:])
	} else {
		return nil
	}
}

func (this *InfoKube) AddKey(dim, value string) {
	if stringIndex(dim, this.dimensions) >= 0 {
		for i, d := range this.dimensions {
			for _, kube := range this.values[d] {
				kube.AddKey(dim, value)
			}
			if d == dim {
				this.addKey(i, dim, value)
				break
			}
		}
	}
}

func (this *InfoKube) addKey(i int, key, value string) *InfoKube {
	dim, ok := this.values[key]
	if !ok {
		dim = map[string]*InfoKube{}
		this.values[key] = dim
	}
	kube, ok := dim[value]
	if !ok {
		kube = NewInfoKube(this.dimensions[i+1:])
		dim[value] = kube
	}
	return kube
}

func (this *InfoKube) AddElement(elem interface{}, keys ...string) {
	if len(keys) != len(this.dimensions) {
		panic(fmt.Errorf("expected %d coordinates, found %d\n", len(this.dimensions), len(keys)))
	}
	for i, k := range keys {
		this.addKey(i, this.dimensions[i], k).AddElement(elem, keys[i+1:]...)
	}
	this.count++
	if elem != nil {
		this.elems = append(this.elems, elem)
	}
}

func (this *InfoKube) Table(gap string, dims []string, keys Coord) {
	names := NewNames(this.dimensions)
	names.append(dims...)
	for d, _ := range keys {
		names.append(d)
	}
	sort.Sort(names)
	this.table(gap, dims, keys, names)
}

func (this *InfoKube) table(gap string, dims []string, keys Coord, names *dimension_names) {
	switch len(dims) {
	case 0:
		return
	case 1:
		this.table1(nil, gap, dims[0], keys, names)
	case 2:
		this.table2(nil, gap, dims[0], dims[1], keys, names)
	default:
		sub := []Coord{}
		if len(dims)%2 == 1 {
			sub = this.table1(sub, gap, dims[0], keys, names)
			dims = dims[1:]
		} else {
			sub = this.table2(sub, gap, dims[0], dims[1], keys, names)
			dims = dims[2:]
		}
		gap = gap + "-> "
		for _, c := range sub {
			fmt.Println()
			subkeys := keys.Clone()
			sep := gap
			for _, d := range this.dimensions {
				if c[d] != "" {
					fmt.Printf("%s%s=%s", sep, d, c[d])
					sep = ", "
					subkeys[d] = c[d]
				}
			}
			fmt.Println()
			this.table(gap, dims, subkeys, names)
		}
	}
	return
}

func (this *InfoKube) Table1(gap, dx string, keys Coord) {
	names := NewNames(this.dimensions)
	names.append(dx)
	for d, _ := range keys {
		names.append(d)
	}
	sort.Sort(names)
	this.table1(nil, gap, dx, keys, names)
}

func (this *InfoKube) table1(list []Coord, gap, dx string, keys Coord, names *dimension_names) []Coord {
	kube := this
	initial := []string{}
	var end []string

	for _, d := range names.names {
		switch {
		case d == dx:
			end = []string{}
		case keys[d] != "":
			if end == nil {
				initial = append(initial, d)
				kube = kube.values[d][keys[d]]
				if kube == nil {
					fmt.Printf("%sno entry\n", gap)
					return list
				}
			} else {
				end = append(end, d)
			}
		}
	}

	line := []string{}
	head := []string{}

	sub, ok := kube.values[dx]
	if !ok {
		fmt.Printf("%sno entry\n", gap)
	} else {
		for kx, kubex := range sub {
			head = append(head, "-"+kx)
			kubex = kubex.proceed(end, keys)
			if kubex != nil {
				line = append(line, strconv.Itoa(kubex.count))
				if kubex.count > 0 && list != nil {
					list = append(list, map[string]string{dx: kx})
				}
			} else {
				line = append(line, "0")
			}
		}
	}

	FormatTable(gap, [][]string{head, line})
	return list
}

func (this *InfoKube) Table2(gap, dy, dx string, keys Coord) {
	names := NewNames(this.dimensions)
	names.append(dx, dy)
	for d, _ := range keys {
		names.append(d)
	}
	sort.Sort(names)
	//fmt.Printf("sorted: %v\n", names.names)
	this.table2(nil, gap, dy, dx, keys, names)
}

func (this *InfoKube) table2(list []Coord, gap, dy, dx string, keys Coord, names *dimension_names) []Coord {
	kube := this
	mode := NONE
	initial := []string{}
	var middle, end []string

	for _, d := range names.names {
		switch {
		case d == dy:
			if mode == NONE {
				mode = ROW
				middle = []string{}
			} else {
				end = []string{}
			}
		case d == dx:
			if mode == NONE {
				mode = COL
				middle = []string{}
			} else {
				end = []string{}
			}
		case keys[d] != "":
			if mode == NONE {
				//fmt.Printf("down %s for %s\n", keys[d], d)
				initial = append(initial, d)
				kube = kube.values[d][keys[d]]
				if kube == nil {
					fmt.Printf("%sno entry\n", gap)
					return list
				}
			} else {
				if end != nil {
					end = append(end, d)
				} else if middle != nil {
					middle = append(middle, d)
				}
			}
		}
	}

	//fmt.Printf("mode: %v\n", mode)
	//fmt.Printf("initial: %v\n", initial)
	//fmt.Printf("middle:  %v\n", middle)
	//fmt.Printf("end:     %v\n", end)

	if mode == COL {
		return kube.col_table2(list, gap, dy, dx, middle, end, keys)
	} else {
		return kube.row_table2(list, gap, dy, dx, middle, end, keys)
	}
}

func (this *InfoKube) col_table2(list []Coord, gap, dy, dx string, middle, end []string, keys Coord) []Coord {
	indices := map[string]int{}
	data := [][]string{}
	head := []string{dy + "\\" + dx}

	values, ok := this.values[dx]
	if !ok {
		fmt.Printf("%sno entry\n", gap)
	} else {
		for kx, kubex := range values {
			kubex := kubex.proceed(middle, keys)
			if kubex != nil {
				sub, ok := kubex.values[dy]
				if ok {
					head = append(head, "-"+kx)
					for ky, kubey := range sub {
						i, ok := indices[ky]
						if !ok {
							i = len(data)
							indices[ky] = i
							data = append(data, []string{ky})
						}
						line := data[i]
						for len(head) > len(line) {
							line = append(line, "0")
							data[i] = line
						}
						kubey = kubey.proceed(end, keys)
						if kubey != nil {
							line[len(head)-1] = strconv.Itoa(kubey.count)
							if kubey.count > 0 && list != nil {
								list = append(list, map[string]string{dx: kx, dy: ky})
							}
						}
					}
				}
			}
		}
		for k, v := range data {
			for len(v) < len(head) {
				v = append(v, "0")
			}
			data[k] = v
		}
		FormatTable(gap, append([][]string{head}, data...))
	}
	return list
}

func (this *InfoKube) row_table2(list []Coord, gap, dy, dx string, middle, end []string, keys Coord) []Coord {
	indices := map[string]int{}
	data := [][]string{}
	head := []string{dy + "\\" + dx}

	values, ok := this.values[dy]
	if !ok {
		fmt.Printf("%sno entry\n", gap)
	} else {
		for ky, kubey := range values {
			line := []string{ky}
			kubey := kubey.proceed(middle, keys)
			if kubey != nil {
				sub, ok := kubey.values[dx]
				if ok {
					for kx, kubex := range sub {
						i, ok := indices[kx]
						if !ok {
							i = len(head)
							indices[kx] = i
							head = append(head, "-"+kx)
						}
						for i >= len(line) {
							line = append(line, "0")
						}
						kubex = kubex.proceed(end, keys)
						if kubex != nil {
							line[i] = strconv.Itoa(kubex.count)
							if kubex.count > 0 && list != nil {
								list = append(list, map[string]string{dx: kx, dy: ky})
							}
						}
					}
					data = append(data, line)
				}
			}
		}
		for k, v := range data {
			for len(v) < len(head) {
				v = append(v, "0")
			}
			data[k] = v
		}
		FormatTable(gap, append([][]string{head}, data...))
	}
	return list
}

func (this *InfoKube) proceed(dims []string, keys Coord) *InfoKube {
	sub := this
	ok := true
	for _, d := range dims {
		sub, ok = sub.values[d][keys[d]]
		if !ok {
			return nil
		}
	}
	return sub
}

func stringIndex(s string, a []string) int {
	for i, v := range a {
		if v == s {
			return i
		}
	}
	return -1
}
