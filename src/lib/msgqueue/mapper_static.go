package msgqueue

import (
	"fmt"
	"encoding/json"
	"bitbucket.org/minamartinteam/myevents/src/contracts"
)

type StaticEventMapper struct {}

func (e *StaticEventMapper) MapEvent(eventName string, serialized []byte) (Event, error) {
	var event Event

	switch eventName {
	case "eventCreated":
		event = &contracts.EventCreatedEvent{}
	case "locationCreated":
		event = &contracts.LocationCreatedEvent{}
	case "eventBooked":
		event = &contracts.EventBookedEvent{}
	default:
		return nil, fmt.Errorf("unknown event type %s", eventName)
	}

	err := json.Unmarshal(serialized, event)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal event %s: %s", eventName, err)
	}

	return event, nil
}