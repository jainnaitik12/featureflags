package repository

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"

	"featureflags/flags-service/internal/model"
)

const (
	flagsHashKey      = "flags"
	flagChangesPubSub = "flag-changes"
)

// FlagsRedis persists flags in a Redis hash and publishes toggle events.
type FlagsRedis struct {
	rdb *redis.Client
}

func NewFlagsRedis(rdb *redis.Client) *FlagsRedis {
	return &FlagsRedis{rdb: rdb}
}

func (s *FlagsRedis) ListAll(ctx context.Context) ([]model.Flag, error) {
	data, err := s.rdb.HGetAll(ctx, flagsHashKey).Result()
	if err != nil {
		return nil, err
	}
	flags := make([]model.Flag, 0, len(data))
	for name, val := range data {
		flags = append(flags, model.Flag{Name: name, Enabled: parseEnabled(val)})
	}
	return flags, nil
}

func (s *FlagsRedis) Get(ctx context.Context, name string) (model.Flag, error) {
	val, err := s.rdb.HGet(ctx, flagsHashKey, name).Result()
	if err != nil {
		return model.Flag{}, err
	}
	return model.Flag{Name: name, Enabled: parseEnabled(val)}, nil
}

func (s *FlagsRedis) Set(ctx context.Context, flag model.Flag) error {
	val := "false"
	if flag.Enabled {
		val = "true"
	}
	return s.rdb.HSet(ctx, flagsHashKey, flag.Name, val).Err()
}

func parseEnabled(val string) bool {
	return val == "true" || val == "1"
}

func (s *FlagsRedis) Delete(ctx context.Context, name string) (deleted bool, err error) {
	n, err := s.rdb.HDel(ctx, flagsHashKey, name).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (s *FlagsRedis) PublishToggle(ctx context.Context, event model.ToggleEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return s.rdb.Publish(ctx, flagChangesPubSub, payload).Err()
}
