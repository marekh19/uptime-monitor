package store

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string   `json:"id"`
	Username  string   `json:"username"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Password  password `json:"-"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, user *User) error {
	query := `
    INSERT INTO users (id, username, password)
    VALUES ($1, $2, $3)
    RETURNING id, username, created_at, updated_at
  `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.ID,
		user.Username,
		user.Password.hash,
	).Scan(&user.ID, &user.Username, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
