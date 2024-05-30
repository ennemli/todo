package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ennemli/todo/user/internal/errors"
	"github.com/ennemli/todo/user/internal/models/user"
	"github.com/ennemli/todo/user/pkg/crypto"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Handlers interface {
	GetUsers(w http.ResponseWriter, r *http.Request)
	GetUserById(w http.ResponseWriter, r *http.Request)
	GetUserByName(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	DeleteUserById(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	store user.Store
}

func NewUserHandler(store user.Store) Handlers {
	return &userHandler{store: store}
}

func (h *userHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.GetUsers(r.Context())
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.JSON(w, r, users)
}

func (h *userHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		renderError(w, r, http.StatusBadRequest, "invalid id")
		return
	}
	user, err := h.store.GetUserById(r.Context(), uint(id))
	if err != nil {
		renderError(w, r, http.StatusNotFound, fmt.Sprintf("user with id %d not found", id))
		return
	}
	render.JSON(w, r, user)
}

func (h *userHandler) GetUserByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if len(name) <= 0 {
		renderError(w, r, http.StatusBadRequest, "invalid Name")
		return
	}
	userName := bytes.ReplaceAll([]byte(name), []byte("-"), []byte(" "))
	user, err := h.store.GetUserByName(r.Context(), string(userName))
	if err != nil {
		renderError(w, r, http.StatusNotFound, fmt.Sprintf("%s not found", userName))
		return
	}
	render.JSON(w, r, user)
}

func (h *userHandler) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	user := &user.User{}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, "Invalid ID")
		return
	}
	user, err = h.store.DeleteUserById(r.Context(), uint(id))
	if err != nil {
		renderError(w, r, http.StatusNotFound, fmt.Sprintf("User with ID %d not found", id))
		return
	}
	render.JSON(w, r, user)
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var credential user.Credential
	err := json.NewDecoder(r.Body).Decode(&credential)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if len(credential.Password) < 8 {
		renderError(w, r, http.StatusBadRequest, "Password must at least contaians 8 charachters")
		return
	}
	ep, err := crypto.HashPassword(credential.Password)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, "Something went wrong")
		return
	}
	userItem := &user.User{
		Name:     credential.Name,
		Password: ep,
	}
	userItem, err = h.store.CreateUser(r.Context(), userItem)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, "Something Went Wrong")
		return
	}
	render.JSON(w, r, userItem)
}

func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, "Invalid ID")
		return
	}
	var existingUser *user.User
	existingUser, err = h.store.GetUserById(r.Context(), uint(id))
	if err != nil {
		renderError(w, r, http.StatusNotFound, fmt.Sprintf("User with ID %d not found", id))
		return
	}

	updatedFields := make(map[string]interface{})
	if err := json.NewDecoder(r.Body).Decode(&updatedFields); err != nil {
		renderError(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	existingUser, err = h.store.UpdateUser(r.Context(), existingUser, updatedFields)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, r, existingUser)
}

func renderError(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, errors.ErrorResponse{Message: message})
}
