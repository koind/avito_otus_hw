package app

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestAppEventCrud(t *testing.T) {
	logFile, err := os.CreateTemp("", "log")
	if err != nil {
		t.Errorf("failed to open test log file: %s", err)
	}

	logg, err := logger.New(config.LoggerConf{
		Level:    config.Debug,
		Filename: logFile.Name(),
	})
	if err != nil {
		t.Errorf("failed to open test log file: %s", err)
	}

	inMemoryStorage := memorystorage.New()
	ctx := context.Background()
	testApp := New(logg, inMemoryStorage)

	event1 := storage.Event{
		ID:          parseUUID(t, "4927aa58-a175-429a-a125-c04765597152"),
		Title:       "Event Title 1",
		StartedAt:   parseDate(t, "2022-03-20T12:30:00Z"),
		FinishedAt:  parseDate(t, "2022-03-21T12:30:00Z"),
		Description: "Event Description 1",
		UserID:      parseUUID(t, "b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee"),
		Notify:      parseDate(t, "2022-03-19T12:30:00Z"),
	}
	err = testApp.CreateEvent(ctx, event1)
	require.Nil(t, err)

	// +1 week
	event2 := storage.Event{
		ID:          parseUUID(t, "4927aa58-a175-429a-a125-c04765597153"),
		Title:       "Event Title 2",
		StartedAt:   parseDate(t, "2022-03-27T12:30:00Z"),
		FinishedAt:  parseDate(t, "2022-03-28T12:30:00Z"),
		Description: "Event Description 2",
		UserID:      parseUUID(t, "b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee"),
		Notify:      parseDate(t, "2022-03-26T12:30:00Z"),
	}
	err = testApp.CreateEvent(ctx, event2)
	require.Nil(t, err)

	// +1 month
	event3 := storage.Event{
		ID:          parseUUID(t, "4927aa58-a175-429a-a125-c04765597154"),
		Title:       "Event Title 3",
		StartedAt:   parseDate(t, "2022-04-20T12:30:00Z"),
		FinishedAt:  parseDate(t, "2022-04-21T12:30:00Z"),
		Description: "Event Description 3",
		UserID:      parseUUID(t, "b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee"),
		Notify:      parseDate(t, "2022-04-19T12:30:00Z"),
	}
	err = testApp.CreateEvent(ctx, event3)
	require.Nil(t, err)

	// -1 day
	event4 := storage.Event{
		ID:          parseUUID(t, "4927aa58-a175-429a-a125-c04765597155"),
		Title:       "Event Title 4",
		StartedAt:   parseDate(t, "2022-03-19T12:30:00Z"),
		FinishedAt:  parseDate(t, "2022-03-20T12:30:00Z"),
		Description: "Event Description 4",
		UserID:      parseUUID(t, "b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee"),
		Notify:      parseDate(t, "2022-03-18T12:30:00Z"),
	}
	err = testApp.CreateEvent(ctx, event4)
	require.Nil(t, err)

	startedAt := parseDate(t, "2022-03-20T12:30:00Z")

	events, err := testApp.GetEventsMonth(ctx, startedAt)
	require.Nil(t, err)
	require.Len(t, events, 3)
	require.Equal(t, "4927aa58-a175-429a-a125-c04765597152", events[0].ID.String())
	require.Equal(t, "4927aa58-a175-429a-a125-c04765597153", events[1].ID.String())
	require.Equal(t, "4927aa58-a175-429a-a125-c04765597154", events[2].ID.String())
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
