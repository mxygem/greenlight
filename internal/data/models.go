package data

import (
	"database/sql"
	"errors"
	"testing"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

// Create a Models struct which wraps the MovieModel. We'll add other models to this,
// like a UserModel and PermissionModel, as our build progresses.
type Models struct {
	Movies Movies
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
	}
}

func NewMockModels(t *testing.T) Models {
	return Models{
		Movies: NewMockMovies(t),
	}
}
