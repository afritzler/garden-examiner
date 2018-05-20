package data

type FilterFunction func(interface{}) bool
type MappingFunction func(interface{}) interface{}
type CompareFunction func(interface{}, interface{}) int

func Process(data Iterable) *_SequentialProcessing {
	return (&_SequentialProcessing{}).new(data)
}

/////////////////////////////////////////////////////////////////////////////

type _SequentialProcessing struct {
	data Iterable
}

var _ Iterable = &_SequentialProcessing{}

func (this *_SequentialProcessing) new(data Iterable) *_SequentialProcessing {
	this.data = data
	return this
}

func (this *_SequentialProcessing) Map(m MappingFunction) *_SequentialStep {
	return (&_SequentialStep{}).new(this.data, mapper(m))
}
func (this *_SequentialProcessing) Filter(f FilterFunction) *_SequentialStep {
	return (&_SequentialStep{}).new(this.data, filter(f))
}
func (this *_SequentialProcessing) Sort(c CompareFunction) *_SequentialProcessing {
	return &_SequentialProcessing{this.AsSlice().Sort(c)}
}

func (this *_SequentialProcessing) Parallel(n int) *_ParallelProcessing {
	return (&_ParallelProcessing{}).new(newEntryIterableFromIterable(this.data), NewProcessorPool(n), NewOrderedContainer)
}

func (this *_SequentialProcessing) Iterator() Iterator {
	return this.data.Iterator()
}
func (this *_SequentialProcessing) AsSlice() IndexedSliceAccess {
	return IndexedSliceAccess(Slice(this.data))
}

////////////////////////////////////////////////////////////////////////////

type operation interface {
	process(e interface{}) (interface{}, bool)
}

type mapper MappingFunction

func (this mapper) process(e interface{}) (interface{}, bool) {
	return this(e), true
}

type filter FilterFunction

func (this filter) process(e interface{}) (interface{}, bool) {
	if this(e) {
		return e, true
	}
	return nil, false
}

////////////////////////////////////////////////////////////////////////////

type _SequentialStep struct {
	_SequentialProcessing
	op operation
}

func (this *_SequentialStep) new(data Iterable, op operation) *_SequentialStep {
	slice := []interface{}{}
	i := data.Iterator()
	for i.HasNext() {
		e, ok := op.process(i.Next())
		if ok {
			slice = append(slice, e)
		}
	}
	this.data = NewIndexedSliceAccess(slice)
	return this
}
