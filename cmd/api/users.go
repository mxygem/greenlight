package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mxygem/greenlight/internal/data"
	"github.com/mxygem/greenlight/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string
		Email    string
		Password string
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("reading new user input: %w", err))
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	if err := user.Password.Set(input.Password); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("setting user password: %w", err))
		return
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Users.Insert(user); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	app.background(func() {
		if err := app.mailer.Send(user.Email, "user_welcome.tmpl", user); err != nil {
			// app.serverErrorResponse(w, r, fmt.Errorf("sending welcome email: %w", err))
			app.logger.Error("sending email: %s", err)
		}
	})

	if err := app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}