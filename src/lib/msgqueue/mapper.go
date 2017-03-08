package msgqueue

type EventMapper interface {
	MapEvent(string, []byte) (Event, error)
}

func NewEventMapper() EventMapper {
	return &StaticEventMapper{}
}