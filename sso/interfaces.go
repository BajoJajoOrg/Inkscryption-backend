package sso

import (
	"context"

	"github.com/BajoJajoOrg/Inkscryption-backend/config"
)

type Usecase interface {
	Create(ctx context.Context, user User) (*int, error)
	Get(ctx context.Context, email string) (*User, error)
}

type Repository interface {
	Create(ctx context.Context, user User) (*int, error)
	Get(ctx context.Context, email string) (*User, error)
}

type Handlers interface {
	MapHandlers() error
	ListenAndServe(cfg config.HTTPServer) error
}
