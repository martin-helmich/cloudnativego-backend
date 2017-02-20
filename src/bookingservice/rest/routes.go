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
	r.Methods("get").Handler("/events/{eventID}/bookings", &ListBookingHandler{})
	r.Methods("post").Handler("/events/{eventID}/bookings", &CreateBookingHandler{})

	srv := http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:3000",
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	srv.ListenAndServe()
}
