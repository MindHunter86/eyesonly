package dynamic

import (
	"context"
	"sync"
)

type CallbackFunc func(context.Context)

var cm sync.RWMutex
var callbacks = make([]CallbackFunc, 0, 32)

func RegisterCallback(fn CallbackFunc) {
	cm.Lock()
	defer cm.Unlock()
	callbacks = append(callbacks, fn)
}

func forEachCallback(c context.Context) {
	cm.RLock()
	defer cm.RUnlock()

	if len(callbacks) == 0 {
		return
	}

	for _, cb := range callbacks {
		cb(c)
	}
}
