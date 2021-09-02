package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/leonard475192/go-stations/model"
	"github.com/leonard475192/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	switch r.Method {
	case "POST":
		w.Header().Set("Content-Type", "application/json")
		var req model.CreateTODORequest
		json.NewDecoder(r.Body).Decode(&req)
		res, err := h.Create(ctx, &req)
		if err != nil {
			w.WriteHeader(400)
			log.Print(err)
		}
		res_json, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(400)
			log.Print(err)
		}

		w.Write(res_json)
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	query, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return &model.CreateTODOResponse{}, err
	}
	res := model.CreateTODOResponse{
		TODO: *query,
	}
	return &res, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}
