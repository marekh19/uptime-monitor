package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Monitor struct {
	ID        string `json:"id"`
	UserId    string `json:"user_id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Method    string `json:"method"`
	Kind      string `json:"kind"`
	Config    string `json:"config"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Interval  int    `json:"interval"`
	Version   int    `json:"version"`
}

type MonitorStore struct {
	db *sql.DB
}

func (s *MonitorStore) Create(ctx context.Context, monitor *Monitor) error {
	query := `
    INSERT INTO monitors (id, user_id, name, address, interval, method, kind, config)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING id, created_at, updated_at;
  `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		monitor.ID,
		monitor.UserId,
		monitor.Name,
		monitor.Address,
		monitor.Interval,
		monitor.Method,
		monitor.Kind,
		monitor.Config,
	).Scan(&monitor.ID, &monitor.CreatedAt, &monitor.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *MonitorStore) GetByID(ctx context.Context, id string) (*Monitor, error) {
	query := `
    SELECT id, user_id, name, address, method, kind, config, created_at, updated_at, interval, version
    FROM monitors
    WHERE id = $1;
  `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var monitor Monitor

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&monitor.ID,
		&monitor.UserId,
		&monitor.Name,
		&monitor.Address,
		&monitor.Method,
		&monitor.Kind,
		&monitor.Config,
		&monitor.CreatedAt,
		&monitor.UpdatedAt,
		&monitor.Interval,
		&monitor.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &monitor, nil
}

func (s *MonitorStore) List(ctx context.Context) ([]*Monitor, error) {
	query := `
    SELECT id, user_id, name, address, method, kind, config, created_at, updated_at, interval, version
    FROM monitors;
  `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch monitors: %w", err)
	}
	defer rows.Close()

	var monitors []*Monitor
	for rows.Next() {
		var monitor Monitor
		err := rows.Scan(
			&monitor.ID,
			&monitor.UserId,
			&monitor.Name,
			&monitor.Address,
			&monitor.Method,
			&monitor.Kind,
			&monitor.Config,
			&monitor.CreatedAt,
			&monitor.UpdatedAt,
			&monitor.Interval,
			&monitor.Version,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan monitor: %w", err)
		}
		monitors = append(monitors, &monitor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return monitors, nil
}

func (s *MonitorStore) Delete(ctx context.Context, id string) error {
	query := `
    DELETE FROM monitors
    WHERE id = $1;
  `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *MonitorStore) Update(ctx context.Context, monitor *Monitor) error {
	query := `
    UPDATE monitors
    SET
      name = COALESCE($1, name),
      address = COALESCE($2, address),
      interval = COALESCE($3, interval),
      method = COALESCE($4, method),
      kind = COALESCE($5, kind),
      config = COALESCE($6, config),
      version = version + 1
    WHERE id = $7 AND version = $8
    RETURNING version;
  `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		monitor.Name,
		monitor.Address,
		monitor.Interval,
		monitor.Method,
		monitor.Kind,
		monitor.Config,
		monitor.ID,
		monitor.Version,
	).Scan(&monitor.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}
