package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"featureflags/flags-service/internal/model"
	"featureflags/flags-service/internal/repository"
)

var (
	ErrFlagNotFound = errors.New("flag not found")
	ErrInvalidName  = errors.New("name is required")
)

// FlagService contains business rules for feature flags.
type FlagService struct {
	repo *repository.FlagsRedis
}

func NewFlagService(repo *repository.FlagsRedis) *FlagService {
	return &FlagService{repo: repo}
}

func (s *FlagService) List(ctx context.Context) ([]model.Flag, error) {
	return s.repo.ListAll(ctx)
}

func (s *FlagService) Get(ctx context.Context, name string) (model.Flag, error) {
	flag, err := s.repo.Get(ctx, name)
	if errors.Is(err, redis.Nil) {
		return model.Flag{}, ErrFlagNotFound
	}
	if err != nil {
		return model.Flag{}, err
	}
	return flag, nil
}

func (s *FlagService) Create(ctx context.Context, flag model.Flag) (model.Flag, error) {
	flag.Name = strings.TrimSpace(flag.Name)
	if flag.Name == "" {
		return model.Flag{}, ErrInvalidName
	}
	if err := s.repo.Set(ctx, flag); err != nil {
		return model.Flag{}, err
	}
	return flag, nil
}

func (s *FlagService) Toggle(ctx context.Context, name string) (model.Flag, error) {
	current, err := s.repo.Get(ctx, name)
	if errors.Is(err, redis.Nil) {
		return model.Flag{}, ErrFlagNotFound
	}
	if err != nil {
		return model.Flag{}, err
	}

	newVal := !current.Enabled
	if err := s.repo.Set(ctx, model.Flag{Name: name, Enabled: newVal}); err != nil {
		return model.Flag{}, err
	}

	event := model.ToggleEvent{
		FlagName:  name,
		Action:    "toggle",
		OldValue:  current.Enabled,
		NewValue:  newVal,
		ChangedAt: time.Now().UTC(),
	}
	if err := s.repo.PublishToggle(ctx, event); err != nil {
		log.Printf("failed to publish toggle event: %v", err)
	}

	return model.Flag{Name: name, Enabled: newVal}, nil
}

func (s *FlagService) Delete(ctx context.Context, name string) error {
	ok, err := s.repo.Delete(ctx, name)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFlagNotFound
	}
	return nil
}
