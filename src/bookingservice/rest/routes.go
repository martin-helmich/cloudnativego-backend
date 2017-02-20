package rest

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"sync"
)

func ServeAPI(wg *sync.WaitGroup) {
	defer wg.Done()

	r := mux.NewRouter()
	r.Methods("get").Path("/events/{eventID}/bookings").Handler(&ListBookingHandler{})
	r.Methods("post").Path("/events/{eventID}/bookings").Handler(&CreateBookingHandler{})

	srv := http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:3000",
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	srv.ListenAndServe()
}
