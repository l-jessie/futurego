package future

import (
	"errors"
	"time"
)

func newFuture[T any]() *Future[T] {
	return &Future[T]{
		done: make(chan struct{}),
	}
}

// VoidAsync 异步执行一个无返回值但可能出错的函数。
// 返回一个 Future[struct{}]，可用于等待或检查错误。
func VoidAsync(fn func() error) *Future[struct{}] {
	return newFuture[struct{}]().asyncVoid(func() error {
		return fn()
	})
}

// Async 异步执行一个带返回值和错误的函数。
// 返回一个 Future[T]，通过 Get 获取结果。
func Async[T any](fn func() (T, error)) *Future[T] {
	return newFuture[T]().async(fn)
}

// WaitAll 会阻塞直到所有传入的 Future 完成。
func WaitAll(futures ...future) {
	for _, future := range futures {
		<-future.sDone()
	}
}

// WaitAllWithTimeout 会阻塞直到所有 Future 完成或超时。
// 超时返回错误 "wait all timeout"。
func WaitAllWithTimeout(timeout time.Duration, futures ...future) error {
	done := make(chan struct{})
	go func() {
		for _, f := range futures {
			<-f.sDone()
		}
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return errors.New("wait all timeout")
	}
}

func (c *Future[T]) asyncVoid(fn func() error) *Future[T] {
	c.onc.Do(func() {
		go func() {
			err := fn()

			c.mu.Lock()
			c.err = err
			c.mu.Unlock()
			close(c.done)
		}()
	})
	return c
}

func (c *Future[T]) async(fn func() (T, error)) *Future[T] {
	c.onc.Do(func() {
		go func() {
			result, err := fn()

			c.mu.Lock()
			c.result = result
			c.err = err
			c.mu.Unlock()
			close(c.done)
		}()
	})
	return c
}

// Get 阻塞直到异步任务完成并返回结果与错误。
func (c *Future[T]) Get() (T, error) {
	<-c.done
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.result, c.err
}

// GetWithTimout 在指定超时时间内等待异步结果。
// 如果超时，返回零值和超时错误。
func (c *Future[T]) GetWithTimout(afterTime time.Duration) (T, error) {
	select {
	case <-c.done:
		c.mu.Lock()
		defer c.mu.Unlock()
		return c.result, c.err
	case <-time.After(afterTime):
		var zero T
		return zero, errors.New("future get timeout~")
	}
}

// Error 返回异步任务完成后的错误信息。
func (c *Future[T]) Error() error {
	_, err := c.Get()
	return err
}

// IsDone 判断当前 Future 是否已完成。
func (c *Future[T]) IsDone() bool {
	select {
	case <-c.sDone():
		return true
	default:
		return false
	}
}
func (f *Future[T]) sDone() <-chan struct{} {
	return f.done
}
