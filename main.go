package futurego

import (
	"errors"
	"time"
)

func newFuture[T any]() *Future[T] {
	return &Future[T]{
		done: make(chan struct{}),
	}
}

// VoidAsync runs a function asynchronously that returns no value but may return an error.
// Returns a Future[struct{}], which can be used to wait or check for errors.
func VoidAsync(fn func() error) *Future[struct{}] {
	return newFuture[struct{}]().asyncVoid(func() error {
		return fn()
	})
}

// Async runs a function asynchronously that returns a value and an error.
// Returns a Future[T], and the result can be retrieved via Get.
func Async[T any](fn func() (T, error)) *Future[T] {
	return newFuture[T]().async(fn)
}

// WaitAll blocks until all provided Futures are completed.
func WaitAll(futures ...future) {
	for _, future := range futures {
		<-future.sDone()
	}
}

// WaitAllWithTimeout blocks until all Futures are completed or the timeout is reached.
// Returns error "wait all timeout" if the timeout occurs.
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

// Get blocks until the asynchronous task is complete and returns the result and error.
func (c *Future[T]) Get() (T, error) {
	<-c.done
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.result, c.err
}

// GetWithTimout waits for the asynchronous result up to the specified timeout duration.
// Returns zero value and a timeout error if the wait exceeds the limit.
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

// Error returns the error after the asynchronous task is completed.
func (c *Future[T]) Error() error {
	_, err := c.Get()
	return err
}

// IsDone checks whether the current Future is completed.
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
