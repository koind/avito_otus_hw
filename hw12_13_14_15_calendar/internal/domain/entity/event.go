package entity

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/domain/presenter"
)

var (
	ErrDuplicateEvent = errors.New("error duplicate event")
	ErrNotExistEvent  = errors.New("event not exist")
)

type Event struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	StartedAt   time.Time
	FinishedAt  time.Time
	Description string
	NotifyAt    time.Time
}

func NewEvent(
	userID uuid.UUID,
	title string,
	startedAt time.Time,
	finishedAt time.Time,
	description string,
	notifyAt time.Time,
) *Event {
	return &Event{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       title,
		StartedAt:   startedAt,
		FinishedAt:  finishedAt,
		Description: description,
		NotifyAt:    notifyAt,
	}
}

func (e *Event) ToPresenter() presenter.Event {
	return presenter.Event{
		ID:          e.ID.String(),
		UserID:      e.UserID.String(),
		Title:       e.Title,
		StartedAt:   e.StartedAt.Format(time.RFC3339),
		FinishedAt:  e.FinishedAt.Format(time.RFC3339),
		Description: e.Description,
		NotifyAt:    e.NotifyAt.Format(time.RFC3339),
	}
}

type Notification struct {
	EventID  uuid.UUID
	UserID   uuid.UUID
	Title    string
	DateTime time.Time
}

func (n Notification) String() string {
	var builder strings.Builder

	builder.WriteString(
		fmt.Sprintf(
			"New notification: %s at %s",
			n.Title,
			n.DateTime.Format(time.RFC3339),
		),
	)

	return builder.String()
}
