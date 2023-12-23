package data

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

var (
	usersTableRows       = sqlmock.NewRows([]string{"name", "email", "password_hash", "activated"})
	usersTableReturnRows = sqlmock.NewRows([]string{"id", "created_at", "version"})
	testCreatedAt        = time.Date(2000, time.May, 23, 7, 2, 0, 0, time.Local)
)

func TestSet(t *testing.T) {
	testCases := []struct {
		name      string
		plaintext string
		expected  error
	}{
		{
			name:      "empty password",
			plaintext: "",
			expected:  fmt.Errorf("password input required"),
		},
		{
			name:      "password too long",
			plaintext: pass(t, 10001),
			expected:  fmt.Errorf("generating hash from password: bcrypt: password length exceeds 72 bytes"),
		},
		{
			name:      "success",
			plaintext: "Th3F1ct10nSh4llM33tTh3R34l",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pw := password{}

			err := pw.Set(tc.plaintext)

			checkErr(t, tc.expected, err)
			if tc.expected == nil {
				require.NotNil(t, pw.plaintext)
				assert.Equal(t, tc.plaintext, *pw.plaintext)
				assert.NoError(t, bcrypt.CompareHashAndPassword(pw.hash, []byte(tc.plaintext)))
			}
			t.Fail()
		})
	}
}

func TestMatch(t *testing.T) {
	testCases := []struct {
		name        string
		plaintext   string
		password    password
		expected    bool
		expectedErr error
	}{
		{
			name:      "empty input",
			plaintext: "",
			password: password{
				plaintext: strPtr(t, "Th3F1ct10nSh4llM33tTh3R34l"),
				hash:      []byte("$2a$12$OojW3AfLxGDwWrhse73JlejlMMFabpJf9kci.9fLVAUdL9HAAA87W"),
			},
			expectedErr: fmt.Errorf("password input required"),
		},
		{
			name:      "successful match",
			plaintext: "Th3F1ct10nSh4llM33tTh3R34l",
			password: password{
				plaintext: strPtr(t, "Th3F1ct10nSh4llM33tTh3R34l"),
				hash:      []byte("$2a$12$OojW3AfLxGDwWrhse73JlejlMMFabpJf9kci.9fLVAUdL9HAAA87W"),
			},
			expected: true,
		},
		{
			name:      "no match",
			plaintext: "foobar",
			password: password{
				plaintext: strPtr(t, "Th3F1ct10nSh4llM33tTh3R34l"),
				hash:      []byte("$2a$12$OojW3AfLxGDwWrhse73JlejlMMFabpJf9kci.9fLVAUdL9HAAA87W"),
			},
			expected: false,
		},
		{
			name:        "password object empty",
			plaintext:   "foobar",
			password:    password{},
			expectedErr: fmt.Errorf("matching password: crypto/bcrypt: hashedSecret too short to be a bcrypted password"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.password.Matches(tc.plaintext)

			assert.Equal(t, tc.expected, actual)
			checkErr(t, tc.expectedErr, err)
		})
	}
}

func TestUserInsert(t *testing.T) {
	insertQuery := `
		insert into users (name, email, password_hash, activated)
		values ($1, $2, $3, $4)
		return id, created_at, version`
	testCases := []struct {
		name        string
		user        *User
		mockSetup   func(t *testing.T, db sqlmock.Sqlmock, user *User)
		expected    *User
		expectedErr error
	}{
		{
			name: "TBD - user data empty",
			user: &User{},
			mockSetup: func(t *testing.T, db sqlmock.Sqlmock, user *User) {
				t.SkipNow()

				var hash []uint8
				db.ExpectQuery(regexp.QuoteMeta(insertQuery)).
					WithArgs("", "", hash, false).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "version"}))
			},
			expectedErr: fmt.Errorf("nope"),
		},
		{
			name: "successful creation",
			user: &User{
				Name:      "u202",
				Email:     "u202@email.com",
				Password:  password{hash: []byte("$2a$12$OojW3AfLxGDwWrhse73JlejlMMFabpJf9kci.9fLVAUdL9HAAA87W")},
				Activated: false,
			},
			mockSetup: func(t *testing.T, db sqlmock.Sqlmock, user *User) {
				db.ExpectQuery(regexp.QuoteMeta(insertQuery)).
					WithArgs(user.Name, user.Email, user.Password.hash, user.Activated).
					WillReturnRows(usersTableReturnRows.AddRow(
						2113, testCreatedAt, 1,
					))
			},
			expected: &User{
				ID:        2113,
				CreatedAt: testCreatedAt,
				Version:   1,
				Name:      "u202",
				Email:     "u202@email.com",
				Password:  password{hash: []byte("$2a$12$OojW3AfLxGDwWrhse73JlejlMMFabpJf9kci.9fLVAUdL9HAAA87W")},
				Activated: false,
			},
		},
		{
			name: "email conflict",
			user: &User{
				Name:      "u409",
				Email:     "u409@email.com",
				Password:  password{hash: []byte("!1234567890!")},
				Activated: false,
			},
			mockSetup: func(t *testing.T, db sqlmock.Sqlmock, user *User) {
				db.ExpectQuery(regexp.QuoteMeta(insertQuery)).
					WithArgs(user.Name, user.Email, user.Password.hash, user.Activated).
					WillReturnError(fmt.Errorf(`pq: duplicate key value violates unique constraint "users_email_key"`))
			},
			expectedErr: fmt.Errorf("duplicate email"),
		},
		{
			name: "TODO - invalid email",
			user: &User{
				Name:  "u422",
				Email: "u422",
			},
			mockSetup: func(t *testing.T, db sqlmock.Sqlmock, user *User) {
				t.SkipNow()
				db.ExpectQuery(regexp.QuoteMeta(insertQuery)).
					WithArgs(user.Name, user.Email, user.Password.hash, user.Activated).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "version"}))
			},
			expectedErr: fmt.Errorf("must be a valid email address"),
		},
		{
			name: "query throws generic error",
			user: &User{},
			mockSetup: func(t *testing.T, db sqlmock.Sqlmock, user *User) {
				db.ExpectQuery(regexp.QuoteMeta(insertQuery)).
					WithArgs(user.Name, user.Email, user.Password.hash, user.Activated).
					WillReturnError(fmt.Errorf("it broke"))
			},
			expectedErr: fmt.Errorf("unexpected error inserting user: it broke"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			tc.mockSetup(t, mock, tc.user)
			um := UserModel{DB: db}

			err = um.Insert(tc.user)

			checkErr(t, tc.expectedErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func checkErr(t *testing.T, expected, actual error) {
	t.Helper()

	if expected != nil {
		require.Error(t, actual)
		assert.Equal(t, expected.Error(), actual.Error())
	} else {
		assert.NoError(t, actual)
	}
}

func pass(t *testing.T, len int) string {
	var out strings.Builder
	out.Grow(len)

	for i := 0; i < len; i++ {
		_, err := out.WriteString("a")
		require.NoError(t, err)
	}

	return out.String()
}

func strPtr(t *testing.T, str string) *string {
	return &str
}