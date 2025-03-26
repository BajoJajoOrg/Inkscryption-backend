package usecase

import (
	"context"
	"mime/multipart"
	"strconv"

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

func (service *UseCase) GetImage(dates []string, name string, id string, ctx context.Context) ([]structures.Canvas, error) {

	if id != "" {
		canvas_id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return []structures.Canvas{}, err
		}
		canvases, err := service.imageStorage.GetById(ctx, canvas_id)
		if err != nil {
			return []structures.Canvas{}, err
		}
		return canvases, nil
	}

	images, err := service.imageStorage.Get(ctx, dates, name)
	if err != nil {
		return []structures.Canvas{}, err
	}

	return images, err
}

func (service *UseCase) AddImage(userCanvas structures.Canvas, ctx context.Context) (int64, error) {

	id, err := service.imageStorage.Add(ctx, userCanvas)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (service *UseCase) AddML(userImage structures.Canvas, img multipart.File, ctx context.Context) error {

	err := service.imageStorage.AddML(ctx, userImage, img)
	if err != nil {
		return err
	}

	return nil
}

func (service *UseCase) DeleteCanvas(canvas structures.Canvas, ctx context.Context) error {
	err := service.imageStorage.Delete(ctx, canvas)
	if err != nil {
		return err
	}

	return nil
}

func (service *UseCase) UpdateCanvas(canvas structures.Canvas, img multipart.File, ctx context.Context) error {

	if canvas.Name != "" {
		err := service.imageStorage.UpdateName(ctx, canvas.Name, canvas.Id)
		if err != nil {
			return nil
		}
	}

	//fmt.Print("wiiide")

	if img != nil {
		err := service.imageStorage.Update(ctx, canvas, img)
		// fmt.Print("Updatin img")
		if err != nil {
			return err
		}
	}

	return nil
}
