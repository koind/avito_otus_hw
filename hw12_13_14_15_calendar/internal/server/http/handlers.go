package internalhttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/storage"
)

type ServerHandlers struct {
	app *app.App
}

func NewServerHandlers(a *app.App) *ServerHandlers {
	return &ServerHandlers{app: a}
}

func (s *ServerHandlers) HelloWorld(w http.ResponseWriter, r *http.Request) {
	msg := []byte("Hello, world!\n")
	w.WriteHeader(200)
	w.Write(msg)
}

func (s *ServerHandlers) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var dto EventDto
	err := ParseRequest(r, &dto)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	event, err := dto.GetModel()
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	err = s.app.CreateEvent(r.Context(), *event)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err)
		return
	}

	responseData, _ := json.Marshal(dto)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseData)
}

func (s *ServerHandlers) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var dto EventDto
	err := ParseRequest(r, &dto)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	vars := mux.Vars(r)
	dto.ID = vars["id"]

	event, err := dto.GetModel()
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	err = s.app.UpdateEvent(r.Context(), *event)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err)
		return
	}

	responseData, err := json.Marshal(dto)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseData)
}

func (s *ServerHandlers) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		RespondError(w, http.StatusBadRequest, err)
	}

	if err = s.app.DeleteEvent(r.Context(), id); err != nil {
		RespondError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (s *ServerHandlers) ListEvents(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	interval := r.URL.Query().Get("interval")

	withDate := false
	dtStart, err := time.Parse("2006-01-02", date)
	if err == nil {
		withDate = true
	}

	var events []storage.Event
	if withDate {
		switch interval {
		case "day":
			events, err = s.app.GetEventsStartedIn(r.Context(), dtStart, dtStart.AddDate(0, 0, 1).Sub(dtStart))
		case "week":
			events, err = s.app.GetEventsStartedIn(r.Context(), dtStart, dtStart.AddDate(0, 0, 7).Sub(dtStart))
		case "month":
			events, err = s.app.GetEventsStartedIn(r.Context(), dtStart, dtStart.AddDate(0, 1, 0).Sub(dtStart))
		default:
			events, err = s.app.GetEventsStartedIn(r.Context(), dtStart, dtStart.AddDate(0, 0, 1).Sub(dtStart))
		}
	} else {
		events, err = s.app.GetEvents(r.Context())
	}

	if err != nil {
		RespondError(w, http.StatusInternalServerError, err)
	}

	eventDtos := make([]EventDto, 0, len(events))
	for _, e := range events {
		eventDtos = append(eventDtos, CreateEventDtoFromModel(e))
	}

	responseData, err := json.Marshal(eventDtos)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseData)
}

func ParseRequest(r *http.Request, dto interface{}) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	err = json.Unmarshal(data, dto)
	if err != nil {
		return fmt.Errorf("failed to decode JSON request: %w", err)
	}

	return nil
}

func RespondError(w http.ResponseWriter, code int, err error) {
	data, err := json.Marshal(ErrorDto{
		false,
		err.Error(),
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to marshall error dto"))
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
