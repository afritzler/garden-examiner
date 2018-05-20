package data

type IndexedAccess interface {
	Size() int
	Get(int) interface{}
}

type IndexedIterator struct {
	access  IndexedAccess
	current int
}

var _ ResettableIterator = &IndexedIterator{}

func NewIndexedIterator(a IndexedAccess) *IndexedIterator {
	return (&IndexedIterator{}).new(a)
}

func (this *IndexedIterator) new(a IndexedAccess) *IndexedIterator {
	this.access = a
	this.current = -1
	return this
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

type IndexedSliceAccess []interface{}

var _ IndexedAccess = IndexedSliceAccess{}
var _ Iterable = IndexedSliceAccess{}

func NewIndexedSliceAccess(slice []interface{}) IndexedSliceAccess {
	return IndexedSliceAccess(slice)
}

func (this IndexedSliceAccess) Size() int {
	return len(this)
}

func (this IndexedSliceAccess) Get(i int) interface{} {
	return this[i]
}

func (this IndexedSliceAccess) Iterator() Iterator {
	return NewIndexedIterator(this)
}

func (this IndexedSliceAccess) Sort(cmp CompareFunction) IndexedSliceAccess {
	Sort(this, cmp)
	return this
}

func (this IndexedSliceAccess) entry_iterator() entry_iterator {
	return (&_slice_entry_iterator{}).new(this)
}

func NewSliceIterator(slice []interface{}) *IndexedIterator {
	return NewIndexedIterator(IndexedSliceAccess(slice))
}

type _slice_entry_iterator struct {
	IndexedIterator
}

func (this *_slice_entry_iterator) new(a IndexedSliceAccess) *_slice_entry_iterator {
	this.IndexedIterator.new(a)
	return this
}

func (this *_slice_entry_iterator) next() entry {
	e := this.Next()
	return entry{this.current, true, e}
}
