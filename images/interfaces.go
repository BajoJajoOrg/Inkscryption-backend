package images

import (
	"context"
	"mime/multipart"
)

type (
	UseCase interface {
		GetImage(userID int64, ctx context.Context) ([]Canvas, error)
		AddImage(userImage Canvas, img multipart.File, ctx context.Context) error
	}

	ImgStorage interface {
		Get(ctx context.Context, userID int64) ([]Canvas, error)
		Add(ctx context.Context, image Canvas, img multipart.File) error
	}
)
