package service

import (
	"context"
	"time"

	"github.com/marc-pango/marc-chatbot/backend/internal/bot"
	"github.com/marc-pango/marc-chatbot/backend/internal/domain"
)

func (s *Service) upsertBotGroup(ctx context.Context, g bot.Group) (domain.Group, error) {
	now := time.Now()
	existing, err := s.store.GetGroupByProviderID(ctx, g.ProviderID)
	if err == nil {
		existing.Name = g.Name
		existing.MemberCount = g.MemberCount
		existing.Available = true
		existing.UpdatedAt = now
		return s.store.UpsertGroup(ctx, existing)
	}
	return s.store.UpsertGroup(ctx, domain.Group{ID: newID("group"), ProviderGroupID: g.ProviderID, Name: g.Name, MemberCount: g.MemberCount, Enabled: false, Available: true, DiscoveredAt: now, UpdatedAt: now})
}
func (s *Service) SyncGroups(ctx context.Context) ([]domain.Group, error) {
	groups, err := s.gateway.ListGroups(ctx)
	if err != nil {
		return nil, err
	}
	for _, g := range groups {
		if _, err := s.upsertBotGroup(ctx, g); err != nil {
			return nil, err
		}
	}
	return s.store.ListGroups(ctx, nil)
}
func (s *Service) ListGroups(ctx context.Context, enabled *bool) ([]domain.Group, error) {
	return s.store.ListGroups(ctx, enabled)
}
func (s *Service) SetGroupEnabled(ctx context.Context, id string, enabled bool, actor string) (domain.Group, error) {
	g, err := s.store.SetGroupEnabled(ctx, id, enabled)
	if err != nil {
		return g, err
	}
	action := "group.disable"
	if enabled {
		action = "group.enable"
	}
	s.audit(ctx, actor, action, "group", id, domain.AuditSuccess, g.Name)
	s.events.Publish(AppEvent{Type: "group.updated", Data: g})
	return g, nil
}
