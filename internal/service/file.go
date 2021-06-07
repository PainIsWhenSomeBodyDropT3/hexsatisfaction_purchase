package service

import (
	"context"

	"github.com/JesusG2000/hexsatisfaction/pkg/grpc/api"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/model"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/repository"
	"github.com/pkg/errors"
)

// FileService is a file service.
type FileService struct {
	repository.File
}

// NewFileService is a FileService service constructor.
func NewFileService(file repository.File, client api.ExistanceClient) *FileService {
	return &FileService{file}
}

// Create creates new file and returns id.
func (f FileService) Create(ctx context.Context, request model.CreateFileRequest) (string, error) {
	file := model.FileDTO{
		Name:        request.Name,
		Description: request.Description,
		Size:        request.Size,
		Path:        request.Path,
		AddDate:     request.AddDate,
		UpdateDate:  request.UpdateDate,
		Actual:      request.Actual,
		AuthorID:    request.AuthorID,
	}
	id, err := f.File.Create(ctx, file)
	if err != nil {
		return "", errors.Wrap(err, "couldn't create file")
	}

	return id, nil
}

// Update updates file and returns id.
func (f FileService) Update(ctx context.Context, request model.UpdateFileRequest) (string, error) {
	file := model.FileDTO{
		Name:        request.Name,
		Description: request.Description,
		Size:        request.Size,
		Path:        request.Path,
		AddDate:     request.AddDate,
		UpdateDate:  request.UpdateDate,
		Actual:      request.Actual,
		AuthorID:    request.AuthorID,
	}
	id, err := f.File.Update(ctx, request.ID, file)
	if err != nil {
		return "", errors.Wrap(err, "couldn't update file")
	}

	return id, nil
}

// Delete deletes file and returns deleted id.
func (f FileService) Delete(ctx context.Context, request model.DeleteFileRequest) (string, error) {
	id, err := f.File.Delete(ctx, request.ID)
	if err != nil {
		return "", errors.Wrap(err, "couldn't delete file")
	}

	return id, nil
}

// FindByID finds file by id.
func (f FileService) FindByID(ctx context.Context, request model.IDFileRequest) (*model.FileDTO, error) {
	file, err := f.File.FindByID(ctx, request.ID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find file")
	}

	return file, nil
}

// FindByName finds files by name.
func (f FileService) FindByName(ctx context.Context, request model.NameFileRequest) ([]model.FileDTO, error) {
	files, err := f.File.FindByName(ctx, request.Name)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find files")
	}

	return files, nil
}

// FindAll finds files.
func (f FileService) FindAll(ctx context.Context) ([]model.FileDTO, error) {
	files, err := f.File.FindAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find files")
	}

	return files, nil
}

// FindByAuthorID finds files by author id.
func (f FileService) FindByAuthorID(ctx context.Context, request model.AuthorIDFileRequest) ([]model.FileDTO, error) {
	files, err := f.File.FindByAuthorID(ctx, request.ID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find files")
	}

	return files, nil
}

// FindNotActual finds not actual files.
func (f FileService) FindNotActual(ctx context.Context) ([]model.FileDTO, error) {
	files, err := f.File.FindNotActual(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find files")
	}

	return files, nil
}

// FindActual finds actual files.
func (f FileService) FindActual(ctx context.Context) ([]model.FileDTO, error) {
	files, err := f.File.FindActual(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find files")
	}

	return files, nil
}

// FindAddedByPeriod finds added files by date period.
func (f FileService) FindAddedByPeriod(ctx context.Context, request model.AddedPeriodFileRequest) ([]model.FileDTO, error) {
	files, err := f.File.FindAddedByPeriod(ctx, request.Start, request.End)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find files")
	}

	return files, nil
}

// FindUpdatedByPeriod finds updated files by date period.
func (f FileService) FindUpdatedByPeriod(ctx context.Context, request model.UpdatedPeriodFileRequest) ([]model.FileDTO, error) {
	files, err := f.File.FindUpdatedByPeriod(ctx, request.Start, request.End)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find files")
	}

	return files, nil
}
