package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	"example.com/pkg/domain"
	"github.com/golang-jwt/jwt"
)

const (
	signingKey = "pudgebooster"
	tokenTTL   = 12 * time.Hour
)

type Controller struct {
	repo domain.UserRepo
}

type tokenClaims struct {
	jwt.StandardClaims
	UserID int32 `json:"user_id"`
}

func NewController(repo domain.UserRepo) (*Controller, error) {
	return &Controller{repo: repo}, nil
}

func (c *Controller) CreateUser(ctx context.Context, user *domain.User) (int32, error) {
	return c.repo.CreateUser(ctx, user)
}

func (c *Controller) GetUser(ctx context.Context, id int32) (*domain.User, error) {
	return c.repo.GetUserById(ctx, id)
}

func (c *Controller) UpdateUserDetails(ctx context.Context, user *domain.User) error {
	return c.repo.UpdateUserDetails(ctx, user)
}

func (c *Controller) GenerateToken(ctx context.Context, username string, password string) (string, error) {
	user, err := c.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("error getting user: %w", err)
	}
	if user.Username != username || user.Password != password {
		return "", domain.ErrWrongCredentials
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})
	return token.SignedString([]byte(signingKey))
}

func (c *Controller) ParseToken(ctx context.Context, token string) (int32, error) {
	accessToken, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := accessToken.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}
	return claims.UserID, nil
}
