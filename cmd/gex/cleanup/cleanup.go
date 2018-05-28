package cleanup

import (
	"sync"
)

type cleanupRequest struct {
	f func()
}

func (this *cleanupRequest) cleanup() {
	this.f()
}

var cleanupRequests = []*cleanupRequest{}
var lock sync.Mutex

func RegisterCleanup(f func()) *cleanupRequest {
	cf := &cleanupRequest{f}
	lock.Lock()
	defer lock.Unlock()

	cleanupRequests = append(cleanupRequests, cf)
	return cf
}

func Cleanup(f func()) func() {
	cf := RegisterCleanup(f)
	return func() {
		lock.Lock()

		for i, p := range cleanupRequests {
			if p == cf {
				cleanupRequests = append(cleanupRequests[0:i], cleanupRequests[i+1:]...)
				break
			}
		}
		lock.Unlock()
		cf.cleanup()
	}
}

func cleanup() {
	lock.Lock()
	defer lock.Unlock()

	for _, cf := range cleanupRequests {
		cf.cleanup()
	}
	cleanupRequests = []*cleanupRequest{}
}
