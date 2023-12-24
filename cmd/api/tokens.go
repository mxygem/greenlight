package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mxygem/greenlight/internal/data"
	"github.com/mxygem/greenlight/internal/validator"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("reading create auth token request body: %w", err))
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			fmt.Println("record not found")
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, fmt.Errorf("unexpected error getting user by email: %w", err))
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("matching password: %w", err))
		return
	}

	if !match {
		fmt.Println("no match")
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("creating new token: %w", err))
		return
	}

	if err := app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("writing token creation response: %w", err))
	}
}
