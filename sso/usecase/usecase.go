package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/BajoJajoOrg/Inkscryption-backend/sso"
)

type usecase struct {
	userRepo  sso.Repository
	redisRepo sso.RedisRepository
	logger    *slog.Logger
}

func New(userRepo sso.Repository, redisRepo sso.RedisRepository, logger *slog.Logger) sso.Usecase {
	return &usecase{
		userRepo:  userRepo,
		redisRepo: redisRepo,
		logger:    logger,
	}
}

func (u *usecase) Create(ctx context.Context, user sso.User) (*int, error) {
	return u.userRepo.Create(ctx, user)
}

func (u *usecase) Get(ctx context.Context, email string) (*sso.User, error) {
	return u.userRepo.Get(ctx, email)
}

func (u *usecase) AddSession(ctx context.Context, jwt string, id string) error {
	return u.redisRepo.Add(ctx, jwt, id)
}

func (u *usecase) GetSession(ctx context.Context, jwt string, id string) error {
	res, err := u.redisRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	if res != jwt {
		return fmt.Errorf("wrong refresh token")
	}

	return nil
}

func (u *usecase) DeleteSession(ctx context.Context, id string) error {
	return u.redisRepo.Delete(ctx, id)
}
