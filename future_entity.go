package futurego

import "sync"

type Future[T any] struct {
	mu     sync.Mutex
	onc    sync.Once
	result T
	err    error
	done   chan struct{}
}
