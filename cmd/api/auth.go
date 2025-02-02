package main

import (
	"net/http"

	"github.com/marekh19/uptime-ume/internal/store"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,min=3,max=40"`
	Password string `json:"password" validate:"required,min=8,max=40"`
}

// RegisterUser godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		main.RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	store.User					"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/auth/register [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
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

	user := &store.User{
		ID:       id,
		Username: payload.Username,
	}

	// Hash the password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Users.Create(ctx, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
