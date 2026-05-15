package service

import (
	"context"
	"errors"

	"featureflags/audit-service/internal/model"
	"featureflags/audit-service/internal/repository"
)

const defaultListLimit = 50

var ErrValidation = errors.New("flag_name and action are required")

// AuditService validates and applies audit use cases.
type AuditService struct {
	repo *repository.AuditPostgres
}

func NewAuditService(repo *repository.AuditPostgres) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) Record(ctx context.Context, event model.AuditEvent) error {
	if event.FlagName == "" || event.Action == "" {
		return ErrValidation
	}
	return s.repo.Insert(ctx, event)
}

func (s *AuditService) Recent(ctx context.Context) ([]model.AuditEvent, error) {
	return s.repo.ListRecent(ctx, defaultListLimit)
}
