package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

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

	if err := app.models.Permissions.AddForUser(user.ID, "movies:read"); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("creating new token: %w", err))
		return
	}

	app.background(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
			"name":            user.Name,
		}

		if err := app.mailer.Send(user.Email, "user_welcome.tmpl", data); err != nil {
			app.logger.Error("sending email: %s", err)
		}
	})

	if err := app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("reading request body: %w", err))
		return
	}

	v := validator.New()

	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, fmt.Errorf("unexpected error getting user by token: %w", err))
		}

		return
	}

	if err := app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("deleting tokens for activated user: %w", err))
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("writing activation response: %w", err))
	}
}
