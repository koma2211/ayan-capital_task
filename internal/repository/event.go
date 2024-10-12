package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/koma2211/ayan-capital_task/internal/entities"
)

const (
	prepareEvent = "addEvent"
	cursorEvent  = "eventCursor"
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

func (er *EventRepository) DeclareEventCursor(ctx context.Context) (pgx.Tx, error) {
	tx, err := er.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("DECLARE %s CURSOR FOR SELECT session_id, order_type, card, event_date, website_url FROM %s WHERE DATE(event_date) = CURRENT_DATE", cursorEvent, eventsTable)
	_, err = tx.Exec(ctx, query)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (er *EventRepository) CloseEventCursor(ctx context.Context, tx pgx.Tx) error {
	query := fmt.Sprintf("CLOSE %s", cursorEvent)
	_, err := er.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (er *EventRepository) GetEventByCursor(ctx context.Context, tx pgx.Tx) (entities.Event, error) {
	var event entities.Event
	query := fmt.Sprintf("FETCH NEXT FROM %s", cursorEvent)
	err := tx.QueryRow(ctx, query).Scan(&event.SessionID, &event.OrderType, &event.Card, &event.Date, &event.WebSiteURL)
	if err != nil {
		return entities.Event{}, err
	}
	
	return event, nil
}
