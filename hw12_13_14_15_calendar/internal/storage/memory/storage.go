package memory

import (
	"sync"

	"github.com/google/uuid"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/domain/entity"
)

type Storage struct {
	mu     sync.RWMutex
	events map[uuid.UUID]entity.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[uuid.UUID]entity.Event),
	}
}

func (s *Storage) Insert(event entity.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; ok {
		return entity.ErrDuplicateEvent
	}

	s.events[event.ID] = event
	return nil
}

func (s *Storage) Update(event entity.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[event.ID] = event
	return nil
}

func (s *Storage) Delete(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return entity.ErrNotExistEvent
	}

	delete(s.events, id)
	return nil
}

func (s *Storage) Select() ([]entity.Event, error) {
	events := make([]entity.Event, 0, len(s.events))
	for _, event := range s.events {
		events = append(events, event)
	}
	return events, nil
}

func (s *Storage) SelectOne(id uuid.UUID) (*entity.Event, error) {
	event, has := s.events[id]
	if !has {
		return nil, entity.ErrNotExistEvent
	}

	return &event, nil
}
