package future

import "sync"

type Future[T any] struct {
	// 锁相关
	mu  sync.Mutex
	onc sync.Once

	// 结果相关
	result T
	err    error

	// 等待相关
	done chan struct{}
}
