package data

func Slice(s Iterable) []interface{} {
	a := []interface{}{}
	i := s.Iterator()
	for i.HasNext() {
		a = append(a, i.Next())
	}
	return a
}
