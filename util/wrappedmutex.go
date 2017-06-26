package util

import (
	"context"
	"github.com/sasha-s/go-deadlock"
	"sync"
)

type WrappedMutex struct {
	Name     string
	mutex    sync.Mutex
	devMutex deadlock.Mutex
}

func (wrapped *WrappedMutex) Lock(ctx context.Context) context.Context {
	if Testing {
		wrapped.devMutex.Lock()
	} else {
		wrapped.mutex.Lock()
	}
	return ctx
}

func (wrapped *WrappedMutex) Unlock(ctx context.Context) {
	if Testing {
		wrapped.devMutex.Unlock()
	} else {
		wrapped.mutex.Unlock()
	}
}

///////////////////////////////////////////////////////////////////////////////

type WrappedRWMutex struct {
	Name     string
	mutex    sync.RWMutex
	devMutex deadlock.RWMutex
}

func (wrapped *WrappedRWMutex) Lock(ctx context.Context) context.Context {
	if Testing {
		wrapped.devMutex.Lock()
	} else {
		wrapped.mutex.Lock()
	}
	return ctx
}

func (wrapped *WrappedRWMutex) Unlock(ctx context.Context) {
	if Testing {
		wrapped.devMutex.Unlock()
	} else {
		wrapped.mutex.Unlock()
	}
}

func (wrapped *WrappedRWMutex) RLock(ctx context.Context) context.Context {
	if Testing {
		wrapped.devMutex.RLock()
	} else {
		wrapped.mutex.RLock()
	}
	return ctx
}

func (wrapped *WrappedRWMutex) RUnlock(ctx context.Context) {
	if Testing {
		wrapped.devMutex.RUnlock()
	} else {
		wrapped.mutex.RUnlock()
	}
}
