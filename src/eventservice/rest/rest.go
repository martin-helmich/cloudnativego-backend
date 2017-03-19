package rest

import (
	"net/http"

	"bitbucket.org/minamartinteam/myevents/src/lib/persistence"
	"github.com/gorilla/mux"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
)

func ServeAPI(endpoint string, dbHandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) error {

	handler := newEventHandler(dbHandler, eventEmitter)
	r := mux.NewRouter()
	eventsrouter := r.PathPrefix("/events").Subrouter()
	eventsrouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.findEventHandler)
	eventsrouter.Methods("GET").Path("/all").HandlerFunc(handler.allEventHandler)
	eventsrouter.Methods("POST").Path("/New/").HandlerFunc(handler.newEventHandler)
	return http.ListenAndServe(endpoint, r)
}
