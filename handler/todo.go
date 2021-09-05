package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/schema"
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
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		if err := r.ParseForm(); err != nil {
			log.Printf("error ParseForm:%v", err)
		}

		req := model.ReadTODORequest{}
		if err := schema.NewDecoder().Decode(&req, r.Form); err != nil {
			log.Printf("error Decoder:%v", err)
		}

		res, err := h.Read(ctx, &req)
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
	case "POST":
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
	case "PUT":
		var req model.UpdateTODORequest
		json.NewDecoder(r.Body).Decode(&req)
		res, err := h.Update(ctx, &req)
		// ここ聞く
		switch err {
		case nil:
			res_json, err := json.Marshal(res)
			if err != nil {
				w.WriteHeader(400)
				log.Print(err)
			}
			w.Write(res_json)
		case model.ErrNotFound{}:
			w.WriteHeader(404)
			w.Write(nil)
		default:
			log.Print(err)
			w.WriteHeader(400)
			w.Write(nil)
		}
	case "DELETE":
		var req model.DeleteTODORequest
		json.NewDecoder(r.Body).Decode(&req)
		// ここじゃなくて、model.error に ErrEmptyRequest とか作ったほうがきれいな気がしました。
		if len(req.IDs) == 0 {
			w.WriteHeader(400)
			w.Write(nil)
		}
		res, err := h.Delete(ctx, &req)
		switch err {
		case nil:
			res_json, err := json.Marshal(res)
			if err != nil {
				w.WriteHeader(400)
				log.Print(err)
			}
			w.Write(res_json)
		case model.ErrNotFound{}:
			w.WriteHeader(404)
			w.Write(nil)
		default:
			log.Print(err)
			w.WriteHeader(400)
			w.Write(nil)
		}
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
	prevID, size := int64(0), int64(5)
	if req.PrevID != 0 {
		prevID = req.PrevID
	}
	if req.Size != 0 {
		size = req.Size
	}
	query, err := h.svc.ReadTODO(ctx, prevID, size)
	if err != nil {
		return &model.ReadTODOResponse{}, err
	}
	res := model.ReadTODOResponse{
		TODOs: query,
	}
	return &res, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	query, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return &model.UpdateTODOResponse{}, err
	}
	res := model.UpdateTODOResponse{
		TODO: *query,
	}
	return &res, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDs)
	if err != nil {
		return &model.DeleteTODOResponse{}, err
	}
	res := model.DeleteTODOResponse{}
	return &res, nil
}
