package service

import (
	"context"

	"github.com/koma2211/ayan-capital_task/internal/entities"
	"github.com/koma2211/ayan-capital_task/internal/repository"
	cacherepository "github.com/koma2211/ayan-capital_task/internal/repository/cache_repository"
)

type Eventer interface {
	AddEvents(ctx context.Context, events []entities.Event) error
	NotifyAllEvents(ctx context.Context) error 
}

type Service struct {
	Eventer
}

func NewService(
	repo *repository.Repository,
	cacheRepo *cacherepository.CacheRepository,
) *Service {
	return &Service{
		Eventer: NewEventService(repo.Eventer, cacheRepo.Eventer),
	}
}
