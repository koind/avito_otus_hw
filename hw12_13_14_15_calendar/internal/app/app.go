package app

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/domain/entity"
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Error(format string, params ...interface{})
	Info(format string, params ...interface{})
}

type Storage interface {
	SelectOne(id uuid.UUID) (*entity.Event, error)
	Select() ([]entity.Event, error)
	Insert(e entity.Event) error
	Update(e entity.Event) error
	Delete(id uuid.UUID) error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event entity.Event) error {
	_, err := a.Storage.SelectOne(event.ID)
	if err != nil && !errors.Is(err, entity.ErrNotExistEvent) {
		a.Logger.Error("Error CreateEvent: Event exists: %s", err)

		return err
	}

	if err := a.Storage.Insert(event); err != nil {
		a.Logger.Error("Error CreateEvent. Can't create event: %s", err)

		return err
	}

	return nil
}

func (a *App) UpdateEvent(ctx context.Context, event entity.Event) error {
	var dbEvent *entity.Event
	var err error

	if dbEvent, err = a.Storage.SelectOne(event.ID); err != nil || dbEvent == nil {
		a.Logger.Error("Error UpdateEvent: Event is absent: %s", err)
		return err
	}

	if err = a.Storage.Update(event); err != nil {
		a.Logger.Error("Error UpdateEvent. Can't update event: %s", err)
		return err
	}

	return nil
}

func (a *App) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	var dbEvent *entity.Event
	var err error

	a.Logger.Info("DeleteEvent %s", id)

	if dbEvent, err = a.Storage.SelectOne(id); err != nil || dbEvent == nil {
		a.Logger.Error("Error DeleteEvent. Event is absent: %s", err)
		return err
	}

	if err = a.Storage.Delete(id); err != nil {
		a.Logger.Error("Error DeleteEvent. Can't delete event: %s", err)
		return err
	}

	return nil
}

func (a *App) GetEvents(ctx context.Context) ([]entity.Event, error) {
	return a.Storage.Select()
}

func (a *App) GetDayIntervalEvents(ctx context.Context, day time.Time, interval time.Duration) ([]entity.Event, error) { // nolint:lll
	events := make([]entity.Event, 0)
	day = day.Truncate(time.Minute * 1440)

	items, err := a.Storage.Select()
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		diff := item.StartedAt.Sub(day)
		if diff >= 0 && diff < interval {
			events = append(events, item)
		}
	}

	return events, nil
}

func (a *App) GetDayEvents(ctx context.Context, day time.Time) ([]entity.Event, error) {
	return a.GetDayIntervalEvents(ctx, day, day.AddDate(0, 0, 1).Sub(day))
}

func (a *App) GetWeekEvents(ctx context.Context, day time.Time) ([]entity.Event, error) {
	return a.GetDayIntervalEvents(ctx, day, day.AddDate(0, 0, 7).Sub(day))
}

func (a *App) GetMonthEvents(ctx context.Context, day time.Time) ([]entity.Event, error) {
	return a.GetDayIntervalEvents(ctx, day, day.AddDate(0, 1, 0).Sub(day))
}
