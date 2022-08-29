package storage

import (
	"context"
	"log"

	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/storage/postgres"
)

func New(ctx context.Context, cfg config.StorageConf) app.Storage {
	var storage app.Storage

	switch cfg.Type {
	case "memory":
		storage = memory.New()
	case "postgres":
		storage = postgres.New(ctx, cfg.Dsn).Connect(ctx)
	default:
		log.Fatalln("Unknown type of storage: " + cfg.Type)
	}

	return storage
}
