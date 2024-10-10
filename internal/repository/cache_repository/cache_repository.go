package cacherepository

import (
	"context"
	"time"

	"github.com/koma2211/ayan-capital_task/internal/entities"
	"github.com/redis/go-redis/v9"
)

type Eventer interface {
	AddEvents(ctx context.Context, events []entities.Event) error
}

type CacheRepository struct {
	Eventer
}

func NewCacheRepository(
	cache *redis.Client,
	cacheTTL time.Duration,
) *CacheRepository {
	return &CacheRepository{
		Eventer: NewEventCacheRepository(cache, cacheTTL),
	}
}
