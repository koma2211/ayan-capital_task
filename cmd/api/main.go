package main

import (
	"log"

	"github.com/koma2211/ayan-capital_task/internal/config"
	"github.com/koma2211/ayan-capital_task/internal/handler"
	"github.com/koma2211/ayan-capital_task/internal/repository"
	cacherepository "github.com/koma2211/ayan-capital_task/internal/repository/cache_repository"
	"github.com/koma2211/ayan-capital_task/internal/scheduler"
	"github.com/koma2211/ayan-capital_task/internal/server"
	"github.com/koma2211/ayan-capital_task/internal/service"
	"github.com/koma2211/ayan-capital_task/pkg/cache/redis"
	"github.com/koma2211/ayan-capital_task/pkg/database/migrate"
	"github.com/koma2211/ayan-capital_task/pkg/database/postgres"
)

func main() {
	cfg := config.MustLoad()

	if err := migrate.MigrateUp(cfg.MigrateSource, cfg.DBSoruce); err != nil {
		log.Fatalf("error when to init migration: %s", err.Error())
	}

	db, err := postgres.DBConn(cfg.DBSoruce)
	if err != nil {
		log.Fatalf("error when to connect postgres-db: %s", err.Error())
	}

	cache, err := redis.CacheConn(cfg.RedisSource)
	if err != nil {
		log.Fatalf("eror when  to connect redis: %s", err.Error())
	}

	cacheRepo := cacherepository.NewCacheRepository(cache, cfg.CacheTTL)

	repos := repository.NewRepository(db)
	serv := service.NewService(repos, cacheRepo)
	handlers := handler.NewHandler(serv)

	schedulers := scheduler.NewJobSheduler(serv)

	server := server.SetupServer(handlers, &cfg.HTTPServer, db, schedulers)

	if err := server.Run(); err != nil {
		log.Fatalf("error when to run server: %s", err.Error())
	}
}
