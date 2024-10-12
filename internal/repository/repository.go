package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/koma2211/ayan-capital_task/internal/entities"
)

const (
	eventsTable = "events"
)

type Eventer interface {
	AddEvents(ctx context.Context, event []entities.Event) error
	DeclareEventCursor(ctx context.Context) (pgx.Tx, error)
	CloseEventCursor(ctx context.Context, tx pgx.Tx) error
	GetEventByCursor(ctx context.Context, tx pgx.Tx) (entities.Event, error)
}

type Repository struct {
	Eventer
}

func NewRepository(
	db *pgx.Conn,
) *Repository {
	return &Repository{
		Eventer: NewEventRepository(db),
	}
}
