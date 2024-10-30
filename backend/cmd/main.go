package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/john-vh/college_testing/backend/cmd/api"
	"github.com/john-vh/college_testing/backend/db"
	"github.com/john-vh/college_testing/backend/env"
	"github.com/redis/go-redis/v9"
)

var EXT_ENVIRONMENT string = env.DEV

var logLevels = map[string]slog.Level{
	"dev": slog.LevelDebug,
}

func main() {
	cfg := env.GetConfig(EXT_ENVIRONMENT)

	slog.SetLogLoggerLevel(logLevels[EXT_ENVIRONMENT])

	databaseURL := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", cfg.POSTGRES_USER, cfg.POSTGRES_PASSWORD, cfg.POSTGRES_HOST, cfg.POSTGRES_PORT, cfg.POSTGRES_DB)
	pgConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	pg, err := db.NewPgxStorage(context.Background(), pgConfig)
	if err != nil {
		log.Fatal(err)
	}

	opts := redis.Options{
		Addr:     cfg.REDIS_ADDR,
		Password: "",
		DB:       0,
	}
	redis := redis.NewClient(&opts)

	gob.Register(uuid.UUID{})
	server := api.NewAPIServer(":8080", cfg, pg, redis)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
