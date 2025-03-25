package images

import (
	"context"
	"mime/multipart"
)

type (
	UseCase interface {
		GetImage(dates []string, name string, id string, ctx context.Context) ([]Canvas, error)
		AddImage(userCanvas Canvas, ctx context.Context) (int64, error)
		DeleteCanvas(canvas Canvas, ctx context.Context) error
		UpdateCanvas(canvas Canvas, img multipart.File, ctx context.Context) error
		AddML(userImage Canvas, img multipart.File, ctx context.Context) error
	}

	ImgStorage interface {
		Get(ctx context.Context, dates []string, name string) ([]Canvas, error)
		GetById(ctx context.Context, id int64) (Canvas, error)
		Add(ctx context.Context, canvas Canvas) (id int64, err error)
		Delete(ctx context.Context, canvas Canvas) error
		Update(ctx context.Context, canvas Canvas, img multipart.File) error
		UpdateName(ctx context.Context, name string, id int64) error
		AddML(ctx context.Context, canvas Canvas, img multipart.File) error
	}
)
