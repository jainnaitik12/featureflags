package repository

import (
	"context"

	"github.com/redis/go-redis/v9"

	"featureflags/sdk-service/internal/model"
)

const flagsHashKey = "flags"

// FlagsRedisRead provides read-only access to flags in Redis.
type FlagsRedisRead struct {
	rdb *redis.Client
}

func NewFlagsRedisRead(rdb *redis.Client) *FlagsRedisRead {
	return &FlagsRedisRead{rdb: rdb}
}

func (s *FlagsRedisRead) ListAll(ctx context.Context) ([]model.Flag, error) {
	all, err := s.rdb.HGetAll(ctx, flagsHashKey).Result()
	if err != nil {
		return nil, err
	}
	flags := make([]model.Flag, 0, len(all))
	for name, value := range all {
		flags = append(flags, model.Flag{Name: name, Enabled: parseEnabled(value)})
	}
	return flags, nil
}

func (s *FlagsRedisRead) Get(ctx context.Context, name string) (model.Flag, error) {
	val, err := s.rdb.HGet(ctx, flagsHashKey, name).Result()
	if err != nil {
		return model.Flag{}, err
	}
	return model.Flag{Name: name, Enabled: parseEnabled(val)}, nil
}

func parseEnabled(val string) bool {
	return val == "true" || val == "1"
}
