package rest

import (
	"net/http"

	"bitbucket.org/minamartinteam/myevents/src/lib/persistence"

	"fmt"

	"encoding/hex"
	"strings"

	"encoding/json"

	"github.com/gorilla/mux"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
	"bitbucket.org/minamartinteam/myevents/src/contracts"
	"time"
)

type eventServiceHandler struct {
	dbhandler persistence.DatabaseHandler
	eventEmitter msgqueue.EventEmitter
}

func newEventHandler(databasehandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) *eventServiceHandler {
	return &eventServiceHandler{
		dbhandler: databasehandler,
		eventEmitter: eventEmitter,
	}
}

func (eh *eventServiceHandler) findEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	criteria, ok := vars["SearchCriteria"]
	if !ok {
		fmt.Fprint(w, `No search criteria found, you can either search by id via /id/4
						to search by name via /name/coldplayconcert`)
		return
	}

	searchkey, ok := vars["search"]
	if !ok {
		fmt.Fprint(w, `No search keys found, you can either search by id via /id/4
						to search by name via /name/coldplayconcert`)
		return
	}

	var event persistence.Event
	var err error
	switch strings.ToLower(criteria) {
	case "name":
		event, err = eh.dbhandler.FindEventByName(searchkey)
	case "id":
		id, err := hex.DecodeString(searchkey)
		if nil == err {
			event, err = eh.dbhandler.FindEvent(id)
		}
	}
	if err != nil {
		fmt.Fprintf(w, "Error occured %s", err)
	}
	json.NewEncoder(w).Encode(&event)
}

func (eh *eventServiceHandler) allEventHandler(w http.ResponseWriter, r *http.Request) {
	events, err := eh.dbhandler.FindAllAvailableEvents()
	if err != nil {
		fmt.Fprintf(w, "Error occrued while trying to find all available events %s", err)
		return
	}
	err = json.NewEncoder(w).Encode(&events)
	if err != nil {
		fmt.Fprintf(w, "Error occrued while trying encode events to JSON %s", err)
	}
}

func (eh *eventServiceHandler) newEventHandler(w http.ResponseWriter, r *http.Request) {
	event := persistence.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if nil != err {
		fmt.Fprintf(w, "error occured while decoding event data %s", err)
		return
	}
	id, err := eh.dbhandler.AddEvent(event)
	if nil != err {
		fmt.Fprintf(w, "error occured while decoding event data %s", err)
		return
	}

	msg := contracts.EventCreatedEvent{
		ID: hex.EncodeToString(id),
		Name: event.Name,
		Start: time.Unix(event.StartDate, 0),
		End: time.Unix(event.EndDate, 0),
		// LocationID: event.Location.ID,   TODO: Locations need an ID!
	}
	eh.eventEmitter.Emit(&msg)

	fmt.Fprint(w, hex.EncodeToString(id))
}
