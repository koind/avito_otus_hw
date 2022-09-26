package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/server/grpcs"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/server/https"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/storage"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/calendar_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := config.LoadCalendar(configFile)
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

	calendar := app.New(logg, storage.New(ctx, cfg.Storage))

	// gRPC
	serverGrpc := grpcs.NewServer(logg, calendar, cfg.GRPC.Host, cfg.GRPC.Port)

	go func() {
		if err := serverGrpc.Start(); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
		}
	}()

	go func() {
		<-ctx.Done()
		serverGrpc.Stop()
	}()

	// HTTP
	server := https.NewServer(logg, calendar, cfg.HTTP.Host, cfg.HTTP.Port)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop https server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start https server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
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
