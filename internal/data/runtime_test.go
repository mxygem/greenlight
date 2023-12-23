package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalJSON(t *testing.T) {
	testCases := []struct {
		name        string
		runtime     Runtime
		expected    string
		expectedErr error
	}{
		{
			name:     "successful",
			runtime:  42,
			expected: `"42 mins"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.runtime.MarshalJSON()

			assert.Equal(t, []byte(tc.expected), actual)
			checkErr(t, tc.expectedErr, err)
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name        string
		jsonValue   []byte
		expected    Runtime
		expectedErr error
	}{
		{
			name:        "nil input",
			jsonValue:   nil,
			expectedErr: ErrInvalidRuntimeFormat,
		},
		{
			name:        "input unquoted",
			jsonValue:   []byte(`2113`),
			expectedErr: ErrInvalidRuntimeFormat,
		},
		{
			name:        "missing mins",
			jsonValue:   []byte(`"2113"`),
			expectedErr: ErrInvalidRuntimeFormat,
		},
		{
			name:        "time not numbers",
			jsonValue:   []byte(`"Twenty mins"`),
			expectedErr: ErrInvalidRuntimeFormat,
		},
		{
			name:      "successful",
			jsonValue: []byte(`"2113 mins"`),
			expected:  Runtime(2113),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runtime := Runtime(0)

			err := runtime.UnmarshalJSON(tc.jsonValue)

			assert.Equal(t, tc.expected, runtime)
			checkErr(t, tc.expectedErr, err)
		})
	}
}
