package folder

import (
	"context"

	"github.com/BajoJajoOrg/Inkscryption-backend/canvas"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/filter"
)

type Repository interface {
	Get(ctx context.Context, filterOptions filter.Options, folder_id int, user_id int) (*canvas.FolderContent, error)
	GetFolder(ctx context.Context, folder_id int, user_id int) (*canvas.FolderBase, error)
	Create(ctx context.Context, folder canvas.FolderBase, user_id int) (*int, error)
	Delete(ctx context.Context, folder_id int) error
	ChangeFolderParent(ctx context.Context, folder_id int, new_parent_id int) error
	ChangeCanvasParent(ctx context.Context, canvas_id int, new_parent_id int) error
}

type UseCase interface {
	Create(ctx context.Context, folder canvas.FolderBase, user_id int) (*int, error)
	Get(ctx context.Context, filterOptions filter.Options, folder_id int, user_id int) (*canvas.FolderContent, error)
	Delete(ctx context.Context, folder_id int) error
	ChangeParent(ctx context.Context, identity string, id int, new_parent_id int) error
}

type Handlers interface {
	MapHandlers() error
}
