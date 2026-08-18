package api_test

import (
	"context"
	"encoding/json"
	"github.com/marc-pango/marc-chatbot/backend/internal/domain"
	"github.com/marc-pango/marc-chatbot/backend/internal/testutil"
	"net/http"
	"testing"
)

func TestGroupsEndpoint(t *testing.T) {
	app := testutil.New(t)
	_ = app.Service.MockLogin(context.Background())
	testutil.Eventually(t, func() bool { g, _ := app.Service.ListGroups(context.Background(), nil); return len(g) > 0 })
	rr := app.Request(t, http.MethodGet, "/api/v1/groups", nil, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d", rr.Code)
	}
	var groups []domain.Group
	_ = json.Unmarshal(rr.Body.Bytes(), &groups)
	if len(groups) == 0 {
		t.Fatal("missing groups")
	}
	rr = app.Request(t, http.MethodPatch, "/api/v1/groups/"+groups[0].ID, map[string]bool{"enabled": true}, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
}
