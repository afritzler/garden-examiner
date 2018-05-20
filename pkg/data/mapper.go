package data

import (
	"sync"
)

////////////////////////////////////////////////////////////////////////////

type SequentialMapper struct {
}

func NewSequentialMapper() *SequentialMapper {
	return &SequentialMapper{}
}

func (this *SequentialMapper) Map(s Iterable, m MappingFunction) *SequentialMapping {
	return (&SequentialMapping{}).new(s, m)
}

type SequentialMapping struct {
	data []interface{}
}

func (this *SequentialMapping) new(s Iterable, m MappingFunction) *SequentialMapping {
	this.data = []interface{}{}
	i := s.Iterator()
	for i.HasNext() {
		this.data = append(this.data, m(i.Next()))
	}
	return this
}

func (this *SequentialMapping) Map(m MappingFunction) *SequentialMapping {
	return (&SequentialMapping{}).new(NewIndexedSliceAccess(this.data), m)
}

func (this *SequentialMapping) Iterator() Iterator {
	return NewSliceIterator(this.data)
}

////////////////////////////////////////////////////////////////////////////

type ParallelOrderedMapper struct {
}

func NewParallelOrderedMapper() *ParallelOrderedMapper {
	return &ParallelOrderedMapper{}
}

func (this *ParallelOrderedMapper) Map(s Iterable, m MappingFunction) *ParallelOrderedMapping {
	return (&ParallelOrderedMapping{}).new(s, m)
}

type ParallelOrderedMapping struct {
	SequentialMapping
}

type entry struct {
	index int
	ok    bool
	entry interface{}
}

func (this *ParallelOrderedMapping) new(s Iterable, m MappingFunction) *ParallelOrderedMapping {
	this.data = []interface{}{}

	mc := make(chan entry)
	go func() {
		i := s.Iterator()
		idx := 0
		wg := sync.WaitGroup{}
		for i.HasNext() {
			e := i.Next()
			wg.Add(1)
			go func(idx int) {
				r := m(e)
				mc <- entry{idx, true, r}
				wg.Done()
			}(idx)
			idx++
		}
		wg.Wait()
		close(mc)
	}()
	for e := range mc {
		if len(this.data) <= e.index {
			t := make([]interface{}, e.index+1)
			copy(t, this.data)
			this.data = t
		}
		this.data[e.index] = e.entry
	}
	return this
}

////////////////////////////////////////////////////////////////////////////

type Stoppable interface {
	Stop()
}

type request struct {
	m     *LimitedParallelOrderedMapping
	entry entry
}

type LimitedParallelOrderedMapper struct {
	limit    int
	requests chan request
	closed   bool
}

func NewLimitedParallelOrderedMapper(limit int) *LimitedParallelOrderedMapper {
	return (&LimitedParallelOrderedMapper{}).new(limit)
}

func (this *LimitedParallelOrderedMapper) new(limit int) *LimitedParallelOrderedMapper {
	this.limit = limit
	this.requests = make(chan request, limit*4)
	for i := 0; i < limit; i++ {
		go func() {
			for r := range this.requests {
				r.entry.entry = r.m.mapping(r.entry.entry)
				r.m.mc <- r.entry
				r.m.wg.Done()
			}
		}()
	}
	mappers = append(mappers, this)
	return this
}

func (this *LimitedParallelOrderedMapper) Stop() {
	if !this.closed {
		close(this.requests)
		this.closed = true
	}
}

func (this *LimitedParallelOrderedMapper) Map(s Iterable, m MappingFunction) *LimitedParallelOrderedMapping {
	return (&LimitedParallelOrderedMapping{}).new(this, s, m)
}

type LimitedParallelOrderedMapping struct {
	SequentialMapping
	mapper  *LimitedParallelOrderedMapper
	mapping MappingFunction
	mc      chan entry
	wg      sync.WaitGroup
}

func (this *LimitedParallelOrderedMapping) new(mapper *LimitedParallelOrderedMapper, s Iterable, m MappingFunction) *LimitedParallelOrderedMapping {
	this.mapper = mapper
	this.mapping = m
	this.data = []interface{}{}
	this.mc = make(chan entry)

	go func() {
		i := s.Iterator()
		idx := 0
		for i.HasNext() {
			e := i.Next()
			this.wg.Add(1)
			mapper.requests <- request{this, entry{idx, true, e}}
			idx++
		}
		this.wg.Wait()
		close(this.mc)
	}()
	for e := range this.mc {
		if len(this.data) <= e.index {
			t := make([]interface{}, e.index+1)
			copy(t, this.data)
			this.data = t
		}
		this.data[e.index] = e.entry
	}
	return this
}

var mappers []Stoppable = []Stoppable{}

func Stop() {
	for _, s := range mappers {
		s.Stop()
	}
}
