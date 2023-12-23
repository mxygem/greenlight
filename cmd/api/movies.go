package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mxygem/greenlight/internal/data"
	"github.com/mxygem/greenlight/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("application createMovieHandler start")
	defer app.logger.Info("application createMovieHandler end")

	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Movies.Insert(movie); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("creating movie: %w", err))
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("writing create movie response: %w", err))
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("application showMovieHandler start")
	defer app.logger.Info("application showMovieHandler end")

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, fmt.Errorf("unknown error getting movie: %w", err))
		}

		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("writing show movie response: %w", err))
	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("application updateMovieHandler start")
	defer app.logger.Info("application updateMovieHandler end")

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, fmt.Errorf("unknown error getting movie to update: %w", err))
		}
		return
	}

	var input struct {
		Title   string
		Year    int32
		Runtime data.Runtime
		Genres  []string
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Movies.Update(movie); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("updating movie: %w", err))
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("writing update movie response: %w", err))
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("application deleteMovieHandler start")
	defer app.logger.Info("application deleteMovieHandler end")

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	if err := app.models.Movies.Delete(id); err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, fmt.Errorf("unknown error deleting movie: %w", err))
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("writing delete movie response: %w", err))
	}
}
