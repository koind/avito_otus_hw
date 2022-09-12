package storage

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestEvent(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		userID := uuid.New()
		startedAt, err := time.Parse("2006-01-02 15:04:05", "2022-03-08 12:00:00")
		if err != nil {
			t.FailNow()
			return
		}
		finishedAt, err := time.Parse("2006-01-02 15:04:05", "2022-03-09 12:00:00")
		if err != nil {
			t.FailNow()
			return
		}
		notifyAt, err := time.Parse("2006-01-02 15:04:05", "2022-03-07 12:00:00")
		if err != nil {
			t.FailNow()
			return
		}

		event := NewEvent(
			"Event title",
			startedAt,
			finishedAt,
			"Event description",
			userID,
			notifyAt,
		)

		require.Len(t, event.ID, 16)
		require.Equal(t, "Event title", event.Title)
		require.Equal(t, startedAt, event.StartedAt)
		require.Equal(t, finishedAt, event.FinishedAt)
		require.Equal(t, "Event description", event.Description)
		require.Equal(t, userID, event.UserID)
		require.Equal(t, notifyAt, event.Notify)
	})
}
