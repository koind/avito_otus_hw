package https

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/domain/entity"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/domain/presenter"
)

type EventHandler struct {
	app *app.App
}

func NewEventHandler(app *app.App) *EventHandler {
	return &EventHandler{app: app}
}

func (s *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	form := EventForm{}

	if err := parseRequest(r, &form); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	event, err := form.ToEntity()
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	err = s.app.CreateEvent(r.Context(), *event)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	response, err := json.Marshal(event.ToPresenter())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondSuccess(w, http.StatusCreated, response)
}

func (s *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	form := EventForm{}

	if err := parseRequest(r, &form); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	vars := mux.Vars(r)
	form.ID = vars["id"]

	event, err := form.ToEntity()
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	err = s.app.UpdateEvent(r.Context(), *event)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	response, err := json.Marshal(event.ToPresenter())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondSuccess(w, http.StatusCreated, response)
}

func (s *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
	}

	err = s.app.DeleteEvent(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondSuccess(w, http.StatusNoContent, nil)
}

func (s *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	interval := r.URL.Query().Get("interval")

	withDate := false
	dtStart, err := time.Parse("2006-01-02", date)
	if err == nil {
		withDate = true
	}

	var events []entity.Event
	if withDate {
		switch interval {
		case "day":
			events, err = s.app.GetDayEvents(r.Context(), dtStart)
		case "week":
			events, err = s.app.GetWeekEvents(r.Context(), dtStart)
		case "month":
			events, err = s.app.GetMonthEvents(r.Context(), dtStart)
		default:
			events, err = s.app.GetDayEvents(r.Context(), dtStart)
		}
	} else {
		events, err = s.app.GetEvents(r.Context())
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	list := make([]presenter.Event, 0, len(events))

	for _, event := range events {
		list = append(list, event.ToPresenter())
	}

	response, err := json.Marshal(list)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondSuccess(w, http.StatusCreated, response)
}

func parseRequest(r *http.Request, form *EventForm) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	err = json.Unmarshal(data, form)
	if err != nil {
		return fmt.Errorf("failed to decode JSON request: %w", err)
	}

	return nil
}

func respondError(w http.ResponseWriter, code int, err error) {
	data, err := json.Marshal(Error{
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

func respondSuccess(w http.ResponseWriter, code int, response []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
