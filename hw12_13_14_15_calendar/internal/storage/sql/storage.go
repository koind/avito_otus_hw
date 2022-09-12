package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v4"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	ctx  context.Context
	conn *pgx.Conn
	dsn  string
}

func New(ctx context.Context, dsn string) *Storage {
	return &Storage{
		ctx: ctx,
		dsn: dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) (app.Storage, error) {
	conn, err := pgx.Connect(ctx, s.dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect database: %v\n", err)
		return nil, err
	}
	s.conn = conn
	return s, nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}

func (s *Storage) Create(e storage.Event) error {
	sql := `insert into events (id, title, started_at, finished_at, description, user_id, notify) 
			values ($1, $2, $3, $4, $5, $6, $7)
			`
	_, err := s.conn.Exec(
		s.ctx,
		sql,
		e.ID.String(),
		e.Title,
		e.StartedAt.Format(time.RFC3339),
		e.FinishedAt.Format(time.RFC3339),
		e.Description,
		e.UserID,
		e.Notify.Format(time.RFC3339),
	)

	return err
}

func (s *Storage) Update(e storage.Event) error {
	sql := `update events 
			set
    		title = $2,
    		started_at = $3,
    		finished_at = $4,
    		description = $5,
    		user_id = $6,
    		notify = $7
			where
			id = $1
			`
	_, err := s.conn.Exec(
		s.ctx,
		sql,
		e.ID.String(),
		e.Title,
		e.StartedAt.Format(time.RFC3339),
		e.FinishedAt.Format(time.RFC3339),
		e.Description,
		e.UserID,
		e.Notify.Format(time.RFC3339),
	)

	return err
}

func (s *Storage) Delete(id uuid.UUID) error {
	sql := `DELETE FROM events WHERE id = $1`
	_, err := s.conn.Exec(s.ctx, sql, id)

	return err
}

func (s *Storage) Find(id uuid.UUID) (*storage.Event, error) {
	var e storage.Event

	sql := `select id, title, started_at, finished_at, description, user_id, notify
			from events
			where id = $1
			`
	err := s.conn.QueryRow(s.ctx, sql, id).Scan(
		&e.ID,
		&e.Title,
		&e.StartedAt,
		&e.FinishedAt,
		&e.Description,
		&e.UserID,
		&e.Notify,
	)

	if err == nil {
		return &e, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.ErrorEventNotFound
	}

	return nil, fmt.Errorf("failed to scan SQL result into struct: %w", err)
}

func (s *Storage) FindAll() ([]storage.Event, error) {
	var events []storage.Event

	sql := `select id, title, started_at, finished_at, description, user_id, notify 
			from events
			order by date
			`
	rows, err := s.conn.Query(s.ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e storage.Event
		if err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.StartedAt,
			&e.FinishedAt,
			&e.Description,
			&e.UserID,
			&e.Notify,
		); err != nil {
			return nil, fmt.Errorf("unable to transform array result into struct: %w", err)
		}

		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) FindAtDay(day time.Time) ([]storage.Event, error) { //nolint:dupl
	var events []storage.Event

	from := day.AddDate(0, 0, 1).Format(time.RFC3339)
	to := day.AddDate(0, 0, 1).Format(time.RFC3339)

	rows, err := s.findAtDate(from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e storage.Event
		if err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.StartedAt,
			&e.FinishedAt,
			&e.Description,
			&e.UserID,
			&e.Notify,
		); err != nil {
			return nil, fmt.Errorf("unable to transform array result into struct: %w", err)
		}

		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) FindAtWeek(dayStart time.Time) ([]storage.Event, error) { //nolint:dupl
	var events []storage.Event

	from := dayStart.AddDate(0, 0, 7).Format(time.RFC3339)
	to := dayStart.AddDate(0, 0, 7).Format(time.RFC3339)

	rows, err := s.findAtDate(from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e storage.Event
		if err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.StartedAt,
			&e.FinishedAt,
			&e.Description,
			&e.UserID,
			&e.Notify,
		); err != nil {
			return nil, fmt.Errorf("unable to transform array result into struct: %w", err)
		}

		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) FindAtMonth(dayStart time.Time) ([]storage.Event, error) { //nolint:dupl
	var events []storage.Event

	from := dayStart.AddDate(0, 1, 0).Format(time.RFC3339)
	to := dayStart.AddDate(0, 1, 0).Format(time.RFC3339)

	rows, err := s.findAtDate(from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e storage.Event
		if err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.StartedAt,
			&e.FinishedAt,
			&e.Description,
			&e.UserID,
			&e.Notify,
		); err != nil {
			return nil, fmt.Errorf("unable to transform array result into struct: %w", err)
		}

		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) findAtDate(from, to string) (pgx.Rows, error) {
	const searchSQL = `select id, title, started_at, finished_at, description, user_id, notify 
						from events
						where started_at >= $1 and started_at <= $2
						order by date
						`
	return s.conn.Query(s.ctx, searchSQL, from, to)
}

func (s *Storage) GetEventsReadyToNotify(dt time.Time) ([]storage.Event, error) {
	sql := `
 select id, title, started_at, finished_at, description, user_id, notify 
 from events 
 where notify <= $1
 `
	rows, err := s.conn.Query(s.ctx, sql, dt.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return createEventsFromRows(rows)
}

func (s *Storage) GetEventsOlderThan(dt time.Time) ([]storage.Event, error) {
	sql := `
 select id, title, started_at, finished_at, description, user_id, notify 
 from events 
 where started_at <= $1
 `
	rows, err := s.conn.Query(s.ctx, sql, dt.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return createEventsFromRows(rows)
}

func createEventsFromRows(rows pgx.Rows) ([]storage.Event, error) {
	var events []storage.Event

	for rows.Next() {
		var e storage.Event
		if err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.StartedAt,
			&e.FinishedAt,
			&e.Description,
			&e.UserID,
			&e.Notify,
		); err != nil {
			return nil, fmt.Errorf("failed to scan SQL result into struct: %w", err)
		}

		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
