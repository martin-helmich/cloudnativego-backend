package events

import "time"

type EventCreatedEvent struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Start time.Time `json:"start_date"`
	End   time.Time `json:"end_date"`
}

func (c *EventCreatedEvent) EventName() string {
	return "eventCreated";
}