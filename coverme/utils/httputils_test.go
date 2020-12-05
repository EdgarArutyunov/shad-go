package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestRespondJSON(t *testing.T) {

	for _, tc := range []struct {
		endpoint     string
		description  string
		data         interface{}
		expectedCode int
		expectedErr  bool
	}{
		{
			endpoint:    "RespondJSON",
			description: "ok data",
			data: map[string]string{
				"meet_id": "123",
			},
			expectedCode: http.StatusOK,
			expectedErr:  false,
		},
		{
			endpoint:     "RespondJSON",
			description:  "err-marshal",
			data:         make(chan struct{}),
			expectedCode: http.StatusOK,
			expectedErr:  true,
		},
	} {
		t.Run(tc.endpoint+"-"+tc.description, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := RespondJSON(w, tc.expectedCode, tc.data)

			if tc.expectedErr {
				assert.Error(t, err, "expected error in respond")
			} else {
				assert.NoError(t, err, "expected noerror in respond")
			}

			resp := w.Result()
			require.Equal(
				t,
				tc.expectedCode, resp.StatusCode,
				"Status codes aren't equal",
			)
		})
	}
}

func TestServerError(t *testing.T) {
	w := httptest.NewRecorder()
	ServerError(w)

	resp := w.Result()
	require.Equal(
		t,
		http.StatusInternalServerError, resp.StatusCode,
		"Status codes aren't equal",
	)
}

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	BadRequest(w, "hello")

	resp := w.Result()
	require.Equal(
		t,
		http.StatusBadRequest, resp.StatusCode,
		"Status codes aren't equal",
	)
}
