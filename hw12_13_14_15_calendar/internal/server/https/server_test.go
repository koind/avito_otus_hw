package https

import (
	"bytes"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestHttpServerEvents(t *testing.T) {
	body := bytes.NewBufferString(`{
		"id": "dd983962-b469-11ec-b909-0242ac120002",
		"userId": "fae10b5c-b469-11ec-b909-0242ac120002",
		"title": "Test Title",
		"startedAt": "2022-04-01 00:00:00",
		"finishedAt": "2022-05-01 00:00:00",
		"description": "Test Description",
		"notifyAt": "2022-03-30 00:00:00"
	}`)
	req := httptest.NewRequest("POST", "/events", body)
	w := httptest.NewRecorder()

	application := createApp(t)
	httpHandlers := routes(application)
	httpHandlers.ServeHTTP(w, req)

	resp := w.Result()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		require.NoError(t, err)
		return
	}
	defer resp.Body.Close()

	respExpected := `{"id":"dd983962-b469-11ec-b909-0242ac120002","userId":"fae10b5c-b469-11ec-b909-0242ac120002","title":"Test Title","startedAt":"2022-04-01T00:00:00Z","finishedAt":"2022-05-01T00:00:00Z","description":"Test Description","notifyAt":"2022-03-30T00:00:00Z"}` // nolint:lll
	require.Equal(t, respExpected, string(respBody))

	body = bytes.NewBufferString(`{
		"userId": "e648562c-b46a-11ec-b909-0242ac120002",
		"title": "New Test Title",
		"startedAt": "2023-04-01 00:00:00",
		"finishedAt": "2023-05-01 00:00:00",
		"description": "New Test Description",	
		"notifyAt": "2023-03-30 00:00:00"
	}`)
	req = httptest.NewRequest("PUT", "/events/dd983962-b469-11ec-b909-0242ac120002", body)
	w = httptest.NewRecorder()

	httpHandlers.ServeHTTP(w, req)

	resp = w.Result()
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		require.NoError(t, err)
		return
	}
	defer resp.Body.Close()

	respExpected = `{"id":"dd983962-b469-11ec-b909-0242ac120002","userId":"e648562c-b46a-11ec-b909-0242ac120002","title":"New Test Title","startedAt":"2023-04-01T00:00:00Z","finishedAt":"2023-05-01T00:00:00Z","description":"New Test Description","notifyAt":"2023-03-30T00:00:00Z"}` // nolint:lll
	require.Equal(t, respExpected, string(respBody))
}

func createApp(t *testing.T) *app.App {
	t.Helper()

	logFile, err := os.CreateTemp("", "test.log")
	if err != nil {
		t.Errorf("failed to open test log file: %s", err)
	}

	log, err := logger.New(config.LoggerConf{Level: "info", Filename: logFile.Name()})
	if err != nil {
		t.Errorf("failed to open test log file: %s", err)
	}

	storage := memory.New()

	return app.New(log, storage)
}
