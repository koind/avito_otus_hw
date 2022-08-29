package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/broker/rabbitmq"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/sender_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	configuration, err := config.LoadSender(configFile)
	if err != nil {
		log.Fatalf("Error read configuration: %s", err)
	}

	logg, err := logger.New(configuration.Logger)
	if err != nil {
		log.Fatalf("error create logger: %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	rabbitClient, err := rabbitmq.New(
		ctx,
		configuration.Rabbit.Dsn,
		configuration.Rabbit.Exchange,
		configuration.Rabbit.Queue,
		logg)
	if err != nil {
		cancel()
		log.Fatalf("error create rabbit client: %s", err) //nolint:gocritic
	}

	sender := app.NewSender(rabbitClient, logg)
	sender.Run()

	<-ctx.Done()
}
