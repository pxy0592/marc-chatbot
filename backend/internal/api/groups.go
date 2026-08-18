package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) listGroups(w http.ResponseWriter, r *http.Request) {
	var enabled *bool
	if raw := r.URL.Query().Get("enabled"); raw != "" {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			writeError(w, http.StatusUnprocessableEntity, "validation_error", "enabled must be boolean")
			return
		}
		enabled = &v
	}
	groups, err := h.service.ListGroups(r.Context(), enabled)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, groups)
}
func (h *Handler) updateGroup(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Enabled bool `json:"enabled"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	g, err := h.service.SetGroupEnabled(r.Context(), chi.URLParam(r, "groupId"), body.Enabled, "admin")
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, g)
}
