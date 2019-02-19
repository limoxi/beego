package event

type IBroker interface {
	Send(asyncEvent *AsyncEvent)
}