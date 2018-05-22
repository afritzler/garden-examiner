package data

type IncrementalProcessingSource interface {
	Iterable
	Open()
	Add(e ...interface{}) IncrementalProcessingSource
	Close()
}

type ProcessingSource interface {
	IncrementalProcessingSource
	IndexedAccess
}

type FilterFunction func(interface{}) bool
type MappingFunction func(interface{}) interface{}
type CompareFunction func(interface{}, interface{}) int

type ProcessingResult interface {
	Iterable

	Map(m MappingFunction) ProcessingResult
	Filter(f FilterFunction) ProcessingResult
	Sort(c CompareFunction) ProcessingResult
	Apply(c ProcessChain) ProcessingResult

	WithPool(ProcessorPool) ProcessingResult
	Unordered() ProcessingResult
	Sequential() ProcessingResult
	Parallel(n int) ProcessingResult

	AsSlice() IndexedSliceAccess
}

func Process(data Iterable) ProcessingResult {
	return (&_SequentialProcessing{}).new(data)
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

/////////////////////////////////////////////////////////////////////////////

type _SequentialProcessing struct {
	data Iterable
}

var _ Iterable = &_SequentialProcessing{}

func (this *_SequentialProcessing) new(data Iterable) *_SequentialProcessing {
	this.data = data
	return this
}

func (this *_SequentialProcessing) Map(m MappingFunction) ProcessingResult {
	return (&_SequentialStep{}).new(this.data, mapper(m))
}
func (this *_SequentialProcessing) Filter(f FilterFunction) ProcessingResult {
	return (&_SequentialStep{}).new(this.data, filter(f))
}
func (this *_SequentialProcessing) Sort(c CompareFunction) ProcessingResult {
	return &_SequentialProcessing{this.AsSlice().Sort(c)}
}
func (this *_SequentialProcessing) WithPool(p ProcessorPool) ProcessingResult {
	return (&_ParallelProcessing{}).new(newEntryIterableFromIterable(this.data), p, NewOrderedContainer)
}
func (this *_SequentialProcessing) Parallel(n int) ProcessingResult {
	return this.WithPool(NewProcessorPool(n))
}
func (this *_SequentialProcessing) Sequential() ProcessingResult {
	return this
}
func (this *_SequentialProcessing) Unordered() ProcessingResult {
	return this
}
func (this *_SequentialProcessing) Apply(c ProcessChain) ProcessingResult {
	return c.Process(this)
}

func (this *_SequentialProcessing) Iterator() Iterator {
	return this.data.Iterator()
}
func (this *_SequentialProcessing) AsSlice() IndexedSliceAccess {
	return IndexedSliceAccess(Slice(this.data))
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
