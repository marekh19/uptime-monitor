package main

import (
	"net/http"
)

type HealthCheckPayload struct {
	Status  string `json:"status"`
	Env     string `json:"env"`
	Version string `json:"version"`
}

// HealthCheck godoc
//
//	@Summary		Check the health status
//	@Description	Check the health status of the API
//	@Tags			ops
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	main.HealthCheckPayload
//	@Failure		500	{object}	error
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := HealthCheckPayload{
		Status:  "ok",
		Env:     app.config.env,
		Version: version,
	}

	if err := writeJSON(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
