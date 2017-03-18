package rest

import (
	"net/http"

	"bitbucket.org/minamartinteam/myevents/src/lib/persistence"

	"fmt"

	"encoding/hex"
	"strings"

	"encoding/json"

	"github.com/gorilla/mux"
)

type eventSericeHandler struct {
	dbhandler persistence.DatabaseHandler
}

func newEventHandler(databasehandler persistence.DatabaseHandler) *eventSericeHandler {
	return &eventSericeHandler{
		dbhandler: databasehandler,
	}
}

func (eh *eventSericeHandler) findEventHandler(w http.ResponseWriter, r *http.Request) {
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

func (eh *eventSericeHandler) allEventHandler(w http.ResponseWriter, r *http.Request) {
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

func (eh *eventSericeHandler) newEventHandler(w http.ResponseWriter, r *http.Request) {
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
	fmt.Fprint(w, hex.EncodeToString(id))
}
