package data

import (
	"database/sql"
	"errors"
	"testing"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Create a Models struct which wraps the MovieModel. We'll add other models to this,
// like a UserModel and PermissionModel, as our build progresses.
type Models struct {
	Movies      Movies
	Users       Users
	Tokens      Tokens
	Permissions Permissions
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies:      MovieModel{DB: db},
		Users:       UserModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Permissions: PermissionModel{DB: db},
	}
}

func NewMockModels(t *testing.T) Models {
	return Models{
		Movies:      NewMockMovies(t),
		Users:       NewMockUsers(t),
		Tokens:      NewMockTokens(t),
		Permissions: NewMockPermissions(t),
	}
}
