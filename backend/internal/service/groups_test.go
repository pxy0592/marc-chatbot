package service_test

import (
	"context"
	"github.com/marc-pango/marc-chatbot/backend/internal/testutil"
	"testing"
)

func TestGroupSyncAndEnable(t *testing.T) {
	app := testutil.New(t)
	_ = app.Service.MockLogin(context.Background())
	testutil.Eventually(t, func() bool { g, _ := app.Service.ListGroups(context.Background(), nil); return len(g) == 2 })
	groups, _ := app.Service.ListGroups(context.Background(), nil)
	updated, err := app.Service.SetGroupEnabled(context.Background(), groups[0].ID, true, "admin")
	if err != nil {
		t.Fatal(err)
	}
	if !updated.Enabled {
		t.Fatal("group not enabled")
	}
}
