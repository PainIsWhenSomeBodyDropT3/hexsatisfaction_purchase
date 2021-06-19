package model

import (
	"time"

	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/database/mongo"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Files represents a slice of a file model.
type Files []File

// FilesDTO represents a slice of a dto file model.
type FilesDTO []FileDTO

// File represents a file model.
type File mongo.File

// FileDTO represents dto of a file model.
type FileDTO struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Size        int       `json:"size"`
	Path        string    `json:"path"`
	AddDate     time.Time `json:"addDate"`
	UpdateDate  time.Time `json:"updateDate"`
	Actual      bool      `json:"actual"`
	AuthorID    int       `json:"authorID"`
}

// Entity converts FileDTO to File.
func (f FileDTO) Entity() (*File, error) {
	file := File{
		Name:        f.Name,
		Description: f.Description,
		Size:        f.Size,
		Path:        f.Path,
		AddDate:     f.AddDate,
		UpdateDate:  f.UpdateDate,
		Actual:      f.Actual,
		AuthorID:    f.AuthorID,
	}
	var err error
	if f.ID != "" {
		file.ID, err = primitive.ObjectIDFromHex(f.ID)
		if err != nil {
			return nil, errors.Wrap(err, "invalid id")
		}
	}

	return &file, nil
}

// DTO converts File to FileDTO.
func (f File) DTO() *FileDTO {
	file := FileDTO{
		ID:          f.ID.Hex(),
		Name:        f.Name,
		Description: f.Description,
		Size:        f.Size,
		Path:        f.Path,
		AddDate:     f.AddDate,
		UpdateDate:  f.UpdateDate,
		Actual:      f.Actual,
		AuthorID:    f.AuthorID,
	}

	return &file
}

// Entity converts FilesDTO to Files.
func (f FilesDTO) Entity() (Files, error) {
	var files Files
	for _, file := range f {
		entityFile, err := file.Entity()
		if err != nil {
			return nil, errors.Wrap(err, "couldn't convert file dto to file")
		}
		files = append(files, *entityFile)
	}
	return files, nil
}

// DTO converts Files to FilesDTO
func (f Files) DTO() FilesDTO {
	var files FilesDTO
	for _, file := range f {
		files = append(files, *file.DTO())
	}
	return files
}
