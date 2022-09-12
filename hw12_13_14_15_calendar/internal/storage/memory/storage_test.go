package memorystorage

import (
	"testing"
	"time"

	"github.com/google/uuid"
	memorystorage "github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) { //nolint:funlen,gocognit,nolintlint
	storage := New()

	t.Run("common test", func(t *testing.T) {
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

		event := memorystorage.NewEvent(
			"Event title",
			startedAt,
			finishedAt,
			"Event description",
			userID,
			notifyAt,
		)

		err = storage.Create(*event)
		if err != nil {
			t.FailNow()
			return
		}

		saved, err := storage.FindAll()
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, saved, 1)
		require.Equal(t, *event, saved[0])

		event.Title = "New event title"
		event.Description = "New event description"

		saved, err = storage.FindAll()
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, saved, 1)
		require.NotEqual(t, *event, saved[0])
		require.NotEqual(t, event.Title, saved[0].Title)
		require.NotEqual(t, event.Description, saved[0].Description)

		err = storage.Update(*event)
		if err != nil {
			t.FailNow()
			return
		}

		saved, err = storage.FindAll()
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, saved, 1)
		require.Equal(t, *event, saved[0])
		require.Equal(t, event.Title, saved[0].Title)
		require.Equal(t, event.Description, saved[0].Description)

		err = storage.Delete(event.ID)
		if err != nil {
			t.FailNow()
			return
		}

		saved, err = storage.FindAll()
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, saved, 0)
	})

	t.Run("test notify list", func(t *testing.T) {
		events := []memorystorage.Event{
			{
				ID:        parseUUID(t, "4927aa58-a175-429a-a125-c04765597150"),
				StartedAt: parseDate(t, "2022-04-03T11:59:59Z"),
				Notify:    parseDate(t, "2022-04-03T11:59:59Z"),
			},
			{
				ID:        parseUUID(t, "4927aa58-a175-429a-a125-c04765597151"),
				StartedAt: parseDate(t, "2022-04-03T12:00:00Z"),
				Notify:    parseDate(t, "2022-04-03T12:00:00Z"),
			},
			{
				ID:        parseUUID(t, "4927aa58-a175-429a-a125-c04765597152"),
				StartedAt: parseDate(t, "2022-04-04T12:00:00Z"),
				Notify:    parseDate(t, "2022-04-03T12:00:00Z"),
			},
			{
				ID:        parseUUID(t, "4927aa58-a175-429a-a125-c04765597153"),
				StartedAt: parseDate(t, "2022-04-05T12:00:01Z"),
				Notify:    parseDate(t, "2022-04-04T11:59:01Z"),
			},
		}

		for _, e := range events {
			_ = storage.Create(e)
		}

		readyEvents, err := storage.GetEventsReadyToNotify(parseDate(t, "2022-04-03T12:00:00Z"))
		require.Nil(t, err)

		ids := extractEventIDs(readyEvents)
		idsExpected := []string{
			"4927aa58-a175-429a-a125-c04765597150",
			"4927aa58-a175-429a-a125-c04765597151",
			"4927aa58-a175-429a-a125-c04765597152",
		}
		require.Equal(t, idsExpected, ids)
	})
}

func parseUUID(t *testing.T, str string) uuid.UUID {
	t.Helper()
	id, err := uuid.Parse(str)
	if err != nil {
		t.Errorf("failed to parse UUID: %s", err)
	}
	return id
}

func parseDate(t *testing.T, str string) time.Time {
	t.Helper()
	dt, err := time.Parse(time.RFC3339, str)
	if err != nil {
		t.Errorf("failed to parse date: %s", err)
	}
	return dt
}

func extractEventIDs(events []memorystorage.Event) []string {
	res := make([]string, 0, len(events))
	for _, e := range events {
		res = append(res, e.ID.String())
	}

	return res
}
