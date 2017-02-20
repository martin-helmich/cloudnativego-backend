package rest

import "net/http"

type CreateBookingHandler struct {}

func (h *CreateBookingHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello World!"))
}
