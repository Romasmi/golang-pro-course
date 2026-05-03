package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/domain"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	goose "github.com/pressly/goose/v3"
)

type Storage struct {
	db *sqlx.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, dsn string) error {
	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}
	s.db = db
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Storage) Migrate(ctx context.Context, migrationsDir string) error {
	if migrationsDir == "" {
		return nil
	}
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", migrationsDir)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.UpContext(ctx, s.db.DB, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (s *Storage) AddEvent(ctx context.Context, event domain.Event) error {
	query := `INSERT INTO events (id, title, description, start_at, end_at, user_id) 
              VALUES (:id, :title, :description, :start_at, :end_at, :user_id)`
	_, err := s.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return fmt.Errorf("failed to add event: %w", err)
	}
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event domain.Event) error {
	query := `UPDATE events SET title=:title, description=:description, start_at=:start_at, 
              end_at=:end_at, user_id=:user_id, updated_at=NOW() WHERE id=:id`
	res, err := s.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrEventNotFound
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	query := `DELETE FROM events WHERE id=$1`
	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrEventNotFound
	}
	return nil
}

func (s *Storage) GetEventByID(ctx context.Context, id string) (domain.Event, error) {
	var event domain.Event
	query := `SELECT id, title, description, start_at, end_at, user_id, created_at, updated_at FROM events WHERE id=$1`
	err := s.db.GetContext(ctx, &event, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Event{}, domain.ErrEventNotFound
		}
		return domain.Event{}, fmt.Errorf("failed to get event: %w", err)
	}
	return event, nil
}

func (s *Storage) ListEvents(ctx context.Context) ([]domain.Event, error) {
	var events []domain.Event
	query := `SELECT id, title, description, start_at, end_at, user_id, created_at, updated_at FROM events`
	err := s.db.SelectContext(ctx, &events, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}
	return events, nil
}
