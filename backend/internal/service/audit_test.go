package service_test

import (
	"context"
	"github.com/marc-pango/marc-chatbot/backend/internal/testutil"
	"strings"
	"testing"
)

func TestAuditDoesNotContainTokens(t *testing.T) {
	app := testutil.New(t)
	_, _ = app.Service.StartBinding(context.Background(), "admin")
	events, err := app.Service.ListAudit(context.Background(), 100)
	if err != nil {
		t.Fatal(err)
	}
	for _, event := range events {
		if event.Detail != nil && strings.Contains(*event.Detail, "test-token") {
			t.Fatal("token leaked")
		}
	}
}
