package cacherepository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/koma2211/ayan-capital_task/internal/entities"
	"github.com/redis/go-redis/v9"
)

type EventCacheRepository struct {
	cache    *redis.Client
	cacheTTL time.Duration
}

func NewEventCacheRepository(
	cache *redis.Client,
	cacheTTL time.Duration,
) *EventCacheRepository {
	return &EventCacheRepository{
		cache:    cache,
		cacheTTL: cacheTTL,
	}
}

func (ecr *EventCacheRepository) AddEvents(ctx context.Context, events []entities.Event) error {
	pipe := ecr.cache.TxPipeline()

	for i := 0; i < len(events); i++ {
		body, err := json.Marshal(events[i])
		if err != nil {
			return err
		}

		if err := pipe.Set(ctx, generateEventKey(events[i].SessionID), body, ecr.cacheTTL).Err(); err != nil {
			return err
		}
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return  err
	}

	return nil
}

func generateEventKey(key string) string {
	return fmt.Sprintf("session-Id:%s", key)
}
