package usecase

import (
	"context"
	"log/slog"

	"github.com/BajoJajoOrg/Inkscryption-backend/sso"
)

type usecase struct {
	userRepo sso.Repository
	logger   *slog.Logger
}

func New(userRepo sso.Repository, logger *slog.Logger) sso.Usecase {
	return &usecase{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (u *usecase) Create(ctx context.Context, user sso.User) (*int, error) {
	return u.userRepo.Create(ctx, user)
}

func (u *usecase) Get(ctx context.Context, email string) (*sso.User, error) {
	return u.userRepo.Get(ctx, email)
}
