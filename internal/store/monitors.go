package store

import (
	"context"
	"database/sql"
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
