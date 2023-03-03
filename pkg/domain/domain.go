package domain

import (
	"context"
	"errors"
	"time"
)

var ErrWrongCredentials = errors.New("wrong credentials")

const (
	TokenTTL = 12 * time.Hour
)

type User struct {
	ID          int32
	Username    string
	Password    string
	FirstName   string
	LastName    string
	Email       string
	Status      string
	Description string
}

type UserRepo interface {
	CreateUser(ctx context.Context, user *User) (int32, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserById(ctx context.Context, id int32) (*User, error)
	UpdateUserDetails(ctx context.Context, user *User) error
	DeleteUserByUsername(ctx context.Context, username string) error
}

type UserService interface {
	CreateUser(ctx context.Context, user *User) (int32, error)
	GetUser(ctx context.Context, id int32) (*User, error)
	UpdateUserDetails(ctx context.Context, user *User) error
	GenerateToken(ctx context.Context, username string, password string) (string, error)
	ParseToken(ctx context.Context, token string) (int32, error)
}
