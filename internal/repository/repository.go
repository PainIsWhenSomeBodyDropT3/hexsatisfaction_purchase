package repository

import (
	"context"
	"time"

	"github.com/JesusG2000/hexsatisfaction_purchase/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

// Purchase is an interface for PurchaseRepo methods.
type Purchase interface {
	Create(ctx context.Context, purchase model.PurchaseDTO) (string, error)
	Delete(ctx context.Context, id string) (string, error)
	DeleteByFileID(ctx context.Context, id string) (string, error)
	FindByID(ctx context.Context, id string) (*model.PurchaseDTO, error)
	FindLastByUserID(ctx context.Context, id string) (*model.PurchaseDTO, error)
	FindAllByUserID(ctx context.Context, id string) ([]model.PurchaseDTO, error)
	FindByUserIDAndPeriod(ctx context.Context, id string, start, end time.Time) ([]model.PurchaseDTO, error)
	FindByUserIDAfterDate(ctx context.Context, id string, start time.Time) ([]model.PurchaseDTO, error)
	FindByUserIDBeforeDate(ctx context.Context, id string, end time.Time) ([]model.PurchaseDTO, error)
	FindByUserIDAndFileID(ctx context.Context, userID, fileID string) ([]model.PurchaseDTO, error)
	FindLast(ctx context.Context) (*model.PurchaseDTO, error)
	FindAll(ctx context.Context) ([]model.PurchaseDTO, error)
	FindByPeriod(ctx context.Context, start, end time.Time) ([]model.PurchaseDTO, error)
	FindAfterDate(ctx context.Context, start time.Time) ([]model.PurchaseDTO, error)
	FindBeforeDate(ctx context.Context, end time.Time) ([]model.PurchaseDTO, error)
	FindByFileID(ctx context.Context, id string) ([]model.PurchaseDTO, error)
}

// Comment is an interface for CommentRepo methods.
type Comment interface {
	Create(ctx context.Context, comment model.CommentDTO) (string, error)
	Update(ctx context.Context, id string, comment model.CommentDTO) (string, error)
	Delete(ctx context.Context, id string) (string, error)
	DeleteByPurchaseID(ctx context.Context, id string) (string, error)
	FindByID(ctx context.Context, id string) (*model.CommentDTO, error)
	FindAllByUserID(ctx context.Context, id string) ([]model.CommentDTO, error)
	FindByPurchaseID(ctx context.Context, id string) ([]model.CommentDTO, error)
	FindByUserIDAndPurchaseID(ctx context.Context, userID, purchaseID string) ([]model.CommentDTO, error)
	FindAll(ctx context.Context) ([]model.CommentDTO, error)
	FindByText(ctx context.Context, text string) ([]model.CommentDTO, error)
	FindByPeriod(ctx context.Context, start, end time.Time) ([]model.CommentDTO, error)
}

// File is an interface for FileRepo methods.
type File interface {
	Create(ctx context.Context, file model.FileDTO) (string, error)
	Update(ctx context.Context, id string, file model.FileDTO) (string, error)
	Delete(ctx context.Context, id string) (string, error)
	DeleteByAuthorID(ctx context.Context, id string) (string, error)
	FindByID(ctx context.Context, id string) (*model.FileDTO, error)
	FindByName(ctx context.Context, name string) ([]model.FileDTO, error)
	FindAll(ctx context.Context) ([]model.FileDTO, error)
	FindByAuthorID(ctx context.Context, id string) ([]model.FileDTO, error)
	FindNotActual(ctx context.Context) ([]model.FileDTO, error)
	FindActual(ctx context.Context) ([]model.FileDTO, error)
	FindAddedByPeriod(ctx context.Context, start, end time.Time) ([]model.FileDTO, error)
	FindUpdatedByPeriod(ctx context.Context, start, end time.Time) ([]model.FileDTO, error)
}

// Repositories collects all repository interfaces.
type Repositories struct {
	Purchase Purchase
	Comment  Comment
	File     File
}

// NewRepositories is a Repositories constructor.
func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Purchase: NewPurchaseRepo(db),
		Comment:  NewCommentRepo(db),
		File:     NewFileRepo(db),
	}
}
