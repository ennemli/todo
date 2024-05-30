package user

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockUser struct {
	mock.Mock
}

func (m *MockUser) CreateUser(ctx context.Context, userItem *User) (*User, error) {
	args := m.Called(ctx, userItem)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUser) GetUsers(ctx context.Context) ([]*User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*User), args.Error(1)
}

func (m *MockUser) GetUserById(ctx context.Context, id uint) (*User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUser) GetUserByName(ctx context.Context, name string) (*User, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUser) DeleteUserById(ctx context.Context, id uint) (*User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUser) UpdateUser(ctx context.Context, userItem *User, fields map[string]interface{}) (*User, error) {
	args := m.Called(ctx, userItem, fields)
	return args.Get(0).(*User), args.Error(1)
}
