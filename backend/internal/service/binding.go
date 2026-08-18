package service

import (
	"context"
	"errors"
	"time"

	"github.com/marc-pango/marc-chatbot/backend/internal/domain"
	"github.com/marc-pango/marc-chatbot/backend/internal/store"
)

func (s *Service) GetBot(ctx context.Context) (domain.BotInstance, error) { return s.store.GetBot(ctx) }
func (s *Service) CurrentBinding(ctx context.Context) (domain.BindingSession, error) {
	return s.store.CurrentBinding(ctx)
}
func (s *Service) StartBinding(ctx context.Context, actor string) (domain.BindingSession, error) {
	if current, err := s.store.CurrentBinding(ctx); err == nil {
		if time.Now().Before(current.ExpiresAt) {
			return domain.BindingSession{}, ErrConflict
		}
		now := time.Now()
		current.Status = domain.BindingExpired
		current.QRCode = nil
		current.CompletedAt = &now
		_ = s.store.SaveBinding(ctx, current)
	} else if !errors.Is(err, store.ErrNotFound) {
		return domain.BindingSession{}, err
	}
	now := time.Now()
	binding := domain.BindingSession{ID: newID("bind"), Status: domain.BindingPending, CreatedAt: now, ExpiresAt: now.Add(2 * time.Minute)}
	if err := s.store.SaveBinding(ctx, binding); err != nil {
		return binding, err
	}
	b, err := s.store.GetBot(ctx)
	if err != nil {
		return binding, err
	}
	b.Status = domain.BotConnecting
	b.StatusReason = nil
	b.QRCode = nil
	b.QRExpiresAt = nil
	b.UpdatedAt = now
	if err := s.store.SaveBot(ctx, b); err != nil {
		return binding, err
	}
	if err := s.gateway.BeginBinding(ctx); err != nil {
		failedAt := time.Now()
		binding.Status = domain.BindingFailed
		binding.CompletedAt = &failedAt
		binding.FailureReason = safeDetail(err.Error())
		_ = s.store.SaveBinding(ctx, binding)
		s.audit(ctx, actor, "binding.start", "binding", binding.ID, domain.AuditFailure, err.Error())
		return binding, err
	}
	s.audit(ctx, actor, "binding.start", "binding", binding.ID, domain.AuditSuccess, "binding requested")
	s.events.Publish(AppEvent{Type: "binding.updated", Data: binding})
	return binding, nil
}
func (s *Service) CancelBinding(ctx context.Context, actor string) error {
	binding, err := s.store.CurrentBinding(ctx)
	if errors.Is(err, store.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if err := s.gateway.CancelBinding(ctx); err != nil {
		return err
	}
	now := time.Now()
	binding.Status = domain.BindingCancelled
	binding.QRCode = nil
	binding.CompletedAt = &now
	if err := s.store.SaveBinding(ctx, binding); err != nil {
		return err
	}
	b, err := s.store.GetBot(ctx)
	if err != nil {
		return err
	}
	b.Status = domain.BotUnbound
	b.QRCode = nil
	b.QRExpiresAt = nil
	b.UpdatedAt = now
	if err := s.store.SaveBot(ctx, b); err != nil {
		return err
	}
	s.audit(ctx, actor, "binding.cancel", "binding", binding.ID, domain.AuditSuccess, "binding cancelled")
	s.events.Publish(AppEvent{Type: "bot.status", Data: b})
	return nil
}
func (s *Service) MockLogin(ctx context.Context) error {
	controller, ok := s.gateway.(interface{ SimulateLogin(context.Context) error })
	if !ok {
		return ErrUnavailable
	}
	return controller.SimulateLogin(ctx)
}
func (s *Service) MockLogout(ctx context.Context, reason string) error {
	controller, ok := s.gateway.(interface {
		SimulateLogout(context.Context, string) error
	})
	if !ok {
		return ErrUnavailable
	}
	return controller.SimulateLogout(ctx, reason)
}
