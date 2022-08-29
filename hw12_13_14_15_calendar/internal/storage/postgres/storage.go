package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	pgx4 "github.com/jackc/pgx/v4"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/domain/entity"
)

type Storage struct {
	ctx  context.Context
	conn *pgx4.Conn
	dsn  string
}

func New(ctx context.Context, dsn string) *Storage {
	return &Storage{
		ctx: ctx,
		dsn: dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) app.Storage {
	conn, err := pgx4.Connect(ctx, s.dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connect to database: %v\n", err)
		os.Exit(1)
	}

	s.conn = conn

	return s
}

func (s *Storage) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}

func (s *Storage) Insert(event entity.Event) error {
	fmt.Print("Storage Insert")
	sql := `INSERT INTO events (id, user_id, title, started_at, finished_at, description, notify_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.conn.Exec(
		s.ctx,
		sql,
		event.ID.String(),
		event.UserID,
		event.Title,
		event.StartedAt.Format(time.RFC3339),
		event.FinishedAt.Format(time.RFC3339),
		event.Description,
		event.NotifyAt.Format(time.RFC3339),
	)

	return err
}

func (s *Storage) Update(event entity.Event) error {
	sql := `UPDATE events 
			SET
				user_id = $1,
			    title = $2,
    			started_at = $3,
    			finished_at = $4,
    			description = $5,
    			notify_at = $6
			WHERE id = $7`

	_, err := s.conn.Exec(
		s.ctx,
		sql,
		event.UserID,
		event.Title,
		event.StartedAt.Format(time.RFC3339),
		event.FinishedAt.Format(time.RFC3339),
		event.Description,
		event.NotifyAt.Format(time.RFC3339),
		event.ID.String(),
	)

	return err
}

func (s *Storage) Delete(id uuid.UUID) error {
	sql := "DELETE FROM events WHERE id = $1"
	_, err := s.conn.Exec(s.ctx, sql, id)

	return err
}

func (s *Storage) Select() ([]entity.Event, error) {
	events := make([]entity.Event, 0)

	sql := `SELECT id, user_id, title, started_at, finished_at, description, notify_at 
			FROM events
			ORDER BY id`

	rows, err := s.conn.Query(s.ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event entity.Event
		if err := rows.Scan(
			&event.ID,
			&event.UserID,
			&event.Title,
			&event.StartedAt,
			&event.FinishedAt,
			&event.Description,
			&event.NotifyAt,
		); err != nil {
			return nil, fmt.Errorf("error scan result: %w", err)
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) SelectOne(id uuid.UUID) (*entity.Event, error) {
	var e entity.Event

	sql := `SELECT id, user_id, title, started_at, finished_at, description, notify_at 
			FROM events
			WHERE id = $1`
	err := s.conn.QueryRow(s.ctx, sql, id).Scan(
		&e.ID,
		&e.UserID,
		&e.Title,
		&e.StartedAt,
		&e.FinishedAt,
		&e.Description,
		&e.NotifyAt,
	)
	if err == nil {
		return &e, nil
	}

	if errors.Is(err, pgx4.ErrNoRows) {
		return nil, entity.ErrNotExistEvent
	}

	return nil, fmt.Errorf("error scan result: %w", err)
}

func (s *Storage) GetActualNotifyEvents(notifyTime time.Time) ([]entity.Event, error) {
	sql := `SELECT id, user_id, title, started_at, finished_at, description, notify_at
			FROM events 
			WHERE notify_at = $1`
	rows, err := s.conn.Query(s.ctx, sql, notifyTime.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsToEvents(rows)
}

func (s *Storage) GetOldEvents(timeBefore time.Time) ([]entity.Event, error) {
	sql := `SELECT id, user_id, title, started_at, finished_at, description, notify_at 
			FROM events 
			WHERE started_at <= $1`
	rows, err := s.conn.Query(s.ctx, sql, timeBefore.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsToEvents(rows)
}

func rowsToEvents(rows pgx.Rows) ([]entity.Event, error) {
	var events []entity.Event

	for rows.Next() {
		var e entity.Event
		if err := rows.Scan(
			&e.ID,
			&e.UserID,
			&e.Title,
			&e.StartedAt,
			&e.FinishedAt,
			&e.Description,
			&e.NotifyAt,
		); err != nil {
			return nil, fmt.Errorf("error scan result: %w", err)
		}

		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
