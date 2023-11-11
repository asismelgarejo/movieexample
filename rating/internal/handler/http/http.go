package http

import (
	"encoding/json"
	"log"
	"net/http"

	controller "movieexample.com/rating/internal/controller/rating"
	"movieexample.com/rating/pkg/model"
)

type Handler struct {
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	recordId := r.FormValue("id")
	if recordId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	recordType := r.FormValue("type")
	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		{
			val, err := h.ctrl.GetAggreatedRating(r.Context(), model.RecordType(recordType), model.RecordID(recordId))
			if err != nil {
				log.Printf("Repository error in get: %v", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if err = json.NewEncoder(w).Encode(val); err != nil {
				log.Printf("Response encode error: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	case http.MethodPost:
		{
			var rating model.Rating
			err := json.NewDecoder(r.Body).Decode(&rating)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = h.ctrl.Put(r.Context(), model.RecordType(recordType), model.RecordID(recordId), &rating)
			if err != nil {
				log.Printf("Repository error in put: %v", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)

		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}
}
