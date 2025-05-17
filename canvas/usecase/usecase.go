package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"

	"github.com/BajoJajoOrg/Inkscryption-backend/canvas"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/filter"
)

type usecase struct {
	canvasRepo canvas.Repository
	awsRepo    canvas.AwsRepository
	logger     *slog.Logger
}

func New(canvasRepo canvas.Repository, awsRepo canvas.AwsRepository, logger *slog.Logger) canvas.UseCase {
	return &usecase{
		canvasRepo: canvasRepo,
		awsRepo:    awsRepo,
		logger:     logger,
	}
}

func (u *usecase) Create(ctx context.Context, canvas canvas.CanvasBase) (*int, error) {
	return u.canvasRepo.Create(ctx, canvas)
}

func (u *usecase) GetAll(ctx context.Context, filterOptions filter.Options, folder_id int, user_id int) ([]canvas.CanvasBase, error) {
	return u.canvasRepo.GetAll(ctx, filterOptions, folder_id, user_id)
}

func (u *usecase) GetByID(ctx context.Context, id int, userID int) (*canvas.CanvasBase, error) {
	return u.canvasRepo.GetByID(ctx, id, userID)
}

func (u *usecase) Delete(ctx context.Context, id int, userID int) error {
	return u.canvasRepo.Delete(ctx, id, userID)
}

func (u *usecase) Update(ctx context.Context, id int, userID int, url string, file *multipart.File) (*canvas.CanvasBase, error) {

	if err := u.canvasRepo.Update(ctx, id, url); err != nil {
		return nil, err
	}

	canvas, err := u.canvasRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if canvas == nil {
		return nil, fmt.Errorf("canvas was not found")
	}

	return canvas, u.awsRepo.Update(id, file)
}

func (u *usecase) MLUpdate(ctx context.Context, id int, file *multipart.File) error {
	return u.awsRepo.Update(id, file)
}

func (u *usecase) UpdateText(ctx context.Context, id int, text string) error {
	return u.canvasRepo.UpdateText(ctx, id, text)
}

func (u *usecase) ImageToText(ctx context.Context, id int, file *multipart.File) error {
	return u.awsRepo.Update(id, file)
}
