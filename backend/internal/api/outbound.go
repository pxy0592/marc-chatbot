package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) sendMessage(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Idempotency-Key")
	var body struct {
		Content string `json:"content"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	req, err := h.service.SendGroupMessage(r.Context(), chi.URLParam(r, "groupId"), body.Content, key, "admin")
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusAccepted, req)
}
