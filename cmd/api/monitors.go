package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/marekh19/uptime-monitor/internal/store"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type CreateMonitorPayload struct {
	Name     string `json:"name" validate:"required,max=100"`
	Address  string `json:"address" validate:"required,url"`
	Method   string `json:"method" validate:"omitempty,oneof=GET POST PUT PATCH DELETE HEAD OPTIONS"`
	Kind     string `json:"kind"`
	Config   string `json:"config"`
	Interval int    `json:"interval" validate:"required,gt=0"`
}

func (app *application) createMonitorHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateMonitorPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	id, err := gonanoid.New()
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// @TODO: Change this once auth is implemented
	userId := "1"

	monitor := &store.Monitor{
		ID:       id,
		UserId:   userId,
		Name:     payload.Name,
		Address:  payload.Address,
		Interval: payload.Interval,
		Method:   payload.Method,
		Kind:     payload.Kind,
		Config:   payload.Config,
	}

	ctx := r.Context()

	if err := app.store.Monitors.Create(ctx, monitor); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, monitor); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getMonitorHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		app.badRequestError(w, r, errors.New("missing id parameter"))
		return
	}

	ctx := r.Context()

	monitor, err := app.store.Monitors.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := writeJSON(w, http.StatusOK, monitor); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) listMonitorsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	monitors, err := app.store.Monitors.List(ctx)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	if err := writeJSON(w, http.StatusOK, monitors); err != nil {
		app.internalServerError(w, r, err)
	}
}
