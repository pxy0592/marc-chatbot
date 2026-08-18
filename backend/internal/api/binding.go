package api

import (
	"net/http"
)

func (h *Handler) getBot(w http.ResponseWriter, r *http.Request) {
	b, err := h.service.GetBot(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, b)
}
func (h *Handler) startBinding(w http.ResponseWriter, r *http.Request) {
	b, err := h.service.StartBinding(r.Context(), "admin")
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, b)
}
func (h *Handler) cancelBinding(w http.ResponseWriter, r *http.Request) {
	if err := h.service.CancelBinding(r.Context(), "admin"); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) mockLogin(w http.ResponseWriter, r *http.Request) {
	if err := h.service.MockLogin(r.Context()); err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "accepted"})
}
func (h *Handler) mockLogout(w http.ResponseWriter, r *http.Request) {
	if err := h.service.MockLogout(r.Context(), "mock logout"); err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "accepted"})
}
