package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/slon/shad-go/coverme/mocks"
	"gitlab.com/slon/shad-go/coverme/models"
)

var (
	errStoreGetAll  = errors.New("some error in getting all records")
	errStoreCreate  = errors.New("some error in creating record")
	errStoreGetTodo = errors.New("some error mb not found")
)

func TestList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStorage(ctrl)
	app := New(mockStore)
	app.initRoutes()

	gomock.InOrder(
		mockStore.EXPECT().
			GetAll().
			Return(
				make([]*models.Todo, 0),
				nil,
			),

		mockStore.EXPECT().
			GetAll().
			Return(
				nil,
				errStoreGetAll,
			),
	)

	// ok result
	req := httptest.NewRequest("GET", "http://example.com/todo", nil)
	w := httptest.NewRecorder()
	app.router.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	todos := make([]models.Todo, 0)
	assert.NoError(t, json.Unmarshal(body, &todos), "List returns no parsing data")
	assert.Equal(t, 0, len(todos), "Mock store returns [] but actual len != 0")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected OK status")

	// error result
	req = httptest.NewRequest("GET", "http://example.com/todo", nil)
	w = httptest.NewRecorder()
	app.router.ServeHTTP(w, req)

	resp = w.Result()
	assert.Equal(
		t,
		http.StatusInternalServerError, resp.StatusCode,
		"Expected internal server error",
	)
}

func TestAddTodo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStorage(ctrl)

	gomock.InOrder(
		mockStore.EXPECT().
			AddTodo("title", "content").
			Return(
				&models.Todo{},
				nil,
			),

		mockStore.EXPECT().
			AddTodo("title", "content").
			Return(
				nil,
				errStoreCreate,
			),
	)

	app := New(mockStore)
	app.initRoutes()

	for _, tc := range []struct {
		endpoint     string
		description  string
		payload      interface{}
		expectedCode int
	}{
		{
			endpoint:     "add-todo-handler-->",
			description:  "[> empty-body <]",
			payload:      "empty",
			expectedCode: http.StatusBadRequest,
		},
		{
			endpoint:    "add-todo-handler-->",
			description: "[> title-empty <]",
			payload: map[string]string{
				"a": "b",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			endpoint:    "add-todo-handler-->",
			description: "[> go-mock-returns-ok-in-create <]",
			payload: map[string]string{
				"title":   "title",
				"content": "content",
			},
			expectedCode: http.StatusCreated,
		},
		{
			endpoint:    "add-todo-handler-->",
			description: "[> go-mock-returns-err-in-create <]",
			payload: map[string]string{
				"title":   "title",
				"content": "content",
			},
			expectedCode: http.StatusInternalServerError,
		},
	} {
		t.Run(tc.endpoint+"-"+tc.description, func(t *testing.T) {
			body := &bytes.Buffer{}
			json.NewEncoder(body).Encode(tc.payload)
			req := httptest.NewRequest(
				"POST",
				"http://example.com/todo/create",
				body,
			)
			w := httptest.NewRecorder()
			app.router.ServeHTTP(w, req)

			resp := w.Result()
			require.Equal(t, tc.expectedCode, resp.StatusCode, "Status codes aren't equal")
		})
	}
}

func TestGetTodo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStorage(ctrl)

	todoID := 123
	gomock.InOrder(
		mockStore.EXPECT().
			GetTodo(models.ID(todoID)).
			Return(
				nil,
				errStoreGetTodo,
			),

		mockStore.EXPECT().
			GetTodo(models.ID(todoID)).
			Return(
				&models.Todo{},
				nil,
			),
	)

	app := New(mockStore)
	app.initRoutes()

	for _, tc := range []struct {
		endpoint     string
		description  string
		request      *http.Request
		expectedCode int
	}{
		{
			endpoint:    "get-todo-handler-->",
			description: "[> bad-id <]",
			request: httptest.NewRequest(
				"GET",
				"http://example.com/todo/999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999",
				nil,
			),
			expectedCode: http.StatusBadRequest,
		},
		{
			endpoint:    "get-todo-handler-->",
			description: "[> store-error <]",
			request: httptest.NewRequest(
				"GET",
				"http://example.com/todo/"+strconv.Itoa(todoID),
				nil,
			),
			expectedCode: http.StatusInternalServerError,
		},
		{
			endpoint:    "get-todo-handler-->",
			description: "[> ok <]",
			request: httptest.NewRequest(
				"GET",
				"http://example.com/todo/"+strconv.Itoa(todoID),
				nil,
			),
			expectedCode: http.StatusOK,
		},
	} {
		t.Run(tc.endpoint+"-"+tc.description, func(t *testing.T) {
			w := httptest.NewRecorder()
			app.router.ServeHTTP(w, tc.request)

			resp := w.Result()
			require.Equal(t, tc.expectedCode, resp.StatusCode, "Status codes aren't equal")
		})
	}
}

func TestStatus(t *testing.T) {
	w := httptest.NewRecorder()
	app := New(nil)
	app.status(w, nil)
	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode, "Status codes aren't equal")
}
