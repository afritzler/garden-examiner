package util

type Map interface {
	Has(interface{}) bool
	Get(interface{}) interface{}
	Set(interface{}, interface{})
	Iterator() Iterator
	Keys() Iterator
	Values() Iterator
	Size() int
}

type MapEntry struct {
	Key   interface{}
	Value interface{}
}

type _Map map[interface{}]interface{}

func NewMap() Map {
	return _Map{}
}

func (this _Map) Has(key interface{}) bool {
	_, ok := this[key]
	return ok
}

func (this _Map) Get(key interface{}) interface{} {
	v, ok := this[key]
	if ok {
		return v
	}
	return nil
}

func (this _Map) Set(key interface{}, value interface{}) {
	this[key] = value
}

func (this _Map) Size() int {
	return len(this)
}

func (this _Map) Keys() Iterator {
	return &_MapIterator{this, newMapKeyIterator(this)}
}

func (this _Map) Iterator() Iterator {
	return &_EntryIterator{this, newMapKeyIterator(this)}
}

func (this _Map) Values() Iterator {
	return NewMappedIterator(this.Iterator(), func(e interface{}) interface{} {
		return e.(MapEntry).Value
	})
}

type _MapIterator struct {
	data _Map
	Iterator
}

func newMapKeyIterator(m _Map) Iterator {
	keys := make([]interface{}, m.Size())
	i := 0
	for k, _ := range m {
		keys[i] = k
		i++
	}
	return NewSliceIterator(keys)
}

func (this *_MapIterator) Reset() {
	this.Iterator = newMapKeyIterator(this.data)
}

type _EntryIterator _MapIterator

func (this *_EntryIterator) Reset() {
	this.Iterator.Reset()
}

func (this *_EntryIterator) Next() interface{} {
	if this.HasNext() {
		k := this.Iterator.Next()
		v, _ := this.data[k]
		return MapEntry{k, v}
	}
	return nil
}
