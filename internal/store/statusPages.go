package store

import (
	"context"
	"database/sql"
)

type StatusPage struct {
	ID         string   `json:"id"`
	UserID     string   `json:"user_id"`
	Name       string   `json:"name"`
	Slug       string   `json:"slug"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
	MonitorIDs []string `json:"monitors"`
}

type StatusPagesStore struct {
	db *sql.DB
}

func (s *StatusPagesStore) Create(ctx context.Context, statusPage *StatusPage) error {
	query := `
    INSERT INTO status_pages (user_id, name, slug, monitor_ids)
    VALUES ($1, $2, $3, $4)
    RETURNING id, created_at, updated_at
  `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		statusPage.UserID,
		statusPage.Name,
		statusPage.Slug,
		statusPage.MonitorIDs).Scan(&statusPage.ID, &statusPage.CreatedAt, &statusPage.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
