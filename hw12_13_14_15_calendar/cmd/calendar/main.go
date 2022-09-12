package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	internalapp "github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/app"
	internalconfig "github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/server/http"
	internalstore "github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/storage/production"
)

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := internalconfig.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}

	logg, err := internallogger.New(config.Logger)
	if err != nil {
		log.Fatalf("Failed to create logger: %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage, err := internalstore.CreateStorage(ctx, config.Storage)
	if err != nil {
		cancel()
		log.Fatalf("Failed to create storage: %s", err) //nolint:gocritic
	}

	calendar := internalapp.New(logg, storage)

	serverGrpc := internalgrpc.NewServer(logg, calendar, config.GRPC.Host, config.GRPC.Port)

	go func() {
		if err := serverGrpc.Start(); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
		}
	}()

	go func() {
		<-ctx.Done()
		serverGrpc.Stop()
	}()

	serverHttp := internalhttp.NewServer(logg, calendar, config.HTTP.Host, config.HTTP.Port)

	go func() {
		if err := serverHttp.Start(ctx); err != nil {
			logg.Error("failed to start server: " + err.Error())
			cancel()
		}
	}()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := serverHttp.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	<-ctx.Done()
}
