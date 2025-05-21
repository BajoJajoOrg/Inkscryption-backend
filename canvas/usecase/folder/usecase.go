package folder

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/BajoJajoOrg/Inkscryption-backend/canvas"
	"github.com/BajoJajoOrg/Inkscryption-backend/canvas/interfaces/folder"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/filter"
)

type usecase struct {
	folderRepo folder.Repository
	logger     *slog.Logger
}

func New(folderRepo folder.Repository, logger *slog.Logger) folder.UseCase {
	return &usecase{
		folderRepo: folderRepo,
		logger:     logger,
	}
}

func (u *usecase) Create(ctx context.Context, folder canvas.FolderBase, user_id int) (*int, error) {
	return u.folderRepo.Create(ctx, folder, user_id)
}

func (u *usecase) Get(ctx context.Context, filterOptions filter.Options, folder_id int, user_id int) (*canvas.FolderContent, error) {

	// if folder_id == 0 {
	// 	return nil, nil
	// }

	folders, err := u.folderRepo.Get(ctx, filterOptions, folder_id, user_id)
	if err != nil {
		return nil, err
	}

	if folder_id == 0 {
		return folders, nil
	}

	id := *folders.Folder.ID

	breadCrumbs := make([]canvas.BreadCrumb, 0)

	for id != 0 {
		folder, err := u.folderRepo.GetFolder(ctx, id, user_id)
		if err != nil {
			return nil, err
		}

		id = folder.Parent

		print("\n", id, "\n")

		breadCrumb := canvas.BreadCrumb{
			ID:   folder.ID,
			Name: folder.Name,
		}

		breadCrumbs = append(breadCrumbs, breadCrumb)
	}

	folders.BreadCrumbs = breadCrumbs
	return folders, nil
}

func (u *usecase) Delete(ctx context.Context, folder_id int) error {
	return u.folderRepo.Delete(ctx, folder_id)
}

func (u *usecase) ChangeParent(ctx context.Context, identity string, id int, new_parent_id int) error {
	if identity == "canvas" {
		return u.folderRepo.ChangeCanvasParent(ctx, id, new_parent_id)
	} else if identity == "folder" {
		return u.folderRepo.ChangeFolderParent(ctx, id, new_parent_id)
	} else {
		return fmt.Errorf("There is no such identity")
	}
}

func (u *usecase) Update(ctx context.Context, folder_id int, user_id int, name string) error {
	return u.folderRepo.Update(ctx, folder_id, user_id, name)
}
