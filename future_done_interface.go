package futurego

type future interface {
	sDone() <-chan struct{}
}
