package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/marc-pango/marc-chatbot/backend/internal/domain"
	"github.com/marc-pango/marc-chatbot/backend/internal/store"
)

func (s *Service) SendGroupMessage(ctx context.Context, groupID, content, key, actor string) (domain.OutboundMessageRequest, error) {
	content = strings.TrimSpace(content)
	if content == "" || len([]rune(content)) > 2000 || strings.TrimSpace(key) == "" || len(key) > 128 {
		return domain.OutboundMessageRequest{}, ErrValidation
	}
	if existing, err := s.store.GetOutboundByKey(ctx, key); err == nil {
		return existing, nil
	} else if !errors.Is(err, store.ErrNotFound) {
		return domain.OutboundMessageRequest{}, err
	}
	g, err := s.store.GetGroup(ctx, groupID)
	if err != nil {
		return domain.OutboundMessageRequest{}, err
	}
	if !g.Enabled || !g.Available {
		return domain.OutboundMessageRequest{}, ErrUnavailable
	}
	b, err := s.store.GetBot(ctx)
	if err != nil {
		return domain.OutboundMessageRequest{}, err
	}
	if b.Status != domain.BotOnline {
		return domain.OutboundMessageRequest{}, ErrUnavailable
	}
	now := time.Now()
	req := domain.OutboundMessageRequest{ID: newID("out"), IdempotencyKey: key, GroupID: groupID, Content: content, RequestedBy: actor, Status: domain.OutboundSending, CreatedAt: now, UpdatedAt: now}
	if err := s.store.SaveOutbound(ctx, req); err != nil {
		return req, err
	}
	providerID, sendErr := s.gateway.SendText(ctx, g.ProviderGroupID, content)
	req.UpdatedAt = time.Now()
	if sendErr != nil {
		req.Status = domain.OutboundFailed
		req.FailureReason = safeDetail(sendErr.Error())
		_ = s.store.SaveOutbound(ctx, req)
		s.audit(ctx, actor, "message.send", "group", groupID, domain.AuditFailure, sendErr.Error())
		s.events.Publish(AppEvent{Type: "message.sent", Data: req})
		return req, sendErr
	}
	req.Status = domain.OutboundSucceeded
	req.ProviderMessageID = ptr(providerID)
	if err := s.store.SaveOutbound(ctx, req); err != nil {
		return req, err
	}
	message := domain.GroupMessage{ID: newID("msg"), ProviderMessageID: providerID, GroupID: groupID, SenderID: "bot", SenderName: "Bot", Direction: domain.DirectionOutbound, MessageType: domain.MessageText, Content: content, OccurredAt: req.UpdatedAt, ReceivedAt: req.UpdatedAt, ProcessingStatus: domain.ProcessingSent, SelfMessage: true}
	_, _ = s.store.InsertMessage(ctx, message)
	s.audit(ctx, actor, "message.send", "group", groupID, domain.AuditSuccess, "text message sent")
	s.events.Publish(AppEvent{Type: "message.sent", Data: req})
	return req, nil
}
