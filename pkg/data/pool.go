package data

import (
	"sync"
)

type ProcessorPool interface {
	Request()
	Release()
	Exec(func())
}

/////////////////////////////////////////////////////////////////////////////

type _UnlimitedPool struct {
}

func NewUnlimitedProcessorPool() ProcessorPool {
	return &_UnlimitedPool{}
}

func (this *_UnlimitedPool) Request() {
}
func (this *_UnlimitedPool) Release() {
}
func (this *_UnlimitedPool) Exec(f func()) {
	go f()
}

/////////////////////////////////////////////////////////////////////////////

type _ProcessorPool struct {
	n    int
	uses int
	lock sync.Mutex
	set  *processor_set
}

func NewProcessorPool(n int) ProcessorPool {
	return &_ProcessorPool{n: n, uses: 0}
}

func (this *_ProcessorPool) Request() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.uses++
	if this.uses == 1 {
		this.set = (&processor_set{}).new(this.n)
	}
}

func (this *_ProcessorPool) Exec(f func()) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.uses == 0 {
		panic("unrequested processor pool used")
	}
	this.set.exec(f)
}

func (this *_ProcessorPool) Release() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.uses > 0 {
		this.uses--
		if this.uses <= 0 && this.set != nil {
			this.set.stop()
			this.set = nil
		}
	}
}

/////////////////////////////////////////////////////////////////////////////

type processor_set struct {
	request chan func()
}

func (this *processor_set) new(n int) *processor_set {
	this.request = make(chan func(), n*2)
	for i := 0; i < n; i++ {
		go func() {
			for f := range this.request {
				f()
			}
		}()
	}
	return this
}

func (this *processor_set) exec(f func()) {
	this.request <- f
}

func (this *processor_set) stop() {
	close(this.request)
}
