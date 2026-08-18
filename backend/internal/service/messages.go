package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/marc-pango/marc-chatbot/backend/internal/bot"
	"github.com/marc-pango/marc-chatbot/backend/internal/domain"
	"github.com/marc-pango/marc-chatbot/backend/internal/store"
)

func (s *Service) receiveGatewayMessage(ctx context.Context, msg bot.Message) error {
	g, err := s.store.GetGroupByProviderID(ctx, msg.ProviderGroupID)
	if errors.Is(err, store.ErrNotFound) {
		g, err = s.upsertBotGroup(ctx, bot.Group{ProviderID: msg.ProviderGroupID, Name: msg.GroupName})
	}
	if err != nil {
		return err
	}
	kind := domain.MessageUnsupported
	if msg.MessageType == "text" {
		kind = domain.MessageText
	}
	status := domain.ProcessingIgnored
	if g.Enabled && kind == domain.MessageText {
		status = domain.ProcessingReceived
	}
	content := strings.TrimSpace(msg.Text)
	if len([]rune(content)) > 4000 {
		content = string([]rune(content)[:4000])
	}
	now := time.Now()
	occurred := msg.OccurredAt
	if occurred.IsZero() {
		occurred = now
	}
	m := domain.GroupMessage{ID: newID("msg"), ProviderMessageID: msg.ProviderMessageID, GroupID: g.ID, SenderID: msg.SenderID, SenderName: msg.SenderName, Direction: domain.DirectionInbound, MessageType: kind, Content: content, OccurredAt: occurred, ReceivedAt: now, ProcessingStatus: status, SelfMessage: msg.Self}
	inserted, err := s.store.InsertMessage(ctx, m)
	if err != nil {
		return err
	}
	if inserted && status == domain.ProcessingReceived {
		s.events.Publish(AppEvent{Type: "message.received", Data: m})
	}
	return nil
}
func (s *Service) ListMessages(ctx context.Context, groupID string, direction domain.MessageDirection, limit int) ([]domain.GroupMessage, error) {
	if limit < 1 || limit > 200 {
		limit = 100
	}
	return s.store.ListMessages(ctx, groupID, direction, limit)
}
func (s *Service) InjectMockMessage(ctx context.Context, msg bot.Message) error {
	controller, ok := s.gateway.(interface {
		InjectMessage(context.Context, bot.Message) error
	})
	if !ok {
		return ErrUnavailable
	}
	return controller.InjectMessage(ctx, msg)
}
