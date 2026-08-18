package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/marc-pango/marc-chatbot/backend/internal/service"
	"github.com/marc-pango/marc-chatbot/backend/internal/store"
)

type Handler struct {
	service *service.Service
	events  *service.EventBus
}

func NewRouter(s *service.Service, events *service.EventBus, adminToken string, origins []string) http.Handler {
	h := &Handler{service: s, events: events}
	r := chi.NewRouter()
	r.Use(middleware.Recoverer, requestLog, cors(origins))
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(bearerAuth(adminToken))
		r.Get("/bot", h.getBot)
		r.Post("/bot/bindings", h.startBinding)
		r.Delete("/bot/bindings/current", h.cancelBinding)
		r.Get("/groups", h.listGroups)
		r.Patch("/groups/{groupId}", h.updateGroup)
		r.Post("/groups/{groupId}/messages", h.sendMessage)
		r.Get("/messages", h.listMessages)
		r.Get("/audit-events", h.listAudit)
		r.Get("/events", h.streamEvents)
		if s.GatewayDriver() == "mock" {
			r.Post("/mock/login", h.mockLogin)
			r.Post("/mock/logout", h.mockLogout)
			r.Post("/mock/messages", h.mockMessage)
		}
	})
	return r
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]string{"code": code, "message": message})
}
func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}
func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, "not_found", "resource not found")
	case errors.Is(err, service.ErrConflict):
		writeError(w, http.StatusConflict, "conflict", err.Error())
	case errors.Is(err, service.ErrValidation):
		writeError(w, http.StatusUnprocessableEntity, "validation_error", err.Error())
	case errors.Is(err, service.ErrUnavailable):
		writeError(w, http.StatusConflict, "unavailable", err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
	}
}
