package canvas

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/BajoJajoOrg/Inkscryption-backend/config"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/filter"
)

type Repository interface {
	GetAll(ctx context.Context, filterOptions filter.Options) ([]CanvasBase, error)
	GetByID(ctx context.Context, id int) (*CanvasBase, error)
	Create(ctx context.Context, canvas CanvasBase) (*int, error)
	Delete(ctx context.Context, id int) error
}

type AwsRepository interface {
	Update(id int, file *multipart.File) error
}

type UseCase interface {
	GetAll(ctx context.Context, filterOptions filter.Options) ([]CanvasBase, error)
	GetByID(ctx context.Context, id int) (*CanvasBase, error)
	Create(ctx context.Context, canvas CanvasBase) (*int, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, file *multipart.File) (*CanvasBase, error)
	ImageToText(ctx context.Context, id int, file *multipart.File) error
}

type Handlers interface {
	MapHandlers() error
	ListenAndServe(cfg config.HTTPServer) error
	Create(w http.ResponseWriter, r *http.Request)
}
