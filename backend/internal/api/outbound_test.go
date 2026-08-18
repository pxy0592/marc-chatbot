package api_test

import (
	"context"
	"github.com/marc-pango/marc-chatbot/backend/internal/domain"
	"github.com/marc-pango/marc-chatbot/backend/internal/testutil"
	"net/http"
	"testing"
)

func TestOutboundEndpoint(t *testing.T) {
	app := testutil.New(t)
	_ = app.Service.MockLogin(context.Background())
	testutil.Eventually(t, func() bool { b, _ := app.Service.GetBot(context.Background()); return b.Status == domain.BotOnline })
	groups, _ := app.Service.ListGroups(context.Background(), nil)
	g, _ := app.Service.SetGroupEnabled(context.Background(), groups[0].ID, true, "admin")
	headers := map[string]string{"Idempotency-Key": "api-key"}
	rr := app.Request(t, http.MethodPost, "/api/v1/groups/"+g.ID+"/messages", map[string]string{"content": "hello"}, headers)
	if rr.Code != http.StatusAccepted {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	rr2 := app.Request(t, http.MethodPost, "/api/v1/groups/"+g.ID+"/messages", map[string]string{"content": "hello"}, headers)
	if rr2.Code != http.StatusAccepted {
		t.Fatalf("status=%d", rr2.Code)
	}
}
