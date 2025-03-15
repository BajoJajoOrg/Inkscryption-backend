package usecase

import (
	"context"
	"mime/multipart"

	"github.com/BajoJajoOrg/Inkscryption-backend/images"
	structures "github.com/BajoJajoOrg/Inkscryption-backend/images"
	"github.com/BajoJajoOrg/Inkscryption-backend/images/repo"
)

type UseCase struct {
	imageStorage images.ImgStorage
}

func NewImageUseCase(istore images.ImgStorage) *UseCase {
	return &UseCase{
		imageStorage: istore,
	}
}

func GetCore(cfg_sql string) (*UseCase, error) {
	images, err := repo.GetImageRepo(cfg_sql)

	if err != nil {
		return nil, err
	}

	core := UseCase{
		imageStorage: images,
	}
	return &core, nil
}

func (service *UseCase) GetImage(userID int64, ctx context.Context) ([]structures.Canvas, error) {
	images, err := service.imageStorage.Get(ctx, userID)
	if err != nil {
		return []structures.Canvas{}, err
	}

	// if images == "" {
	// 	return structures.Canvas{}, errors.New("no images for user with such sessionID")
	// }

	return images, err
}

func (service *UseCase) AddImage(userImage structures.Canvas, img multipart.File, ctx context.Context) error {

	err := service.imageStorage.Add(ctx, userImage, img)
	if err != nil {
		return err
	}

	return nil
}
