package main

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mxygem/greenlight/internal/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testUser struct {
	data.User
	testPass string
}

func TestRegisterUserHandler(t *testing.T) {
	testCases := []struct {
		name         string
		body         string
		mockUsers    func(t *testing.T) data.Users
		expected     testResp
		expectedCode int
	}{
		{
			name: "valid user",
			body: `{"name":"foo bar", "email":"foo.bar@email.com", "password":"supersecret"}`,
			mockUsers: func(t *testing.T) data.Users {
				pw := "supersecret"
				expected := &testUser{
					User: data.User{
						Name:  "foo bar",
						Email: "foo.bar@email.com",
					},
					testPass: pw,
				}

				mu := data.NewMockUsers(t)
				mu.On("Insert", mock.MatchedBy(expectedUserReceived(t, expected))).
					Run(mockInsertBehavior(t)).
					Return(nil)

				return mu
			},
			expected: testResp{
				User: data.User{
					ID:        100,
					CreatedAt: testCreatedAt,
					Name:      "foo bar",
					Email:     "foo.bar@email.com",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := &application{logger: newTestLogger(t)}
			if tc.mockUsers != nil {
				mu := tc.mockUsers(t)
				app.models = data.Models{Users: mu}
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(tc.body))

			app.registerUserHandler(w, r)

			var actual testResp
			readResp(t, w, &actual)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func newTestLogger(t *testing.T) *slog.Logger {
	return slog.New(slog.NewTextHandler(bytes.NewBuffer([]byte{}), nil))
}

// expectedUserReceived is a test helper that wraps and returns a custom matching function that
// satisfies the signature needed for testify's mock.MatchedBy function. It returns whether or not
// the user object received by the user's model's Insert method matches what we expect. This custom
// matcher ignores checking the user's hashed password since it changes each run.
func expectedUserReceived(t *testing.T, expected *testUser) func(*data.User) bool {
	t.Helper()

	return func(actual *data.User) bool {
		if expected.Name != actual.Name {
			return false
		}

		if expected.Email != actual.Email {
			return false
		}

		ok, err := actual.Password.Matches(expected.testPass)
		require.NoError(t, err)

		return ok
	}
}

func mockInsertBehavior(t *testing.T) func(mock.Arguments) {
	return func(args mock.Arguments) {
		user := args.Get(0).(*data.User)
		user.ID = 100
		user.CreatedAt = testCreatedAt
		user.Version = 1
		args[0] = user
	}
}
