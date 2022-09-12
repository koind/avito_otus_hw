package sqlstorage

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v4"
	sqlstorage "github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v3"
)

const DefaultConfigFile = "../../../configs/calendar_config_test.yaml"

func TestStorage(t *testing.T) { //nolint:funlen,gocognit,nolintlint
	if _, err := os.Stat(DefaultConfigFile); errors.Is(err, os.ErrNotExist) {
		t.Skip(DefaultConfigFile + " file does not exists")
	}

	configContent, _ := os.ReadFile(DefaultConfigFile)
	var config struct {
		Storage struct {
			Dsn string
		}
	}

	err := yaml.Unmarshal(configContent, config)
	if err != nil {
		t.Fatal("Failed to unmarshal config", err)
	}

	ctx := context.Background()
	storage := New(ctx, config.Storage.Dsn)
	if _, err := storage.Connect(ctx); err != nil {
		t.Fatal("Failed to connect to DB server", err)
	}

	t.Run("test SQLStorage CRUDL", func(t *testing.T) {
		tx, err := storage.conn.BeginTx(ctx, pgx.TxOptions{
			IsoLevel:       pgx.Serializable,
			AccessMode:     pgx.ReadWrite,
			DeferrableMode: pgx.NotDeferrable,
		})
		if err != nil {
			t.Fatal("Failed to connect to DB server", err)
		}

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

		event := sqlstorage.NewEvent(
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

		err = tx.Rollback(ctx)
		if err != nil {
			t.Fatal("Failed to rollback tx", err)
		}
	})

	t.Run("test notify list", func(t *testing.T) {
		tx, err := storage.conn.BeginTx(ctx, pgx.TxOptions{
			IsoLevel:       pgx.Serializable,
			AccessMode:     pgx.ReadWrite,
			DeferrableMode: pgx.NotDeferrable,
		})
		if err != nil {
			t.Fatal("Failed to connect to DB server", err)
		}

		userID := uuid.New()
		events := []sqlstorage.Event{
			{
				ID:          parseUUID(t, "4927aa58-a175-429a-a125-c04765597150"),
				StartedAt:   parseDate(t, "2022-04-03 11:59:59"),
				Notify:      parseDate(t, "2022-04-03 11:59:59"),
				UserID:      userID,
				Title:       "Title 1",
				Description: "Desc 1",
				FinishedAt:  parseDate(t, "2022-04-10 12:00:00"),
			},
			{
				ID:          parseUUID(t, "4927aa58-a175-429a-a125-c04765597151"),
				StartedAt:   parseDate(t, "2022-04-03 12:00:00"),
				Notify:      parseDate(t, "2022-04-03 12:00:00"),
				UserID:      userID,
				Title:       "Title 2",
				Description: "Desc 2",
				FinishedAt:  parseDate(t, "2022-04-10 12:00:00"),
			},
			{
				ID:          parseUUID(t, "4927aa58-a175-429a-a125-c04765597152"),
				StartedAt:   parseDate(t, "2022-04-04 12:00:00"),
				Notify:      parseDate(t, "2022-04-03 12:00:00"),
				UserID:      userID,
				Title:       "Title 3",
				Description: "Desc 3",
				FinishedAt:  parseDate(t, "2022-04-10 12:00:00"),
			},
			{
				ID:          parseUUID(t, "4927aa58-a175-429a-a125-c04765597153"),
				StartedAt:   parseDate(t, "2022-04-05 12:00:01"),
				Notify:      parseDate(t, "2022-04-04 11:59:01"),
				UserID:      userID,
				Title:       "Title 4",
				Description: "Desc 4",
				FinishedAt:  parseDate(t, "2022-04-10 12:00:00"),
			},
		}

		for _, e := range events {
			_ = storage.Create(e)
		}

		readyEvents, err := storage.GetEventsReadyToNotify(parseDate(t, "2022-04-03 12:00:00"))
		require.Nil(t, err)

		ids := extractEventIDs(readyEvents)
		idsExpected := []string{
			"4927aa58-a175-429a-a125-c04765597150",
			"4927aa58-a175-429a-a125-c04765597151",
			"4927aa58-a175-429a-a125-c04765597152",
		}
		require.Equal(t, idsExpected, ids)

		err = tx.Rollback(ctx)
		if err != nil {
			t.Fatal("Failed to rollback tx", err)
		}
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

func extractEventIDs(events []sqlstorage.Event) []string {
	res := make([]string, 0, len(events))
	for _, e := range events {
		res = append(res, e.ID.String())
	}

	return res

}
