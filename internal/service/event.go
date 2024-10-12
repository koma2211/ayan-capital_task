package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
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

func (es *EventService) NotifyAllEvents(ctx context.Context) error {
	tx, err := es.eventRepo.DeclareEventCursor(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	for {
		event, err := es.eventRepo.GetEventByCursor(ctx, tx)
		if err != nil {
			if err := es.eventRepo.CloseEventCursor(ctx, tx); err != nil {
				return err
			}

			if err := tx.Commit(ctx); err != nil {
				return err
			}

			if err == pgx.ErrNoRows {
				return nil
			}

			return err
		}

		fmt.Printf("session-id:%s\n", event.SessionID)
	}
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
