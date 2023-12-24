package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mxygem/greenlight/internal/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testCreatedAt = time.Date(2000, time.May, 23, 7, 2, 0, 0, time.Local)
)

type testResp struct {
	Error any        `json:"error,omitempty"`
	Movie data.Movie `json:"movie,omitempty"`
	User  data.User  `json:"user,omitempty"`
}

func TestCreateMovieHandler(t *testing.T) {
	testCases := []struct {
		desc       string
		body       string
		mockModels func(*testing.T) data.Models
		expected   testResp
	}{
		{
			desc: "empty body provided",
			body: `{}`,
			expected: testResp{
				Error: map[string]interface{}{
					"title":   "must be provided",
					"genres":  "must be provided",
					"year":    "must be provided",
					"runtime": "must be provided",
				},
			},
		},
		{
			desc: "provided year is string",
			body: `{"title":"test-title", "year":"999", "runtime":"102 mins", "genres":["test-genre-1", "test-genre-2"]}`,
			expected: testResp{
				Error: `body contains incorrect JSON type for field "year"`,
			},
		},
		{
			desc: "successful create",
			body: `{"title":"test-title", "year":1999, "runtime":"102 mins", "genres":["test-genre-1", "test-genre-2"]}`,
			mockModels: func(t *testing.T) data.Models {
				expected := &data.Movie{
					Title:   "test-title",
					Year:    1999,
					Runtime: 102,
					Genres:  []string{"test-genre-1", "test-genre-2"},
				}

				mov := data.NewMockMovies(t)
				mov.On("Insert", expected).
					Run(func(args mock.Arguments) {
						mv := args.Get(0).(*data.Movie)
						mv.ID = 100
						mv.CreatedAt = testCreatedAt
						mv.Version = 1
						args[0] = mv
					}).
					Return(nil)

				return data.Models{Movies: mov}
			},
			expected: testResp{
				Movie: data.Movie{
					ID:      100,
					Title:   "test-title",
					Year:    1999,
					Runtime: 102,
					Genres:  []string{"test-genre-1", "test-genre-2"},
					Version: 1,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(bytes.NewBuffer([]byte{}), nil))
			app := &application{logger: logger}
			if tc.mockModels != nil {
				app.models = tc.mockModels(t)
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/v1/movies", strings.NewReader(tc.body))

			app.createMovieHandler(w, r)

			var actual testResp
			readResp(t, w, &actual)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func readResp(t *testing.T, w *httptest.ResponseRecorder, actual any) {
	t.Helper()

	b, err := io.ReadAll(io.NopCloser(w.Result().Body))
	require.NoError(t, err)

	require.NoError(t, json.Unmarshal(b, &actual))
}
