package rest

import (
	"net/http"

	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
	"bitbucket.org/minamartinteam/myevents/src/lib/persistence"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

func ServeAPI(endpoint string, dbHandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) error {
	handler := newEventHandler(dbHandler, eventEmitter)

	r := mux.NewRouter()
	eventsrouter := r.PathPrefix("/events").Subrouter()
	eventsrouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.findEventHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.allEventHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.newEventHandler)

	locationRouter := r.PathPrefix("/locations").Subrouter()
	locationRouter.Methods("GET").Path("").HandlerFunc(handler.allLocationsHandler)
	locationRouter.Methods("POST").Path("").HandlerFunc(handler.newLocationHandler)

	server := handlers.CORS()(r)

	return http.ListenAndServe(endpoint, server)
}
