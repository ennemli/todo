package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/ennemli/todo/todo/internal/db"
	"github.com/ennemli/todo/todo/internal/errors"
	"github.com/ennemli/todo/todo/pkg/jwt"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type Handlers interface {
	LoginHandler(w http.ResponseWriter, r *http.Request)
	ValidateHandler(w http.ResponseWriter, r *http.Request)
}
type UserSerivce interface {
	GetUser(ctx context.Context, name string) (*User, error)
}

type AuthHandler struct {
	userService UserSerivce
}
type UserServiceClient struct{}

func NewAuthHandler(u UserSerivce) Handlers {
	return &AuthHandler{
		userService: u,
	}
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"-"`
	ID       uint   `json:"id"`
}

func (u *UserServiceClient) GetUser(ctx context.Context, name string) (*User, error) {
	user := new(User)
	if err := db.GetDB().Where(" name = ? ", name).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

type Credential struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type ResponseMessage struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

func (a *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	credential := new(Credential)
	if err := json.NewDecoder(r.Body).Decode(credential); err != nil || credential.Name == "" || credential.Password == "" {
		renderError(w, r, http.StatusBadRequest, "Bad Request")
		return
	}
	user, err := a.userService.GetUser(r.Context(), credential.Name)
	if err != nil {
		renderError(w, r, http.StatusNotFound, "There is no user with this name")
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credential.Password)); err != nil {

		renderError(w, r, http.StatusUnauthorized, "Please check if your name and password are correct")
		return
	}
	exp := time.Now().Add(time.Hour * 3)
	token, err := jwt.Create(user, exp.Unix())
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, "Something Went Wrong")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: exp,
	})
	render.JSON(w, r, ResponseMessage{
		User:  user,
		Token: token,
	})
}

func (a *AuthHandler) ValidateHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	token := strings.Split(authHeader, " ")
	if len(token) != 2 || token[0] != "Bearer" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	tokenStr := token[1]
	_, err := jwt.Validate(tokenStr)
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	http.Error(w, "", http.StatusOK)
}

func renderError(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, errors.ErrorResponse{Message: message})
}
