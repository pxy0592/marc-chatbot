package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/marc-pango/marc-chatbot/backend/internal/api"
	"github.com/marc-pango/marc-chatbot/backend/internal/bot"
	"github.com/marc-pango/marc-chatbot/backend/internal/config"
	"github.com/marc-pango/marc-chatbot/backend/internal/service"
	"github.com/marc-pango/marc-chatbot/backend/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("configuration error", "error", err)
		os.Exit(1)
	}
	st, err := store.Open(cfg.DatabasePath)
	if err != nil {
		slog.Error("database error", "error", err)
		os.Exit(1)
	}
	defer st.Close()
	var gateway bot.Gateway
	if cfg.BotDriver == "wechaty" {
		gateway = bot.NewWechatyGateway(cfg.WechatyPuppetToken)
	} else {
		gateway = bot.NewMockGateway()
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	events := service.NewEventBus()
	svc, err := service.New(ctx, st, gateway, events)
	if err != nil {
		slog.Error("service initialization failed", "error", err)
		os.Exit(1)
	}
	defer svc.Close()
	if err := gateway.Start(ctx); err != nil {
		slog.Error("bot gateway start failed", "driver", gateway.Driver(), "error", err)
		os.Exit(1)
	}
	defer gateway.Stop(context.Background())
	server := &http.Server{Addr: cfg.HTTPAddr, Handler: api.NewRouter(svc, events, cfg.AdminToken, cfg.CORSOrigins), ReadHeaderTimeout: 5 * time.Second}
	go func() {
		slog.Info("server started", "addr", cfg.HTTPAddr, "bot_driver", gateway.Driver())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			cancel()
		}
	}()
	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = server.Shutdown(shutdownCtx)
}
