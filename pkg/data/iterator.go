package data

type Iterable interface {
	Iterator() Iterator
}

type Iterator interface {
	HasNext() bool
	Next() interface{}
}

type ResettableIterator interface {
	HasNext() bool
	Next() interface{}
	Reset()
}

type MappedIterator struct {
	Iterator
	mapping MappingFunction
}

func NewMappedIterator(iter Iterator, mapping MappingFunction) Iterator {
	return &MappedIterator{iter, mapping}
}

func (this *MappedIterator) Next() interface{} {
	if this.HasNext() {
		return this.mapping(this.Iterator.Next())
	}
	return nil
}
