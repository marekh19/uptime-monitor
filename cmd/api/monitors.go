package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/marekh19/uptime-monitor/internal/store"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type monitorKey string

const monitorCtx monitorKey = "monitor"

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

	if err := app.jsonResponse(w, http.StatusCreated, monitor); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getMonitorHandler(w http.ResponseWriter, r *http.Request) {
	monitor := getMonitorFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, monitor); err != nil {
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

	if err := app.jsonResponse(w, http.StatusOK, monitors); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deleteMonitorHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		app.badRequestError(w, r, errors.New("missing id parameter"))
		return
	}

	ctx := r.Context()

	if err := app.store.Monitors.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdateMonitorPayload struct {
	Name     *string `json:"name" validate:"omitempty,max=100"`
	Address  *string `json:"address" validate:"omitempty,url"`
	Method   *string `json:"method" validate:"omitempty,oneof=GET POST PUT PATCH DELETE HEAD OPTIONS"`
	Kind     *string `json:"kind"`
	Config   *string `json:"config"`
	Interval *int    `json:"interval" validate:"omitempty,gt=0"`
}

func (app *application) updateMonitorHandler(w http.ResponseWriter, r *http.Request) {
	monitor := getMonitorFromContext(r)

	var payload UpdateMonitorPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if payload.Name != nil {
		monitor.Name = *payload.Name
	}

	if payload.Address != nil {
		monitor.Address = *payload.Address
	}

	if payload.Method != nil {
		monitor.Method = *payload.Method
	}

	if payload.Kind != nil {
		monitor.Kind = *payload.Kind
	}

	if payload.Config != nil {
		monitor.Config = *payload.Config
	}

	if payload.Interval != nil {
		monitor.Interval = *payload.Interval
	}

	if err := app.store.Monitors.Update(r.Context(), monitor); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, monitor); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) monitorContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx = context.WithValue(ctx, monitorCtx, monitor)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getMonitorFromContext(r *http.Request) *store.Monitor {
	monitor, _ := r.Context().Value(monitorCtx).(*store.Monitor)
	return monitor
}
