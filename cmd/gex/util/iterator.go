package util

type Iterator interface {
	HasNext() bool
	Next() interface{}
	Reset()
}

////////////////////////////////////////////////////////////////////////////

type IndexedAccess interface {
	Size() int
	Get(int) interface{}
}

type IndexedIterator struct {
	access  IndexedAccess
	current int
}

func NewIndexedIterator(a IndexedAccess) *IndexedIterator {
	return &IndexedIterator{a, -1}
}

func (this *IndexedIterator) HasNext() bool {
	return this.access.Size() > this.current+1
}

func (this *IndexedIterator) Next() interface{} {
	if this.HasNext() {
		this.current++
		return this.access.Get(this.current)
	}
	return nil
}

func (this *IndexedIterator) Reset() {
	this.current = -1
}

////////////////////////////////////////////////////////////////////////////

type IndexedSliceAccess struct {
	slice []interface{}
}

func NewIndexedSliceAccess(slice []interface{}) IndexedAccess {
	return &IndexedSliceAccess{slice}
}

func (this *IndexedSliceAccess) Size() int {
	return len(this.slice)
}

func (this *IndexedSliceAccess) Get(i int) interface{} {
	return this.slice[i]
}

func NewSliceIterator(slice []interface{}) *IndexedIterator {
	return NewIndexedIterator(NewIndexedSliceAccess(slice))
}

////////////////////////////////////////////////////////////////////////////

type MappedIterator struct {
	mapper func(interface{}) interface{}
	Iterator
}

func NewMappedIterator(i Iterator, mapper func(interface{}) interface{}) *MappedIterator {
	return &MappedIterator{mapper, i}
}

func (this *MappedIterator) Next() interface{} {
	if this.HasNext() {
		return this.mapper(this.Iterator.Next())
	}
	return nil
}
