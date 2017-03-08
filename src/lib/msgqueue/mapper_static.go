package msgqueue

import (
	"fmt"
	"encoding/json"
	"bitbucket.org/minamartinteam/myevents/src/contracts/events"
)

type StaticEventMapper struct {}

func (e *StaticEventMapper) MapEvent(eventName string, serialized []byte) (Event, error) {
	var event Event

	switch eventName {
	case "eventCreated":
		event = &events.EventCreatedEvent{}
	case "locationCreated":
		event = &events.LocationCreatedEvent{}
	case "eventBooked":
		event = &events.EventBookedEvent{}
	default:
		return nil, fmt.Errorf("unknown event type %s", eventName)
	}

	err := json.Unmarshal(serialized, event)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal event %s: %s", eventName, err)
	}

	return event, nil
}