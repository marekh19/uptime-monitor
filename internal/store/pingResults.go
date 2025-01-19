package store

import (
	"context"
	"database/sql"
)

type PingResult struct {
	ID           string `json:"id"`
	MonitorID    string `json:"monitor_id"`
	Status       string `jhson:"status"`
	Timestamp    string `json:"timestamp"`
	ResponseTime int    `json:"response_time"`
}

type PingResultStore struct {
	db *sql.DB
}

func (s *PingResultStore) Create(ctx context.Context, pingResult *PingResult) error {
	query := `
    INSERT INTO ping_results (monitor_id, status, response_time, timestamp)
    VALUES ($1, $2, $3, $4)
    RETURNING id, timestamp
  `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		pingResult.MonitorID,
		pingResult.Status,
		pingResult.ResponseTime,
		pingResult.Timestamp,
	).Scan(&pingResult.ID, &pingResult.Timestamp)
	if err != nil {
		return err
	}

	return nil
}
