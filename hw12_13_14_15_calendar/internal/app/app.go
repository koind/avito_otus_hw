package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/domain/entity"
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Error(format string, params ...interface{})
}

type Storage interface {
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

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
