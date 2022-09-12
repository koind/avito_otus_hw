package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Debug(format string, params ...interface{})
	Info(format string, params ...interface{})
	Warn(format string, params ...interface{})
	Error(format string, params ...interface{})
	LogHTTPRequest(r *http.Request, code, length int)
}

type Storage interface {
	Create(e storage.Event) error
	Update(e storage.Event) error
	Delete(id uuid.UUID) error
	Find(id uuid.UUID) (*storage.Event, error)
	FindAll() ([]storage.Event, error)
	FindAtDay(day time.Time) ([]storage.Event, error)
	FindAtWeek(dayStart time.Time) ([]storage.Event, error)
	FindAtMonth(dayStart time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {
	var existingEvent *storage.Event
	var err error

	a.Logger.Debug("App.CreateEvent %s", event.ID)

	if existingEvent, err = a.Storage.Find(event.ID); err != nil {
		a.Logger.Error("App.CreateEvent error: find existing event error: %s", err)
		return err
	}

	if existingEvent != nil {
		a.Logger.Warn("App.CreateEvent error: event with ID %s already exists", event.ID)
		return fmt.Errorf("event with ID %s already exists", event.ID)
	}

	if err = a.Storage.Create(event); err != nil {
		a.Logger.Error("App.CreateEvent error: %s", err)
		return err
	}

	return nil
}

func (a *App) UpdateEvent(ctx context.Context, event storage.Event) error {
	var existingEvent *storage.Event
	var err error

	a.Logger.Debug("App.UpdateEvent %s", event.ID)

	if existingEvent, err = a.Storage.Find(event.ID); err != nil {
		a.Logger.Error("App.UpdateEvent error: find existing event error: %s", err)
		return err
	}

	if existingEvent == nil {
		a.Logger.Warn("App.UpdateEvent error: event with ID %s not found", event.ID)
		return fmt.Errorf("event with ID %s not found", event.ID)
	}

	if err = a.Storage.Update(event); err != nil {
		a.Logger.Error("App.UpdateEvent error: %s", err)
		return err
	}

	return nil
}

func (a *App) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	var existingEvent *storage.Event
	var err error

	a.Logger.Debug("App.DeleteEvent %s", id)

	if existingEvent, err = a.Storage.Find(id); err != nil {
		a.Logger.Error("App.DeleteEvent error: find existing event error: %s", err)
		return err
	}

	if existingEvent == nil {
		a.Logger.Warn("App.UpdateEvent error: event with ID %s not found", id)
		return fmt.Errorf("event with ID %s not found", id)
	}

	if err = a.Storage.Delete(id); err != nil {
		a.Logger.Error("App.DeleteEvent error: %s", err)
		return err
	}

	return nil
}

func (a *App) GetEvents(ctx context.Context) ([]storage.Event, error) {
	return a.Storage.FindAll()
}

func (a *App) GetEventsDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	return a.Storage.FindAtDay(day)
}

func (a *App) GetEventsMonth(ctx context.Context, dayStart time.Time) ([]storage.Event, error) {
	return a.Storage.FindAtMonth(dayStart)
}

func (a *App) GetEventsStartedIn(ctx context.Context, day time.Time, interval time.Duration) ([]storage.Event, error) {
	a.Logger.Debug("App.GetEventsStartedIn: day: %s, interval: %s", day, interval)

	var events []storage.Event
	day = day.Truncate(time.Minute * 1440)

	items, err := a.Storage.FindAll()
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		diff := item.StartedAt.Sub(day)
		if diff >= 0 && diff <= interval {
			events = append(events, item)
		}
	}

	return events, nil
}
