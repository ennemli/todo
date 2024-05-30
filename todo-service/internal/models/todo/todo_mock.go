package todo

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockTodo struct {
	mock.Mock
}

func (m *MockTodo) CreateTodo(ctx context.Context, todoItem *Todo) (*Todo, error) {
	args := m.Called(ctx, todoItem)
	return args.Get(0).(*Todo), args.Error(1)
}

func (m *MockTodo) GetTodos(ctx context.Context) ([]*Todo, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*Todo), args.Error(1)
}

func (m *MockTodo) GetTodoById(ctx context.Context, id uint) (*Todo, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*Todo), args.Error(1)
}

func (m *MockTodo) DeleteTodoById(ctx context.Context, id uint) (*Todo, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*Todo), args.Error(1)
}

func (m *MockTodo) UpdateTodo(ctx context.Context, todoItem *Todo, fields map[string]interface{}) (*Todo, error) {
	args := m.Called(ctx, todoItem, fields)
	return args.Get(0).(*Todo), args.Error(1)
}
