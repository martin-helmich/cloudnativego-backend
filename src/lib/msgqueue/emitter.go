package msgqueue

type EventEmitter interface {
	Emit(e Event) error
}
