package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/marekh19/uptime-monitor/internal/store"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type CreateMonitorPayload struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Method   string `json:"method"`
	Kind     string `json:"kind"`
	Config   string `json:"config"`
	Interval int    `json:"interval"`
}

func (app *application) createMonitorHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateMonitorPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := gonanoid.New()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
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
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, monitor); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) getMonitorHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing id parameter")
		return
	}

	ctx := r.Context()

	monitor, err := app.store.Monitors.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeJSONError(w, http.StatusNotFound, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := writeJSON(w, http.StatusOK, monitor); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
