package service_test

import (
	"context"
	"github.com/marc-pango/marc-chatbot/backend/internal/domain"
	"github.com/marc-pango/marc-chatbot/backend/internal/testutil"
	"testing"
)

func TestOutboundIsIdempotent(t *testing.T) {
	app := testutil.New(t)
	_ = app.Service.MockLogin(context.Background())
	testutil.Eventually(t, func() bool { b, _ := app.Service.GetBot(context.Background()); return b.Status == domain.BotOnline })
	groups, _ := app.Service.ListGroups(context.Background(), nil)
	g, _ := app.Service.SetGroupEnabled(context.Background(), groups[0].ID, true, "admin")
	first, err := app.Service.SendGroupMessage(context.Background(), g.ID, "hello", "key-1", "admin")
	if err != nil {
		t.Fatal(err)
	}
	second, err := app.Service.SendGroupMessage(context.Background(), g.ID, "hello", "key-1", "admin")
	if err != nil {
		t.Fatal(err)
	}
	if first.ID != second.ID {
		t.Fatal("idempotency failed")
	}
	n, _ := app.Store.DebugCount(context.Background(), "group_messages")
	if n != 1 {
		t.Fatalf("messages=%d", n)
	}
}
