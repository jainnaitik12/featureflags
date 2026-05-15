package service

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"

	"featureflags/sdk-service/internal/metrics"
	"featureflags/sdk-service/internal/model"
	"featureflags/sdk-service/internal/repository"
)

var ErrFlagNotFound = errors.New("flag not found")

// EvalService serves read-only flag evaluations for external apps.
type EvalService struct {
	repo *repository.FlagsRedisRead
}

func NewEvalService(repo *repository.FlagsRedisRead) *EvalService {
	return &EvalService{repo: repo}
}

func (s *EvalService) List(ctx context.Context) ([]model.Flag, error) {
	flags, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, f := range flags {
		metrics.FlagEvaluations.WithLabelValues(f.Name).Inc()
	}
	return flags, nil
}

func (s *EvalService) Get(ctx context.Context, name string) (model.FlagEvaluation, error) {
	flag, err := s.repo.Get(ctx, name)
	if errors.Is(err, redis.Nil) {
		return model.FlagEvaluation{}, ErrFlagNotFound
	}
	if err != nil {
		return model.FlagEvaluation{}, err
	}
	metrics.FlagEvaluations.WithLabelValues(name).Inc()
	return model.FlagEvaluation{Flag: name, Enabled: flag.Enabled}, nil
}
