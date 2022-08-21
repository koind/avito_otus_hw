package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/storage/postgres"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := config.Load(configFile)
	if err != nil {
		log.Fatalf("Error read configuration: %s", err)
	}

	logg, err := logger.New(cfg.Logger)
	if err != nil {
		log.Fatalf("Error create logger: %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage := NewStorage(ctx, *cfg)
	calendar := app.New(logg, storage)
	server := http.NewServer(logg, calendar, cfg.HTTP.Host, cfg.HTTP.Port)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func NewStorage(ctx context.Context, cfg config.Config) app.Storage {
	var storage app.Storage

	switch cfg.Storage.Type {
	case "memory":
		storage = memory.New()
	case "postgres":
		storage = postgres.New(ctx, cfg.Storage.Dsn).Connect(ctx)
	default:
		log.Fatalln("Unknown type of storage: " + cfg.Storage.Type)
	}

	return storage
}

func printVersion() {
	if err := json.NewEncoder(os.Stdout).Encode(struct {
		Release   string
		BuildDate string
		GitHash   string
	}{
		Release:   release,
		BuildDate: buildDate,
		GitHash:   gitHash,
	}); err != nil {
		fmt.Printf("error while decode version info: %v\n", err)
	}
}
