package rest

import "net/http"

type ListBookingHandler struct {}

func (h *ListBookingHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello World!"))
}
