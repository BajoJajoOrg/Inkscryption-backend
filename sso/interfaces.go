package sso

import (
	"context"
	"time"

	"github.com/BajoJajoOrg/Inkscryption-backend/config"
	"github.com/go-chi/jwtauth/v5"
)

type Usecase interface {
	Create(ctx context.Context, user User) (*int, error)
	Get(ctx context.Context, email string) (*User, error)
	AddSession(ctx context.Context, jwt string, id string, expiresAt time.Duration) error
	GetSession(ctx context.Context, jwt string, id string) error
	DeleteSession(ctx context.Context, id string) error
}

type Repository interface {
	Create(ctx context.Context, user User) (*int, error)
	Get(ctx context.Context, email string) (*User, error)
}

type RedisRepository interface {
	Add(ctx context.Context, jwt string, id string, expiresAt time.Duration) error
	Get(ctx context.Context, id string) (string, error)
	Delete(ctx context.Context, id string) error
}

type Handlers interface {
	MapHandlers(tokenAuth *jwtauth.JWTAuth) error
	ListenAndServe(cfg config.HTTPServer) error
}
