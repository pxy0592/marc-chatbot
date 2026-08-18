package api_test

import (
	"context"
	"github.com/marc-pango/marc-chatbot/backend/internal/testutil"
	"net/http"
	"strings"
	"testing"
)

func TestAuditEndpoint(t *testing.T) {
	app := testutil.New(t)
	_, _ = app.Service.StartBinding(context.Background(), "admin")
	rr := app.Request(t, http.MethodGet, "/api/v1/audit-events", nil, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "binding.start") {
		t.Fatalf("body=%s", rr.Body.String())
	}
}
