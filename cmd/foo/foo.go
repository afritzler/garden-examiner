package main

import (
	"fmt"

	"github.com/afritzler/garden-examiner/pkg/data"
)

func NewMapping(f int) data.MappingFunction {
	return func(v interface{}) interface{} {
		return v.(int) * f
	}
}

func main() {
	slice := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("%v\n", slice)

	a := data.NewIndexedSliceAccess(slice)
	mapper := data.NewLimitedParallelOrderedMapper(2)
	r := data.Slice(mapper.Map(a, NewMapping(2)))

	mapper.Stop()
	fmt.Printf("%v\n", r)
	data.Stop()
	{
		even := func(e interface{}) bool {
			return e.(int)%2 == 0
		}
		odd := func(e interface{}) bool {
			return e.(int)%2 == 1
		}
		decreasing := func(a interface{}, b interface{}) int {
			return b.(int) - a.(int)
		}
		times2 := NewMapping(2)

		r = data.Process(a).Filter(even).Map(NewMapping(2)).Sort(decreasing).AsSlice()
		fmt.Printf("%v\n", r)

		//r = data.Slice(data.NewContainerFromIterable(a))
		r = data.Process(a).Parallel(5).Filter(even).Map(NewMapping(2)).AsSlice()
		fmt.Printf("ordered  %v\n", r)
		r = data.Process(a).Parallel(5).Unordered().Filter(even).Map(NewMapping(2)).AsSlice()
		fmt.Printf("unordered %v\n", r)
		r = data.Process(a).Parallel(5).Unordered().Filter(even).Map(NewMapping(2)).Sort(decreasing).AsSlice()
		fmt.Printf("sorted %v\n", r)

		base := data.Process(a).Parallel(5)
		unordered := base.Unordered()
		filtered := unordered.Filter(even)

		fmt.Printf("chained: ordered %v\n", base.Filter(odd).Map(times2).AsSlice())
		fmt.Printf("chained: unordered %v\n", filtered.Map(times2).AsSlice())
		fmt.Printf("chained: sorted %v\n", filtered.Map(times2).Sort(decreasing).AsSlice())
	}

}
