package images

import (
	"context"
	"mime/multipart"
)

type (
	UseCase interface {
		GetImage(userID int64, dates []string, ctx context.Context) ([]Canvas, error)
		AddImage(userImage Canvas, img multipart.File, ctx context.Context) error
		DeleteCanvas(canvas Canvas, ctx context.Context) error
		UpdateCanvas(canvas Canvas, img multipart.File, ctx context.Context) error
	}

	ImgStorage interface {
		Get(ctx context.Context, userID int64, dates []string) ([]Canvas, error)
		Add(ctx context.Context, image Canvas, img multipart.File) error
		Delete(ctx context.Context, canvas Canvas) error
		Update(ctx context.Context, canvas Canvas, img multipart.File) error
	}
)
