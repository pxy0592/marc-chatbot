package service

import (
	"context"
	"github.com/marc-pango/marc-chatbot/backend/internal/domain"
)

func (s *Service) ListAudit(ctx context.Context, limit int) ([]domain.AuditEvent, error) {
	if limit < 1 || limit > 200 {
		limit = 100
	}
	return s.store.ListAudit(ctx, limit)
}
