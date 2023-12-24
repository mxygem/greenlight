package data

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mxygem/greenlight/internal/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

var (
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
			plaintext: longStr(t, 10001),
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
		})
	}
}

func TestMatches(t *testing.T) {
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

func TestValidateUser(t *testing.T) {
	testCases := []struct {
		name         string
		user         *User
		expected     bool
		expectedErrs map[string]string
	}{
		{
			name: "valid user",
			user: &User{
				Name:  "Foo Bar",
				Email: "foo.bar@email.com",
				Password: password{
					plaintext: strPtr(t, "abcd1234"),
					hash:      []byte("thisIsAHashIPromise"),
				},
			},
			expected: true,
		},
		{
			name: "no name",
			user: &User{
				Name:     "",
				Email:    "foo.bar@email.com",
				Password: password{plaintext: strPtr(t, "abcd1234"), hash: []byte("t")},
			},
			expected: false,
			expectedErrs: map[string]string{
				"name": "must be provided",
			},
		},
		{
			name: "name too long",
			user: &User{
				Name:     longStr(t, 501),
				Email:    "foo.bar@email.com",
				Password: password{plaintext: strPtr(t, "abcd1234"), hash: []byte("t")},
			},
			expected: false,
			expectedErrs: map[string]string{
				"name": "must not be more than 500 bytes long",
			},
		},
		{
			name: "empty email",
			user: &User{
				Name:     "Foo Bar",
				Email:    "",
				Password: password{plaintext: strPtr(t, "abcd1234"), hash: []byte("t")},
			},
			expected: false,
			expectedErrs: map[string]string{
				"email": "must be provided",
			},
		},
		{
			name: "invalid email - only first name",
			user: &User{
				Name:     "Foo Bar",
				Email:    "foo",
				Password: password{plaintext: strPtr(t, "abcd1234"), hash: []byte("t")},
			},
			expected: false,
			expectedErrs: map[string]string{
				"email": "must be a valid email address",
			},
		},
		{
			name: "invalid email - only name",
			user: &User{
				Name:     "Foo Bar",
				Email:    "foo.bar",
				Password: password{plaintext: strPtr(t, "abcd1234"), hash: []byte("t")},
			},
			expected: false,
			expectedErrs: map[string]string{
				"email": "must be a valid email address",
			},
		},
		{
			name: "invalid email - missing domain",
			user: &User{
				Name:     "Foo Bar",
				Email:    "foo.bar@",
				Password: password{plaintext: strPtr(t, "abcd1234"), hash: []byte("t")},
			},
			expected: false,
			expectedErrs: map[string]string{
				"email": "must be a valid email address",
			},
		},
		{
			name: "invalid email - missing domain end",
			user: &User{
				Name:     "Foo Bar",
				Email:    "bar@email",
				Password: password{plaintext: strPtr(t, "abcd1234"), hash: []byte("t")},
			},
			expected: false,
			expectedErrs: map[string]string{
				"email": "must be a valid email address",
			},
		},
		{
			name: "invalid email - missing domain ending",
			user: &User{
				Name:     "Foo Bar",
				Email:    "foo.bar@email.",
				Password: password{plaintext: strPtr(t, "abcd1234"), hash: []byte("t")},
			},
			expected: false,
			expectedErrs: map[string]string{
				"email": "must be a valid email address",
			},
		},
		{
			name: "invalid email - only domain",
			user: &User{
				Name:     "Foo Bar",
				Email:    "email.com",
				Password: password{plaintext: strPtr(t, "abcd1234"), hash: []byte("t")},
			},
			expected: false,
			expectedErrs: map[string]string{
				"email": "must be a valid email address",
			},
		},
		{
			name: "invalid password - empty",
			user: &User{Name: "Foo Bar", Email: "foo.bar@email.com",
				Password: password{
					plaintext: strPtr(t, ""),
					hash:      []byte("t"),
				},
			},
			expected: false,
			expectedErrs: map[string]string{
				"password": "must be provided",
			},
		},
		{
			name: "invalid password - too short",
			user: &User{Name: "Foo Bar", Email: "foo.bar@email.com",
				Password: password{
					plaintext: strPtr(t, "hi"),
					hash:      []byte("t"),
				},
			},
			expected: false,
			expectedErrs: map[string]string{
				"password": "must be at least 8 bytes long",
			},
		},
		{
			name: "invalid password - too long",
			user: &User{Name: "Foo Bar", Email: "foo.bar@email.com",
				Password: password{
					plaintext: strPtr(t, longStr(t, 73)),
					hash:      []byte("t"),
				},
			},
			expected: false,
			expectedErrs: map[string]string{
				"password": "must not be more than 72 bytes long",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := validator.New()

			ValidateUser(v, tc.user)

			assert.Equal(t, tc.expected, v.Valid())
			if tc.expected {
				assert.Len(t, v.Errors, 0)
			} else {
				assert.Equal(t, tc.expectedErrs, v.Errors)
			}
		})
	}
}

func TestValidateUserPanicWithNoHash(t *testing.T) {
	v := validator.New()

	assert.Panics(t, func() { ValidateUser(v, &User{}) })
}

func TestUserInsert(t *testing.T) {
	insertQuery := `
		insert into users (name, email, password_hash, activated)
		values ($1, $2, $3, $4)
		returning id, created_at, version`
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

func longStr(t *testing.T, len int) string {
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
