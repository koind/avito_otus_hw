package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/clock"
	internalconfig "github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/mq"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/storage/production"
)

func main() {
	config, err := internalconfig.LoadSchedulerConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}

	logg, err := internallogger.New(config.Logger)
	if err != nil {
		log.Fatalf("Failed to create logger: %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage, err := production.CreateStorage(ctx, config.Storage)
	if err != nil {
		cancel()
		log.Fatalf("Failed to create storage: %s", err) //nolint:gocritic
	}

	notificationReceiver, err := mq.NewRabbit(ctx, config.Rabbit.Dsn, config.Rabbit.Exchange, config.Rabbit.Queue, logg)
	if err != nil {
		cancel()
		log.Fatalf("Failed to create NotificationSender (rabbit): %s", err)
	}

	systemClock := clock.NewSystemClock()

	scheduler := app.NewAppScheduler(storage.(app.EventSource), notificationReceiver, systemClock, logg)

	timer := time.Tick(time.Second)
	timerHour := time.Tick(time.Hour)

	go func() {
		for {
			select {
			case <-timer:
				err := scheduler.Notify()
				if err != nil {
					logg.Error("Failed to Notify: %s", err)
				}
			case <-timerHour:
				err := scheduler.RemoveOneYearOld()
				if err != nil {
					logg.Error("Failed to Notify: %s", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	logg.Info("Calendar Scheduler Started!")

	<-ctx.Done()
}
