package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Monitors interface {
		Create(context.Context, *Monitor) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
	PingResults interface {
		Create(context.Context, *PingResult) error
	}
	StatusPages interface {
		Create(context.Context, *StatusPage) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Monitors:    &MonitorStore{db},
		Users:       &UsersStore{db},
		PingResults: &PingResultStore{db},
		StatusPages: &StatusPagesStore{db},
	}
}
