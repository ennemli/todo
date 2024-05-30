package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/ennemli/todo/user/internal/errors"
	handlers "github.com/ennemli/todo/user/internal/handlers/user"
	"github.com/ennemli/todo/user/internal/middlewares"
	"github.com/ennemli/todo/user/internal/models/user"
	"github.com/ennemli/todo/user/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTimeOut(T *testing.T) {
	InitServe()
	r := server.GetRouter()
	r.Use(middlewares.SetTimeOut(errors.TimeoutDuration))
	mt := new(user.MockUser)
	userHandlers := handlers.NewUserHandler(mt)
	r.Get("/", userHandlers.GetUsers)
	mt.On("GetUsers", mock.Anything).Return([]*user.User{}, nil).Run(func(args mock.Arguments) {
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

func TestPasswordStrength(T *testing.T) {
	InitServe()
	mt := new(user.MockUser)
	userHandlers := handlers.NewUserHandler(mt)
	userReq := &user.Credential{
		Name:     "User 1",
		Password: "123",
	}

	mt.On("CreateUser", mock.Anything, mock.Anything).Return(mock.Anything, nil)

	reqBody, _ := json.Marshal(userReq)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	r := server.GetRouter()
	r.Post("/", userHandlers.CreateUser)

	res := MakeRequest(req)

	assert.Equal(T, http.StatusBadRequest, res.Code)

	var errorRes errors.ErrorResponse
	err := json.NewDecoder(res.Body).Decode(&errorRes)
	assert.Nil(T, err)

	assert.Equal(T, "Password must at least contaians 8 charachters", errorRes.Message)
}
