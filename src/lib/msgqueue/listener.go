package msgqueue

import "reflect"

type EventListener interface {
	Listen(events ...string) (<-chan Event, <-chan error, error)
	Map(eventName string, typ reflect.Type)
}