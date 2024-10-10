package service

import (
	"context"

	"github.com/koma2211/ayan-capital_task/internal/entities"
	"github.com/koma2211/ayan-capital_task/internal/repository"
	cacherepository "github.com/koma2211/ayan-capital_task/internal/repository/cache_repository"
)

type EventService struct {
	eventRepo      repository.Eventer
	eventCacheRepo cacherepository.Eventer
}

func NewEventService(
	eventRepo repository.Eventer,
	eventCacheRepo cacherepository.Eventer,
) *EventService {
	return &EventService{eventRepo: eventRepo, eventCacheRepo: eventCacheRepo}
}

func (es *EventService) AddEvents(ctx context.Context, events []entities.Event) error {
	// Sorry for using goroutine.
	// I wanted to use RabbitMQ, but it takes a long time:(
	
	errRepoEventCh := make(chan error, 1)
	go func() {
		err := es.eventRepo.AddEvents(ctx, events)
		if err != nil {
			errRepoEventCh <- err
			return
		}
		errRepoEventCh <- nil
	}()
	

	errCacheRepoEventCh := make(chan error, 1)
	go func() {
		err := es.eventCacheRepo.AddEvents(ctx, events)
		if err != nil {
			errCacheRepoEventCh <- err
			return
		}
		errCacheRepoEventCh <- nil
	}()

	errRepo, errCacheRepo := <-errRepoEventCh, <-errCacheRepoEventCh
	
	if errRepo != nil {
		return errRepo
	}

	if errCacheRepo != nil {
		return errCacheRepo
	}

	return nil
}
