package todo

import (
	"context"
	"time"

	"github.com/ennemli/todo/todo/internal/db"
	"gorm.io/gorm"
)

type Store interface {
	CreateTodo(ctx context.Context, todoItem *Todo) (*Todo, error)
	GetTodos(ctx context.Context) ([]*Todo, error)
	GetTodoById(ctx context.Context, id uint) (*Todo, error)
	DeleteTodoById(ctx context.Context, id uint) (*Todo, error)
	UpdateTodo(ctx context.Context, todoItem *Todo, fields map[string]interface{}) (*Todo, error)
}

// Todo model
type Todo struct {
	gorm.Model
	Date        time.Time `json:"date,omitempty"`
	Name        string    `json:"name" gorm:"index,not null"`
	Description string    `json:"description,omitempty"`
	UserID      uint      `json:"userid" gorm:"index,not null"`
}

type store struct {
	db *gorm.DB
}

func NewStore() Store {
	return &store{
		db: db.GetDB(),
	}
}

func (s *store) CreateTodo(ctx context.Context, todoItem *Todo) (*Todo, error) {
	if err := s.db.WithContext(ctx).Create(todoItem).Error; err != nil {
		return nil, err
	}
	return todoItem, nil
}

func (s *store) GetTodos(ctx context.Context) ([]*Todo, error) {
	todoItems := []*Todo{}
	if err := s.db.WithContext(ctx).Model(&Todo{}).Find(&todoItems).Error; err != nil {
		return nil, err
	}
	return todoItems, nil
}

func (s *store) GetTodoById(ctx context.Context, id uint) (*Todo, error) {
	todoItem := new(Todo)
	if err := s.db.WithContext(ctx).First(todoItem, id).Error; err != nil {
		return nil, err
	}
	return todoItem, nil
}

func (s *store) DeleteTodoById(ctx context.Context, id uint) (*Todo, error) {
	todoItem := new(Todo)
	if err := s.db.WithContext(ctx).Delete(todoItem, id).Error; err != nil {
		return nil, err
	}
	return todoItem, nil
}

func (s *store) UpdateTodo(ctx context.Context, todoItem *Todo, fields map[string]interface{}) (*Todo, error) {
	if err := s.db.WithContext(ctx).Model(todoItem).Updates(fields).Error; err != nil {
		return nil, err
	}
	return todoItem, nil
}
