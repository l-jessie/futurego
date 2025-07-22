package future

type future interface {
	sDone() <-chan struct{}
}
