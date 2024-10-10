package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/koma2211/ayan-capital_task/internal/entities"
)

const (
	prepareEvent = "addEvent"
)

type EventRepository struct {
	db *pgx.Conn
}

func NewEventRepository(db *pgx.Conn) *EventRepository {
	return &EventRepository{db: db}
}

func (er *EventRepository) AddEvents(ctx context.Context, events []entities.Event) error {
	tx, err := er.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO %s (session_id, order_type, card, event_date, website_url) VALUES ($1, $2, $3, $4, $5)", eventsTable)

	stmt, err := tx.Prepare(ctx, prepareEvent, query)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	for i := 0; i < len(events); i++ {
		_, err = tx.Exec(ctx, stmt.Name, events[i].SessionID, events[i].OrderType, events[i].Card, events[i].Date, events[i].WebSiteURL)
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return err
			}
			
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
