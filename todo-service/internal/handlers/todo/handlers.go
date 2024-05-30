package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ennemli/todo/todo/internal/errors"
	"github.com/ennemli/todo/todo/internal/models/todo"
	"github.com/ennemli/todo/todo/pkg/maputil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Handlers interface {
	GetTodos(w http.ResponseWriter, r *http.Request)
	GetTodoById(w http.ResponseWriter, r *http.Request)
	CreateTodo(w http.ResponseWriter, r *http.Request)
	DeleteTodoById(w http.ResponseWriter, r *http.Request)
	UpdateTodo(w http.ResponseWriter, r *http.Request)
}

type todoHandlers struct {
	store todo.Store
}

func NewTodoHandlers(store todo.Store) Handlers {
	return &todoHandlers{store: store}
}

func (h *todoHandlers) GetTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.store.GetTodos(r.Context())
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, todos)
}

func (h *todoHandlers) GetTodoById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, "Invalid ID")
		return
	}
	todo, err := h.store.GetTodoById(r.Context(), uint(id))
	if err != nil {
		renderError(w, r, http.StatusNotFound, fmt.Sprintf("Todo with ID %d not found", id))
		return
	}
	render.JSON(w, r, todo)
}

func (h *todoHandlers) DeleteTodoById(w http.ResponseWriter, r *http.Request) {
	todo := &todo.Todo{}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, "Invalid ID")
		return
	}
	todo, err = h.store.DeleteTodoById(r.Context(), uint(id))
	if err != nil {
		renderError(w, r, http.StatusNotFound, fmt.Sprintf("Todo with ID %d not found", id))
		return
	}
	render.JSON(w, r, todo)
}

func (h *todoHandlers) CreateTodo(w http.ResponseWriter, r *http.Request) {
	todoItem := new(todo.Todo)
	err := json.NewDecoder(r.Body).Decode(todoItem)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	todoItem, err = h.store.CreateTodo(r.Context(), todoItem)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, todoItem)
}

func (h *todoHandlers) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, "Invalid ID")
		return
	}
	var existingTodo *todo.Todo
	existingTodo, err = h.store.GetTodoById(r.Context(), uint(id))
	if err != nil {
		renderError(w, r, http.StatusNotFound, fmt.Sprintf("Todo with ID %d not found", id))
		return
	}

	updatedFields := make(map[string]interface{})
	if err := json.NewDecoder(r.Body).Decode(&updatedFields); err != nil {
		renderError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if !maputil.AnyKeys(updatedFields, "name", "date", "description") {
		renderError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	existingTodo, err = h.store.UpdateTodo(r.Context(), existingTodo, updatedFields)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, r, existingTodo)
}

func renderError(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, errors.ErrorResponse{Message: message})
}
