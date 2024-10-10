package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/koma2211/ayan-capital_task/internal/config"
	"github.com/koma2211/ayan-capital_task/internal/handler"
)

type Server struct {
	httpServer *http.Server
	db         *pgx.Conn
}

func SetupServer(
	handler *handler.Handler,
	conf *config.HTTPServer,
	db *pgx.Conn,
) *Server {
	httpServer := &http.Server{
		Addr:           conf.Address,
		Handler:        handler.Init(),
		ReadTimeout:    conf.ReadTimeOut,
		WriteTimeout:   conf.WriteTimeOut,
		IdleTimeout:    conf.IdleTimeout,
		MaxHeaderBytes: conf.MaxHeaderBytes,
	}

	return &Server{httpServer: httpServer, db: db}
}

func (s *Server) Run() error {
	// Close db connection...
	defer s.db.Close(context.Background())

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return s.httpServer.Shutdown(ctx)
}
