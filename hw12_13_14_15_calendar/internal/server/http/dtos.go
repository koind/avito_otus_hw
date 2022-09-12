package internalhttp

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/storage"
)

type EventDto struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	StartedAt   string `json:"startedAt"`
	FinishedAt  string `json:"finishedAt"`
	Description string `json:"description"`
	UserID      string `json:"userId"`
	Notify      string `json:"notify"`
}

type ErrorDto struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func (e *EventDto) GetModel() (*storage.Event, error) {
	startedAt, err := time.Parse("2006-01-02 15:04:00", e.StartedAt)
	if err != nil {
		return nil, fmt.Errorf("error: StartedAt exprected to be 'yyyy-mm-dd hh:mm:ss', got: %s, %w", e.StartedAt, err)
	}

	finishedAt, err := time.Parse("2006-01-02 15:04:00", e.FinishedAt)
	if err != nil {
		return nil, fmt.Errorf("error: FinishedAt exprected to be 'yyyy-mm-dd hh:mm:ss', got: %s, %w", e.FinishedAt, err)
	}

	notify, err := time.Parse("2006-01-02 15:04:00", e.Notify)
	if err != nil {
		return nil, fmt.Errorf("error: Notify exprected to be 'yyyy-mm-dd hh:mm:ss', got: %s, %w", e.Notify, err)
	}

	id, err := uuid.Parse(e.ID)
	if err != nil {
		return nil, fmt.Errorf("ID exprected to be uuid, got: %s, %w", e.ID, err)
	}

	userID, err := uuid.Parse(e.UserID)
	if err != nil {
		return nil, fmt.Errorf("userID exprected to be uuid, got: %s, %w", e.UserID, err)
	}

	appEvent := storage.NewEvent(e.Title, startedAt, finishedAt, e.Description, userID, notify)
	appEvent.ID = id

	return appEvent, nil
}

func CreateEventDtoFromModel(event storage.Event) EventDto {
	return EventDto{
		ID:          event.ID.String(),
		Title:       event.Title,
		StartedAt:   event.StartedAt.Format(time.RFC3339),
		FinishedAt:  event.FinishedAt.Format(time.RFC3339),
		Description: event.Description,
		UserID:      event.UserID.String(),
		Notify:      event.Notify.Format(time.RFC3339),
	}
}
