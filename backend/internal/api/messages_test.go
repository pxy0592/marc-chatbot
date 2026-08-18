package api_test

import (
	"context"
	"github.com/marc-pango/marc-chatbot/backend/internal/testutil"
	"net/http"
	"testing"
)

func TestMessagesEndpoint(t *testing.T) {
	app := testutil.New(t)
	_ = app.Service.MockLogin(context.Background())
	testutil.Eventually(t, func() bool { g, _ := app.Service.ListGroups(context.Background(), nil); return len(g) > 0 })
	groups, _ := app.Service.ListGroups(context.Background(), nil)
	_, _ = app.Service.SetGroupEnabled(context.Background(), groups[0].ID, true, "admin")
	rr := app.Request(t, http.MethodPost, "/api/v1/mock/messages", map[string]any{"providerMessageId": "api-in-1", "providerGroupId": groups[0].ProviderGroupID, "groupName": groups[0].Name, "senderId": "u", "senderName": "User", "text": "hello"}, nil)
	if rr.Code != http.StatusAccepted {
		t.Fatalf("status=%d", rr.Code)
	}
	testutil.Eventually(t, func() bool {
		m, _ := app.Service.ListMessages(context.Background(), groups[0].ID, "", 100)
		return len(m) == 1
	})
	rr = app.Request(t, http.MethodGet, "/api/v1/messages?groupId="+groups[0].ID, nil, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d", rr.Code)
	}
}
