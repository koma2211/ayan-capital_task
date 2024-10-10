package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func DBConn(dbSource string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), dbSource)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return conn, nil
}
