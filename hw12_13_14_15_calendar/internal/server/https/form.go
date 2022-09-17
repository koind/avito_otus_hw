package https

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/domain/entity"
)

type EventForm struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	Title       string `json:"title"`
	StartedAt   string `json:"startedAt"`
	FinishedAt  string `json:"finishedAt"`
	Description string `json:"description"`
	NotifyAt    string `json:"notifyAt"`
}

type Error struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func (d *EventForm) ToEntity() (*entity.Event, error) {
	id, err := uuid.Parse(d.ID)
	if err != nil {
		return nil, fmt.Errorf("ID must be UUID, your value: %s, %w", d.ID, err)
	}

	userID, err := uuid.Parse(d.UserID)
	if err != nil {
		return nil, fmt.Errorf("userID must be UUID, your value: %s, %w", d.UserID, err)
	}

	startedAt, err := time.Parse("2006-01-02 15:04:00", d.StartedAt)
	if err != nil {
		return nil, fmt.Errorf("startedAt format must be 'yyyy-mm-dd hh:mm:ss', your value: %s, %w", d.StartedAt, err)
	}

	finishedAt, err := time.Parse("2006-01-02 15:04:00", d.FinishedAt)
	if err != nil {
		return nil, fmt.Errorf("finishedAt format must be 'yyyy-mm-dd hh:mm:ss', your value: %s, %w", d.FinishedAt, err)
	}

	notifyAt, err := time.Parse("2006-01-02 15:04:00", d.NotifyAt)
	if err != nil {
		return nil, fmt.Errorf("notifyAt format must be 'yyyy-mm-dd hh:mm:ss', your value: %s, %w", d.NotifyAt, err)
	}

	event := entity.NewEvent(userID, d.Title, startedAt, finishedAt, d.Description, notifyAt)
	event.ID = id

	return event, nil
}
