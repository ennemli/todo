package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/ennemli/todo/todo/internal/errors"
	handlers "github.com/ennemli/todo/todo/internal/handlers/todo"
	"github.com/ennemli/todo/todo/internal/middlewares"
	"github.com/ennemli/todo/todo/internal/models/todo"
	"github.com/ennemli/todo/todo/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTimeOut(T *testing.T) {
	InitServe()
	r := server.GetRouter()
	r.Use(middlewares.SetTimeOut(errors.TimeoutDuration))
	mt := new(todo.MockTodo)
	todoHandlers := handlers.NewTodoHandlers(mt)
	r.Get("/", todoHandlers.GetTodos)
	mt.On("GetTodos", mock.Anything).Return([]*todo.Todo{}, nil).Run(func(args mock.Arguments) {
		time.Sleep(errors.TimeoutDuration + time.Second)
	})
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")

	res := MakeRequest(req)

	assert.Equal(T, http.StatusRequestTimeout, res.Code)
	var errorResponse errors.ErrorResponse
	err := json.NewDecoder(res.Body).Decode(&errorResponse)
	assert.Nil(T, err)
	assert.Equal(T, errors.ResponseRequestTimeout.Message, errorResponse.Message)
	mt.AssertExpectations(T)
}

func TestUpdateExtraFields(T *testing.T) {
	InitServe()
	mt := new(todo.MockTodo)
	todoHandlers := handlers.NewTodoHandlers(mt)
	updatedFields := map[string]interface{}{
		"name":       "Updated Todo",
		"extraFiedl": 12,
	}
	mt.On("GetTodoById", mock.Anything, mock.Anything).Return(&todo.Todo{}, nil)

	mt.On("UpdateTodo", mock.Anything, mock.Anything, mock.Anything).Return(&todo.Todo{}, nil)

	reqBody, _ := json.Marshal(updatedFields)
	req, _ := http.NewRequest("PUT", "/1", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	r := server.GetRouter()
	r.Put("/{id}", todoHandlers.UpdateTodo)

	res := MakeRequest(req)

	assert.Equal(T, http.StatusBadRequest, res.Code)

	var errorResponse errors.ErrorResponse
	err := json.NewDecoder(res.Body).Decode(&errorResponse)
	assert.Nil(T, err)

	assert.Equal(T, "Invalid request payload", errorResponse.Message)
}

func TestUpdateUserId(T *testing.T) {
	InitServe()
	mt := new(todo.MockTodo)
	todoHandlers := handlers.NewTodoHandlers(mt)
	updatedFields := map[string]interface{}{
		"name":   "Updated Todo",
		"userid": 12,
	}
	mt.On("GetTodoById", mock.Anything, mock.Anything).Return(&todo.Todo{}, nil)

	mt.On("UpdateTodo", mock.Anything, mock.Anything, mock.Anything).Return(&todo.Todo{}, nil)

	reqBody, _ := json.Marshal(updatedFields)
	req, _ := http.NewRequest("PUT", "/1", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	r := server.GetRouter()
	r.Put("/{id}", todoHandlers.UpdateTodo)

	res := MakeRequest(req)

	assert.Equal(T, http.StatusBadRequest, res.Code)

	var errorResponse errors.ErrorResponse
	err := json.NewDecoder(res.Body).Decode(&errorResponse)
	assert.Nil(T, err)

	assert.Equal(T, "Invalid request payload", errorResponse.Message)
}
