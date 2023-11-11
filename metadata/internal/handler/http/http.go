package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"movieexample.com/metadata/internal/controller/metadata"
	"movieexample.com/metadata/internal/repository"
)

type Handler struct {
	controller *metadata.Controller
}

func New(c *metadata.Controller) *Handler {
	return &Handler{controller: c}
}

func (h *Handler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("Retrieving movie metadata with %v\n", id)

	ctx := r.Context()
	m, err := h.controller.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {
		log.Printf("Repository get error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	mMarshalled, err := json.Marshal(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(mMarshalled)
}
