package data

import (
	"fmt"
	"sync"
)

var log = false

type _ParallelProcessing struct {
	data    entry_iterable
	pool    ProcessorPool
	creator container_creator
}

var _ Iterable = &_ParallelProcessing{}

func (this *_ParallelProcessing) new(data entry_iterable, pool ProcessorPool, creator container_creator) *_ParallelProcessing {
	this.data = data
	this.pool = pool
	this.creator = creator
	return this
}

func (this *_ParallelProcessing) Map(m MappingFunction) *_ParallelStep {
	return (&_ParallelStep{}).new(this.pool, this.data, mapper(m), this.creator)
}
func (this *_ParallelProcessing) Filter(f FilterFunction) *_ParallelStep {
	return (&_ParallelStep{}).new(this.pool, this.data, filter(f), this.creator)
}
func (this *_ParallelProcessing) Sort(c CompareFunction) *_ParallelProcessing {
	return (&_ParallelProcessing{}).new(this.AsSlice().Sort(c), this.pool, NewOrderedContainer)
}

func (this *_ParallelProcessing) WithPool(p ProcessorPool) *_ParallelProcessing {
	return (&_ParallelProcessing{}).new(this.data, p, this.creator)
}
func (this *_ParallelProcessing) Parallel(n int) *_ParallelProcessing {
	return this.WithPool(NewProcessorPool(n))
}
func (this *_ParallelProcessing) Sequential() *_SequentialProcessing {
	return (&_SequentialProcessing{}).new(this.data)
}
func (this *_ParallelProcessing) Unordered() *_ParallelProcessing {
	data := this.data
	ordered, ok := data.(*ordered_container)
	if ok {
		data = &ordered._container
	}
	return (&_ParallelProcessing{}).new(data, this.pool, NewContainer)
}

func (this *_ParallelProcessing) Iterator() Iterator {
	return this.data.Iterator()
}
func (this *_ParallelProcessing) AsSlice() IndexedSliceAccess {
	return IndexedSliceAccess(Slice(this.data))
}

////////////////////////////////////////////////////////////////////////////

type _ParallelStep struct {
	_ParallelProcessing
	container container
	op        operation
	create    container_creator
}

func (this *_ParallelStep) new(pool ProcessorPool, data entry_iterable, op operation, creator container_creator) *_ParallelStep {
	this.container = creator()
	this._ParallelProcessing.new(this.container, pool, creator)
	go func() {
		if log {
			fmt.Printf("start processing\n")
		}
		this.pool.Request()
		i := data.entry_iterator()
		var wg sync.WaitGroup
		for i.HasNext() {
			e := i.next()
			if log {
				fmt.Printf("start %d\n", e.index)
			}
			wg.Add(1)
			pool.Exec(func() {
				if log {
					fmt.Printf("process %d\n", e.index)
				}
				e.entry, e.ok = op.process(e.entry)
				this.container.add(e)
				if log {
					fmt.Printf("done %d\n", e.index)
				}
				wg.Done()

			})
		}
		wg.Wait()
		this.pool.Release()
		this.container.close()
		if log {
			fmt.Printf("done processing\n")
		}
	}()
	return this
}
