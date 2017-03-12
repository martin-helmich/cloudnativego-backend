package rest

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bitbucket.org/minamartinteam/myevents/src/lib/msgqueue"
	"github.com/nu7hatch/gouuid"
	"bitbucket.org/minamartinteam/myevents/src/contracts"
)

type eventRef struct {
	ID string `json:"id"`
	Name string `json:"name,omitempty"`
}

type createBookingRequest struct {
	Event eventRef `json:"event"`
}

type createBookingResponse struct {
	ID string `json:"id"`
	Event eventRef `json:"event"`
}

type CreateBookingHandler struct{
	eventEmitter msgqueue.EventEmitter
}

func (h *CreateBookingHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil{
		respondWithError(res, "could not read request body", 400)
		return
	}

	request := createBookingRequest{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		respondWithError(res, "could not unserialize body", 400)
		return
	}

	bookingID, _ := uuid.NewV4()
	response := createBookingResponse{
		ID: bookingID.String(),
		Event: request.Event,
	}

	go func() {
		h.eventEmitter.Emit(&contracts.EventBookedEvent{
			EventID: request.Event.ID,
			UserID: "foo", // TODO: Authenticate user
		})
	}()

	responseBody, err := json.Marshal(&response)
	if err != nil {
		respondWithError(res, "could not serialize response", 500)
		return
	}

	res.WriteHeader(201)
	res.Header().Set("Content-Type", "application/json")
	res.Write(responseBody)
}
