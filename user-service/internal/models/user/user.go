package user

import (
	"context"
	"time"

	"github.com/ennemli/todo/user/internal/db"
	"gorm.io/gorm"
)

type Store interface {
	CreateUser(ctx context.Context, userItem *User) (*User, error)
	GetUsers(ctx context.Context) ([]*User, error)
	GetUserById(ctx context.Context, id uint) (*User, error)
	GetUserByName(ctx context.Context, name string) (*User, error)
	DeleteUserById(ctx context.Context, id uint) (*User, error)
	UpdateUser(ctx context.Context, userItem *User, fields map[string]interface{}) (*User, error)
}

// User model
type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"uniqueIndex;not null"`
	Password string `json:"-" gorm:"not null"`
	Todos    []Todo `json:"todos" gorm:"foreignKey:UserID"`
}
type Todo struct {
	gorm.Model
	Date        time.Time
	Name        string
	Description string
	UserID      uint
}

type store struct {
	db *gorm.DB
}

type Credential struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func NewStore() Store {
	return &store{
		db: db.GetDB(),
	}
}

func (s *store) CreateUser(ctx context.Context, userItem *User) (*User, error) {
	if err := s.db.WithContext(ctx).Create(userItem).Error; err != nil {
		return nil, err
	}
	return userItem, nil
}

func (s *store) GetUsers(ctx context.Context) ([]*User, error) {
	userItems := []*User{}
	if err := s.db.WithContext(ctx).Model(&User{}).Preload("Todos").Find(&userItems).Error; err != nil {
		return nil, err
	}
	return userItems, nil
}

func (s *store) GetUserById(ctx context.Context, id uint) (*User, error) {
	userItem := new(User)
	if err := s.db.WithContext(ctx).Preload("Todos").First(userItem, id).Error; err != nil {
		return nil, err
	}
	return userItem, nil
}

func (s *store) GetUserByName(ctx context.Context, name string) (*User, error) {
	userItem := new(User)
	if err := s.db.WithContext(ctx).Preload("Todos").Where("name = ?", name).First(userItem).Error; err != nil {
		return nil, err
	}
	return userItem, nil
}

func (s *store) DeleteUserById(ctx context.Context, id uint) (*User, error) {
	userItem := new(User)
	if err := s.db.WithContext(ctx).Delete(userItem, id).Error; err != nil {
		return nil, err
	}
	return userItem, nil
}

func (s *store) UpdateUser(ctx context.Context, userItem *User, fields map[string]interface{}) (*User, error) {
	if err := s.db.WithContext(ctx).Model(userItem).Updates(fields).Error; err != nil {
		return nil, err
	}
	return userItem, nil
}
