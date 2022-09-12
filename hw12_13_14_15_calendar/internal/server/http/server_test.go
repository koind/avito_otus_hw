package internalhttp

import (
	"bytes"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestHttpServerHelloWorld(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	httpHandlers := NewRouter(createApp(t))
	httpHandlers.ServeHTTP(w, req)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, "Hello, world!\n", string(body))
}

func TestHttpServerEventsCrud(t *testing.T) {
	body := bytes.NewBufferString(`{
		"id": "4927aa58-a175-429a-a125-c04765597152",
		"title": "Event Title 1",
		"startedAt": "2022-03-20 12:30:00",
		"finishedAt": "2022-03-21 12:30:00",
		"description": "Event Description 1",
		"notify": "2022-03-19 12:30:00",
		"userId": "b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee"
	}`)
	req := httptest.NewRequest("POST", "/events", body)
	w := httptest.NewRecorder()

	httpHandlers := NewRouter(createApp(t))
	httpHandlers.ServeHTTP(w, req)

	resp := w.Result()
	respBody, _ := io.ReadAll(resp.Body)
	respExpected := `{"id":"4927aa58-a175-429a-a125-c04765597152","title":"Event Title 1","startedAt":"2022-03-20 12:30:00","finishedAt":"2022-03-21 12:30:00","description":"Event Description 1","userId":"b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee","notify":"2022-03-19 12:30:00"}` // nolint:lll
	require.Equal(t, respExpected, string(respBody))

	body = bytes.NewBufferString(`{
		"id": "4927aa58-a175-429a-a125-c04765597152",
		"title": "Event Title 2",
		"startedAt": "2022-03-20 12:30:00",
		"finishedAt": "2022-03-21 12:30:00",
		"description": "Event Description 2",
		"notify": "2022-03-19 12:30:00",
		"userId": "b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee"
	}`)
	req = httptest.NewRequest("PUT", "/events/4927aa58-a175-429a-a125-c04765597152", body)
	w = httptest.NewRecorder()

	httpHandlers.ServeHTTP(w, req)

	resp = w.Result()
	respBody, _ = io.ReadAll(resp.Body)
	respExpected = `{"id":"4927aa58-a175-429a-a125-c04765597152","title":"Event Title 2","startedAt":"2022-03-20 12:30:00","finishedAt":"2022-03-21 12:30:00","description":"Event Description 2","userId":"b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee","notify":"2022-03-19 12:30:00"}` // nolint:lll
	require.Equal(t, respExpected, string(respBody))
}

func createApp(t *testing.T) *app.App {
	t.Helper()
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

	return app.New(logg, inMemoryStorage)
}
