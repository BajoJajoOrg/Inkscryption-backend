package canvas

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/BajoJajoOrg/Inkscryption-backend/config"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/filter"
)

type Repository interface {
	GetAll(ctx context.Context, filterOptions filter.Options, folder_id int, user_id int) ([]CanvasBase, error)
	GetByID(ctx context.Context, id int, userID int) (*CanvasBase, error)
	Create(ctx context.Context, canvas CanvasBase) (*int, error)
	Delete(ctx context.Context, id int, userID int) error
	Update(ctx context.Context, id int, url string) error
	UpdateName(ctx context.Context, id int, name string) error
	UpdateText(ctx context.Context, id int, text string) error
	ChangeParent(ctx context.Context, canvas_id int, new_parent_id int) error
}

type AwsRepository interface {
	Update(id int, file *multipart.File) error
}

type UseCase interface {
	GetAll(ctx context.Context, filterOptions filter.Options, folder_id int, user_id int) ([]CanvasBase, error)
	GetByID(ctx context.Context, id int, userID int) (*CanvasBase, error)
	Create(ctx context.Context, canvas CanvasBase) (*int, error)
	Delete(ctx context.Context, id int, userID int) error
	Update(ctx context.Context, id int, userID int, url string, name string, file *multipart.File) (*CanvasBase, error)
	MLUpdate(ctx context.Context, id int, file *multipart.File) error
	UpdateText(ctx context.Context, id int, text string) error
	ImageToText(ctx context.Context, id int, file *multipart.File) error
}

type Handlers interface {
	MapHandlers() error
	ListenAndServe(cfg config.HTTPServer) error
	Create(w http.ResponseWriter, r *http.Request)
}
