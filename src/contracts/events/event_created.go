package events

import "time"

// EventCreatedEvent is emitted whenever a new event is created
type EventCreatedEvent struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Start time.Time `json:"start_date"`
	End   time.Time `json:"end_date"`
}

// EventName returns the event's name
func (c *EventCreatedEvent) EventName() string {
	return "eventCreated";
}