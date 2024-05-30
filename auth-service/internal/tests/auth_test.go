package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	handlers "github.com/ennemli/todo/todo/internal/handlers/auth"
	"github.com/ennemli/todo/todo/internal/server"
	"github.com/ennemli/todo/todo/pkg/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UserServiceClientMock struct {
	mock.Mock
}

func (m *UserServiceClientMock) GetUser(ctx context.Context, name string) (*handlers.User, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*handlers.User), args.Error(1)
}

func login(T *testing.T, m *UserServiceClientMock) string {
	credential := handlers.Credential{
		Name:     "Guts",
		Password: "123445",
	}
	hashedPassword, _ := crypto.HashPassword(credential.Password)
	expectedUser := &handlers.User{
		ID:       1,
		Name:     credential.Name,
		Password: hashedPassword,
	}
	m.On("GetUser", mock.Anything, credential.Name).Return(expectedUser, nil)

	body, _ := json.Marshal(credential)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	res := MakeRequest(req)
	assert.Equal(T, http.StatusOK, res.Code)
	cookies := res.Result().Cookies()
	assert.NotEmpty(T, cookies)
	var tokenCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			tokenCookie = cookie
			break
		}
	}
	assert.NotNil(T, tokenCookie)

	m.AssertExpectations(T)
	return tokenCookie.Value
}

func TestLogin(T *testing.T) {
	server.InitServerEnv()
	InitServe()
	r := server.GetRouter()
	m := new(UserServiceClientMock)
	authHandler := handlers.NewAuthHandler(m)
	r.Post("/login", authHandler.LoginHandler)
	login(T, m)
}

func TestValidate(T *testing.T) {
	InitServe()
	r := server.GetRouter()
	m := new(UserServiceClientMock)
	authHandler := handlers.NewAuthHandler(m)
	r.Post("/login", authHandler.LoginHandler)
	r.Post("/auth/validate", authHandler.ValidateHandler)
	tokenString := login(T, m)
	req, _ := http.NewRequest("POST", "/auth/validate", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	res := MakeRequest(req)

	assert.Equal(T, http.StatusOK, res.Code)
}

func TestInvalidToken(T *testing.T) {
	tt := []struct {
		name          string
		authorization string
		expected      int
	}{
		{"EmptyAuthorizationHeader", "", http.StatusBadRequest},
		{"InvalidAuthorizationHeader", "ajshahdashda", http.StatusBadRequest},
		{"InvalidFormatAuthorizationHeader", "ashdh ajshahdashda", http.StatusBadRequest},
		{"BearerWithoutToken", "Bearer ", http.StatusUnauthorized},
		{"InvalidToken", "Bearer       4c7b8a3e7d6d6c7c8f3b4c2f1d2d4d7f7c8f3e4c5f6b7c8f3e4c5f6b7c8f3e4c5f6b7c8f3e4c", http.StatusBadRequest},
	}

	InitServe()
	r := server.GetRouter()
	m := new(UserServiceClientMock)
	authHandler := handlers.NewAuthHandler(m)
	r.Post("/auth/validate", authHandler.ValidateHandler)
	for _, tc := range tt {
		T.Run(tc.name, func(T *testing.T) {
			req, _ := http.NewRequest("POST", "/auth/validate", nil)
			req.Header.Set("Authorization", tc.authorization)
			res := MakeRequest(req)
			assert.Equal(T, tc.expected, res.Code)
		})
	}
}

func MakeRequest(req *http.Request) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	r := server.GetRouter()
	r.ServeHTTP(res, req)
	return res
}

func InitServe() {
	server.NewServer()
}
