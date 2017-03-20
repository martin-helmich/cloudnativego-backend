package rest

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"bitbucket.org/minamartinteam/myevents/src/lib/persistence"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
)

func ServeAPI(listenAddr string, database persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) {
	r := mux.NewRouter()
	r.Methods("get").Path("/events/{eventID}/bookings").Handler(&ListBookingHandler{})
	r.Methods("post").Path("/events/{eventID}/bookings").Handler(&CreateBookingHandler{})

	srv := http.Server{
		Handler:      r,
		Addr:         listenAddr,
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	srv.ListenAndServe()
}
