package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/marc-pango/marc-chatbot/backend/internal/bot"
	"github.com/marc-pango/marc-chatbot/backend/internal/domain"
	"github.com/marc-pango/marc-chatbot/backend/internal/store"
)

var (
	ErrConflict    = errors.New("resource state conflict")
	ErrValidation  = errors.New("validation failed")
	ErrUnavailable = errors.New("resource unavailable")
)

type Service struct {
	store   *store.Store
	gateway bot.Gateway
	events  *EventBus
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func New(ctx context.Context, st *store.Store, gateway bot.Gateway, events *EventBus) (*Service, error) {
	s := &Service{store: st, gateway: gateway, events: events}
	if _, err := st.GetBot(ctx); errors.Is(err, store.ErrNotFound) {
		now := time.Now()
		if err := st.SaveBot(ctx, domain.BotInstance{ID: "primary", Status: domain.BotUnbound, UpdatedAt: now}); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	consumerCtx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.wg.Add(1)
	go s.consume(consumerCtx)
	return s, nil
}
func (s *Service) Close() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
}
func (s *Service) GatewayDriver() string { return s.gateway.Driver() }
func newID(prefix string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return prefix + "_" + hex.EncodeToString(b)
}
func ptr[T any](v T) *T { return &v }
func safeDetail(v string) *string {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	if len(v) > 300 {
		v = v[:300]
	}
	return &v
}
func (s *Service) audit(ctx context.Context, actor, action, targetType, targetID string, result domain.AuditResult, detail string) {
	_ = s.store.InsertAudit(ctx, domain.AuditEvent{ID: newID("audit"), Actor: actor, Action: action, TargetType: targetType, TargetID: targetID, Result: result, Detail: safeDetail(detail), CreatedAt: time.Now()})
}
func (s *Service) consume(ctx context.Context) {
	defer s.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-s.gateway.Events():
			if err := s.handleGatewayEvent(context.Background(), event); err != nil {
				slog.Error("handle bot event", "type", event.Type, "error", err)
			}
		}
	}
}
func (s *Service) handleGatewayEvent(ctx context.Context, event bot.Event) error {
	switch event.Type {
	case bot.EventScan:
		binding, err := s.store.CurrentBinding(ctx)
		if err != nil {
			return err
		}
		binding.QRCode = ptr(event.QRCode)
		binding.ExpiresAt = event.QRExpiresAt
		if err := s.store.SaveBinding(ctx, binding); err != nil {
			return err
		}
		b, err := s.store.GetBot(ctx)
		if err != nil {
			return err
		}
		b.Status = domain.BotAwaitingScan
		b.QRCode = ptr(event.QRCode)
		b.QRExpiresAt = ptr(event.QRExpiresAt)
		b.UpdatedAt = time.Now()
		if err := s.store.SaveBot(ctx, b); err != nil {
			return err
		}
		s.events.Publish(AppEvent{Type: "binding.updated", Data: binding})
		s.events.Publish(AppEvent{Type: "bot.status", Data: b})
	case bot.EventLogin:
		now := time.Now()
		b, err := s.store.GetBot(ctx)
		if err != nil {
			return err
		}
		b.Status = domain.BotOnline
		b.PublicAccountID = ptr(event.AccountID)
		b.DisplayName = ptr(event.DisplayName)
		b.StatusReason = nil
		b.QRCode = nil
		b.QRExpiresAt = nil
		b.LastSeenAt = &now
		b.UpdatedAt = now
		if err := s.store.SaveBot(ctx, b); err != nil {
			return err
		}
		if binding, err := s.store.CurrentBinding(ctx); err == nil {
			binding.Status = domain.BindingConfirmed
			binding.QRCode = nil
			binding.CompletedAt = &now
			_ = s.store.SaveBinding(ctx, binding)
		}
		s.audit(ctx, "system", "bot.login", "bot", b.ID, domain.AuditSuccess, "account connected")
		s.events.Publish(AppEvent{Type: "bot.status", Data: b})
	case bot.EventLogout:
		return s.setBotOffline(ctx, event.Reason, domain.BotOffline, "bot.logout")
	case bot.EventError:
		return s.setBotOffline(ctx, event.Reason, domain.BotError, "bot.error")
	case bot.EventGroups:
		for _, g := range event.Groups {
			if _, err := s.upsertBotGroup(ctx, g); err != nil {
				return err
			}
		}
	case bot.EventMessage:
		if event.Message != nil {
			return s.receiveGatewayMessage(ctx, *event.Message)
		}
	}
	return nil
}
func (s *Service) setBotOffline(ctx context.Context, reason string, status domain.BotStatus, action string) error {
	b, err := s.store.GetBot(ctx)
	if err != nil {
		return err
	}
	b.Status = status
	b.StatusReason = safeDetail(reason)
	b.UpdatedAt = time.Now()
	if err := s.store.SaveBot(ctx, b); err != nil {
		return err
	}
	s.audit(ctx, "system", action, "bot", b.ID, domain.AuditSuccess, reason)
	s.events.Publish(AppEvent{Type: "bot.status", Data: b})
	return nil
}
