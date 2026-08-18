package service_test

import (
	"context"
	"github.com/marc-pango/marc-chatbot/backend/internal/bot"
	"github.com/marc-pango/marc-chatbot/backend/internal/testutil"
	"testing"
	"time"
)

func TestInboundMessageDeduplicates(t *testing.T) {
	app := testutil.New(t)
	_ = app.Service.MockLogin(context.Background())
	testutil.Eventually(t, func() bool { g, _ := app.Service.ListGroups(context.Background(), nil); return len(g) > 0 })
	groups, _ := app.Service.ListGroups(context.Background(), nil)
	g, _ := app.Service.SetGroupEnabled(context.Background(), groups[0].ID, true, "admin")
	msg := bot.Message{ProviderMessageID: "same-id", ProviderGroupID: g.ProviderGroupID, GroupName: g.Name, SenderID: "u1", SenderName: "User", Text: "hello", MessageType: "text", OccurredAt: time.Now()}
	_ = app.Service.InjectMockMessage(context.Background(), msg)
	_ = app.Service.InjectMockMessage(context.Background(), msg)
	testutil.Eventually(t, func() bool { m, _ := app.Service.ListMessages(context.Background(), g.ID, "", 100); return len(m) == 1 })
}
