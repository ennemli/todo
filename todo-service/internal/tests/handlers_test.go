package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	handlers "github.com/ennemli/todo/todo/internal/handlers/todo"
	"github.com/ennemli/todo/todo/internal/models/todo"
	"github.com/ennemli/todo/todo/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func MakeRequest(req *http.Request) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	r := server.GetRouter()
	r.ServeHTTP(res, req)
	return res
}

func InitServe() {
	server.NewServer()
}

func TestCreateTodo(T *testing.T) {
	InitServe()
	mt := new(todo.MockTodo)
	todoHandlers := handlers.NewTodoHandlers(mt)
	expectedTodo := &todo.Todo{
		Date:        time.Date(2024, 5, 24, 9, 57, 38, 0, time.UTC),
		Name:        "Task 5",
		Description: "Do Task 1",
	}

	mt.On("CreateTodo", mock.Anything, expectedTodo).Return(expectedTodo, nil)

	reqBody, _ := json.Marshal(expectedTodo)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	r := server.GetRouter()
	r.Post("/", todoHandlers.CreateTodo)

	res := MakeRequest(req)

	assert.Equal(T, http.StatusOK, res.Code)

	var createdTodo todo.Todo
	err := json.NewDecoder(res.Body).Decode(&createdTodo)
	assert.Nil(T, err)

	assert.Equal(T, expectedTodo.Name, createdTodo.Name)
	assert.Equal(T, expectedTodo.Description, createdTodo.Description)

	mt.AssertExpectations(T)
}

func TestGetTodos(T *testing.T) {
	InitServe()
	r := server.GetRouter()
	mt := new(todo.MockTodo)
	todoHandlers := handlers.NewTodoHandlers(mt)
	expectedTodos := []*todo.Todo{
		{Name: "Task 1", Description: "Description 1"},
		{Name: "Task 2", Description: "Description 2"},
	}
	expectedTodos[0].ID = 1
	expectedTodos[1].ID = 2
	mt.On("GetTodos", mock.Anything).Return(expectedTodos, nil)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")

	r.Get("/", todoHandlers.GetTodos)

	res := MakeRequest(req)

	assert.Equal(T, http.StatusOK, res.Code)

	var todos []*todo.Todo
	err := json.NewDecoder(res.Body).Decode(&todos)
	assert.Nil(T, err)

	assert.Equal(T, len(expectedTodos), len(todos))

	mt.AssertExpectations(T)
}

func TestGetTodoById(T *testing.T) {
	InitServe()
	mt := new(todo.MockTodo)
	todoHandlers := handlers.NewTodoHandlers(mt)
	expectedTodo := &todo.Todo{
		Name:        "Task 1",
		Description: "Description 1",
	}
	todoID := uint(1)
	expectedTodo.ID = todoID
	mt.On("GetTodoById", mock.Anything, todoID).Return(expectedTodo, nil)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/%d", todoID), nil)
	req.Header.Set("Content-Type", "application/json")

	r := server.GetRouter()
	r.Get("/{id}", todoHandlers.GetTodoById)

	res := MakeRequest(req)

	assert.Equal(T, http.StatusOK, res.Code)

	var todoResponse todo.Todo
	err := json.NewDecoder(res.Body).Decode(&todoResponse)
	assert.Nil(T, err)

	assert.Equal(T, expectedTodo.ID, todoResponse.ID)
	assert.Equal(T, expectedTodo.Name, todoResponse.Name)
	assert.Equal(T, expectedTodo.Description, todoResponse.Description)

	mt.AssertExpectations(T)
}

func TestDeleteTodoById(T *testing.T) {
	InitServe()
	mt := new(todo.MockTodo)
	todoHandlers := handlers.NewTodoHandlers(mt)
	expectedTodo := &todo.Todo{
		Name:        "Task 1",
		Description: "Description 1",
	}
	todoID := uint(1)

	expectedTodo.ID = todoID
	mt.On("DeleteTodoById", mock.Anything, todoID).Return(expectedTodo, nil)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/%d", todoID), nil)
	req.Header.Set("Content-Type", "application/json")

	r := server.GetRouter()
	r.Delete("/{id}", todoHandlers.DeleteTodoById)

	res := MakeRequest(req)

	assert.Equal(T, http.StatusOK, res.Code)

	var todoResponse todo.Todo
	err := json.NewDecoder(res.Body).Decode(&todoResponse)
	assert.Nil(T, err)

	assert.Equal(T, expectedTodo.ID, todoResponse.ID)
	assert.Equal(T, expectedTodo.Name, todoResponse.Name)
	assert.Equal(T, expectedTodo.Description, todoResponse.Description)

	mt.AssertExpectations(T)
}

func TestUpdateTodo(T *testing.T) {
	InitServe()
	mt := new(todo.MockTodo)
	todoHandlers := handlers.NewTodoHandlers(mt)
	todoID := uint(1)

	existingTodo := &todo.Todo{
		Name:        "Original Todo",
		Description: "Original Description",
	}
	existingTodo.ID = todoID
	updatedFields := map[string]interface{}{
		"name": "Updated Todo",
	}
	expectedTodo := *existingTodo
	expectedTodo.Name = "Updated Todo"
	mt.On("GetTodoById", mock.Anything, todoID).Return(existingTodo, nil)

	mt.On("UpdateTodo", mock.Anything, existingTodo, updatedFields).Return(&expectedTodo, nil).Run(func(args mock.Arguments) {
		t := args.Get(1).(*todo.Todo)
		fields := args.Get(2).(map[string]interface{})
		t.Name = fields["name"].(string)
	})

	reqBody, _ := json.Marshal(updatedFields)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/%d", todoID), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	r := server.GetRouter()
	r.Put("/{id}", todoHandlers.UpdateTodo)

	res := MakeRequest(req)

	assert.Equal(T, http.StatusOK, res.Code)

	var updatedTodo todo.Todo
	err := json.NewDecoder(res.Body).Decode(&updatedTodo)
	assert.Nil(T, err)

	assert.Equal(T, existingTodo.ID, updatedTodo.ID)
	assert.Equal(T, updatedFields["name"], updatedTodo.Name)

	mt.AssertExpectations(T)
}
