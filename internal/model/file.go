package model

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// File represents a file model.
type File struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Size        int                `bson:"size"`
	Path        string             `bson:"path"`
	AddDate     time.Time          `bson:"addDate"`
	UpdateDate  time.Time          `bson:"updateDate"`
	Actual      bool               `bson:"actual"`
	AuthorID    primitive.ObjectID `bson:"authorID"`
}

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
	AuthorID    string    `json:"authorID"`
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
	}
	var err error
	if f.ID != "" {
		file.ID, err = primitive.ObjectIDFromHex(f.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid id : %v", err)
		}
	}
	if f.AuthorID != "" {
		file.AuthorID, err = primitive.ObjectIDFromHex(f.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("invalid author id : %v", err)
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
		AuthorID:    f.AuthorID.Hex(),
	}

	return &file
}
