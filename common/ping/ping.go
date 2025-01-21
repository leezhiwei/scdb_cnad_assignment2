package ping

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type PongResponse struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	pongResponse := PongResponse{
		Message:   "pong",
		Timestamp: time.Now(),
	}

	jsonResponse, err := json.Marshal(pongResponse)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Fatalln("Error marshalling JSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
