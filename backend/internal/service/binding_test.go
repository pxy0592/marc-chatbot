package service_test

import (
	"context"
	"testing"

	"github.com/marc-pango/marc-chatbot/backend/internal/domain"
	"github.com/marc-pango/marc-chatbot/backend/internal/testutil"
)

func TestBindingAndLoginLifecycle(t *testing.T) {
	app := testutil.New(t)
	binding, err := app.Service.StartBinding(context.Background(), "admin")
	if err != nil {
		t.Fatal(err)
	}
	if binding.Status != domain.BindingPending {
		t.Fatalf("status=%s", binding.Status)
	}
	testutil.Eventually(t, func() bool {
		b, _ := app.Service.GetBot(context.Background())
		return b.Status == domain.BotAwaitingScan && b.QRCode != nil
	})
	if err := app.Service.MockLogin(context.Background()); err != nil {
		t.Fatal(err)
	}
	testutil.Eventually(t, func() bool {
		b, _ := app.Service.GetBot(context.Background())
		return b.Status == domain.BotOnline && b.DisplayName != nil
	})
}
