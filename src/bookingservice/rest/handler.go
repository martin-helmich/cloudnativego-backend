package rest

import (
	"bitbucket.org/minamartinteam/myevents/src/contracts"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
	"bitbucket.org/minamartinteam/myevents/src/lib/persistence"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type bookingHandler struct {
	eventEmitter msgqueue.EventEmitter
	database     persistence.DatabaseHandler
}

type bookingRequest struct {
	Seats int `json:"seats"`
}

func newBookingHandler(eventEmitter msgqueue.EventEmitter, database persistence.DatabaseHandler) (*bookingHandler, error) {
	return &bookingHandler{
		eventEmitter: eventEmitter,
		database:     database,
	}, nil
}
func (h *bookingHandler) handleNewBooking(w http.ResponseWriter, r *http.Request) {
	routeVars := mux.Vars(r)
	eventID, ok := routeVars["eventID"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, "missing route parameter 'eventID'")
		return
	}

	event, err := h.database.FindEvent(eventID)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "event %s could not be loaded: %s", eventID, err)
		return
	}

	bookingRequest := bookingRequest{}
	err = json.NewDecoder(r).Decode(&bookingRequest)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "could not decode JSON body: %s", err)
		return
	}

	if bookingRequest.Seats <= 0 {
		w.WriteHeader(400)
		fmt.Fprintf(w, "seat number must be positive (was %d)", bookingRequest.Seats)
		return
	}

	booking := persistence.Booking{
		Date:    time.Now().Unix(),
		EventID: event.ID,
		Seats:   bookingRequest.Seats,
	}

	msg := contracts.EventBookedEvent{
		EventID: event.ID.Hex(),
		UserID:  "someUserID",
	}
	h.eventEmitter.Emit(&msg)

	h.database.AddBookingForUser([]byte("someUserID"), booking)
}
