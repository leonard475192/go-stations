package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/leonard475192/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := model.HealthzResponse{
		Message: "OK",
	}
	msg_json, err := json.Marshal(msg)
	if err != nil {
		log.Printf("HealthzHandler.Readerror:%v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(msg_json))
}
