package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	handlers "github.com/ennemli/todo/user/internal/handlers/user"
	"github.com/ennemli/todo/user/internal/models/user"
	"github.com/ennemli/todo/user/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
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

func TestCreateUser(T *testing.T) {
	InitServe()
	mt := new(user.MockUser)
	userHandlers := handlers.NewUserHandler(mt)
	userReq := &user.Credential{
		Name:     "User1",
		Password: "123456789",
	}
	expectedUser := &user.User{
		Name: userReq.Name,
	}
	mt.On("CreateUser", mock.Anything, mock.MatchedBy(func(u interface{}) bool {
		createdUser := u.(*user.User)
		expectedUser.Password = createdUser.Password
		err := bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte(userReq.Password))
		return err == nil
	})).Return(expectedUser, nil)

	reqBody, _ := json.Marshal(userReq)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	r := server.GetRouter()
	r.Post("/", userHandlers.CreateUser)

	res := MakeRequest(req)

	assert.Equal(T, http.StatusOK, res.Code)

	var createdUser user.User
	err := json.NewDecoder(res.Body).Decode(&createdUser)
	assert.Nil(T, err)
	assert.Empty(T, createdUser.Password)
	mt.AssertExpectations(T)
}

func TestGetUsers(T *testing.T) {
	InitServe()
	r := server.GetRouter()
	mt := new(user.MockUser)
	userHandlers := handlers.NewUserHandler(mt)
	expectedUsers := []*user.User{
		{Name: "User 1"},
		{Name: "User 2"},
	}
	expectedUsers[0].ID = 1
	expectedUsers[1].ID = 2
	mt.On("GetUsers", mock.Anything).Return(expectedUsers, nil)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")

	r.Get("/", userHandlers.GetUsers)

	res := MakeRequest(req)

	assert.Equal(T, http.StatusOK, res.Code)

	var users []*user.User
	err := json.NewDecoder(res.Body).Decode(&users)
	assert.Nil(T, err)

	assert.Equal(T, len(expectedUsers), len(users))

	mt.AssertExpectations(T)
}

func TestGetUserById(T *testing.T) {
	InitServe()
	mt := new(user.MockUser)
	userHandlers := handlers.NewUserHandler(mt)
	expectedUser := &user.User{
		Name: "User 1",
	}
	userID := uint(1)
	expectedUser.ID = userID
	mt.On("GetUserById", mock.Anything, userID).Return(expectedUser, nil)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/%d", userID), nil)
	req.Header.Set("Content-Type", "application/json")

	r := server.GetRouter()
	r.Get("/{id:^[0-9]+$}", userHandlers.GetUserById)

	res := MakeRequest(req)

	assert.Equal(T, http.StatusOK, res.Code)

	var userResponse user.User
	err := json.NewDecoder(res.Body).Decode(&userResponse)
	assert.Nil(T, err)

	assert.Equal(T, expectedUser.ID, userResponse.ID)
	assert.Equal(T, expectedUser.Name, userResponse.Name)

	mt.AssertExpectations(T)
}

func TestGetUserByName(T *testing.T) {
	InitServe()
	mt := new(user.MockUser)
	userHandlers := handlers.NewUserHandler(mt)
	expectedUser := &user.User{
		Name: "User One",
	}
	mt.On("GetUserByName", mock.Anything, expectedUser.Name).Return(expectedUser, nil)
	userName := bytes.ReplaceAll([]byte(expectedUser.Name), []byte(" "), []byte("-"))
	req, _ := http.NewRequest("GET", fmt.Sprintf("/%s", userName), nil)
	req.Header.Set("Content-Type", "application/json")

	r := server.GetRouter()
	r.Get("/{name:^[a-zA-Z][a-zA-Z-]+$}", userHandlers.GetUserByName)

	res := MakeRequest(req)

	assert.Equal(T, http.StatusOK, res.Code)

	var userResponse user.User
	err := json.NewDecoder(res.Body).Decode(&userResponse)
	assert.Nil(T, err)

	assert.Equal(T, expectedUser.Name, userResponse.Name)

	mt.AssertExpectations(T)
}

func TestDeleteUserById(T *testing.T) {
	InitServe()
	mt := new(user.MockUser)
	userHandlers := handlers.NewUserHandler(mt)
	expectedUser := &user.User{
		Name: "User 1",
	}
	userID := uint(1)

	expectedUser.ID = userID
	mt.On("DeleteUserById", mock.Anything, userID).Return(expectedUser, nil)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/%d", userID), nil)
	req.Header.Set("Content-Type", "application/json")

	r := server.GetRouter()
	r.Delete("/{id}", userHandlers.DeleteUserById)

	res := MakeRequest(req)

	assert.Equal(T, http.StatusOK, res.Code)

	var userResponse user.User
	err := json.NewDecoder(res.Body).Decode(&userResponse)
	assert.Nil(T, err)

	assert.Equal(T, expectedUser.ID, userResponse.ID)
	assert.Equal(T, expectedUser.Name, userResponse.Name)

	mt.AssertExpectations(T)
}

func TestUpdateUser(T *testing.T) {
	InitServe()
	mt := new(user.MockUser)
	userHandlers := handlers.NewUserHandler(mt)
	userID := uint(1)

	existingUser := &user.User{
		Name: "Original User",
	}
	existingUser.ID = userID
	updatedFields := map[string]interface{}{
		"name": "Updated User",
	}
	expectedUser := *existingUser
	expectedUser.Name = "Updated User"
	mt.On("GetUserById", mock.Anything, userID).Return(existingUser, nil)

	mt.On("UpdateUser", mock.Anything, existingUser, updatedFields).Return(&expectedUser, nil).Run(func(args mock.Arguments) {
		t := args.Get(1).(*user.User)
		fields := args.Get(2).(map[string]interface{})
		t.Name = fields["name"].(string)
	})

	reqBody, _ := json.Marshal(updatedFields)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/%d", userID), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	r := server.GetRouter()
	r.Put("/{id}", userHandlers.UpdateUser)

	res := MakeRequest(req)

	assert.Equal(T, http.StatusOK, res.Code)

	var updatedUser user.User
	err := json.NewDecoder(res.Body).Decode(&updatedUser)
	assert.Nil(T, err)

	assert.Equal(T, existingUser.ID, updatedUser.ID)
	assert.Equal(T, updatedFields["name"], updatedUser.Name)

	mt.AssertExpectations(T)
}
