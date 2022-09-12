package storage

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEventAlreadyExists = errors.New("event already exists")
	ErrEventDoesNotExists = errors.New("event does not exist")
	ErrorEventNotFound    = errors.New("error: event not found")
)

type Event struct {
	ID          uuid.UUID
	Title       string
	StartedAt   time.Time
	FinishedAt  time.Time
	Description string
	UserID      uuid.UUID
	Notify      time.Time
}

func NewEvent(
	title string,
	startedAt time.Time,
	finishedAt time.Time,
	description string,
	userID uuid.UUID,
	notify time.Time,
) *Event {
	return &Event{
		ID:          uuid.New(),
		Title:       title,
		StartedAt:   startedAt,
		FinishedAt:  finishedAt,
		Description: description,
		UserID:      userID,
		Notify:      notify,
	}
}
