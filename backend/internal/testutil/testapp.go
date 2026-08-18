package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/marc-pango/marc-chatbot/backend/internal/api"
	"github.com/marc-pango/marc-chatbot/backend/internal/bot"
	"github.com/marc-pango/marc-chatbot/backend/internal/service"
	"github.com/marc-pango/marc-chatbot/backend/internal/store"
)

type App struct {
	Store   *store.Store
	Gateway *bot.MockGateway
	Service *service.Service
	Handler http.Handler
	Token   string
}

func New(t *testing.T) *App {
	t.Helper()
	st, err := store.Open(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatal(err)
	}
	gateway := bot.NewMockGateway()
	if err := gateway.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	events := service.NewEventBus()
	svc, err := service.New(context.Background(), st, gateway, events)
	if err != nil {
		t.Fatal(err)
	}
	app := &App{Store: st, Gateway: gateway, Service: svc, Handler: api.NewRouter(svc, events, "test-token", []string{"http://localhost:5173"}), Token: "test-token"}
	t.Cleanup(func() { svc.Close(); _ = gateway.Stop(context.Background()); _ = st.Close() })
	return app
}
func (a *App) Request(t *testing.T, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	var payload []byte
	if body != nil {
		payload, _ = json.Marshal(body)
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+a.Token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rr := httptest.NewRecorder()
	a.Handler.ServeHTTP(rr, req)
	return rr
}
func Eventually(t *testing.T, fn func() bool) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if fn() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("condition not met")
}
