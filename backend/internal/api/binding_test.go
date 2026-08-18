package api_test

import (
	"encoding/json"
	"github.com/marc-pango/marc-chatbot/backend/internal/domain"
	"github.com/marc-pango/marc-chatbot/backend/internal/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBotRequiresAuthentication(t *testing.T) {
	app := testutil.New(t)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/bot", nil)
	rr := httptest.NewRecorder()
	app.Handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d", rr.Code)
	}
}
func TestStartBindingEndpoint(t *testing.T) {
	app := testutil.New(t)
	rr := app.Request(t, http.MethodPost, "/api/v1/bot/bindings", nil, nil)
	if rr.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var b domain.BindingSession
	if err := json.Unmarshal(rr.Body.Bytes(), &b); err != nil {
		t.Fatal(err)
	}
}
