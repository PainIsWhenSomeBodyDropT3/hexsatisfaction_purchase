package service

import (
	"context"

	"github.com/JesusG2000/hexsatisfaction_purchase/internal/model"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/repository"
	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/auth"
)

// Purchase is an interface for PurchaseService repository methods.
type Purchase interface {
	Create(ctx context.Context, request model.CreatePurchaseRequest) (string, error)
	Delete(ctx context.Context, request model.DeletePurchaseRequest) (string, error)
	FindByID(ctx context.Context, request model.IDPurchaseRequest) (*model.PurchaseDTO, error)
	FindLastByUserID(ctx context.Context, request model.UserIDPurchaseRequest) (*model.PurchaseDTO, error)
	FindAllByUserID(ctx context.Context, request model.UserIDPurchaseRequest) ([]model.PurchaseDTO, error)
	FindByUserIDAndPeriod(ctx context.Context, request model.UserIDPeriodPurchaseRequest) ([]model.PurchaseDTO, error)
	FindByUserIDAfterDate(ctx context.Context, request model.UserIDAfterDatePurchaseRequest) ([]model.PurchaseDTO, error)
	FindByUserIDBeforeDate(ctx context.Context, request model.UserIDBeforeDatePurchaseRequest) ([]model.PurchaseDTO, error)
	FindByUserIDAndFileID(ctx context.Context, request model.UserIDFileIDPurchaseRequest) ([]model.PurchaseDTO, error)
	FindLast(ctx context.Context) (*model.PurchaseDTO, error)
	FindAll(ctx context.Context) ([]model.PurchaseDTO, error)
	FindByPeriod(ctx context.Context, request model.PeriodPurchaseRequest) ([]model.PurchaseDTO, error)
	FindAfterDate(ctx context.Context, request model.AfterDatePurchaseRequest) ([]model.PurchaseDTO, error)
	FindBeforeDate(ctx context.Context, request model.BeforeDatePurchaseRequest) ([]model.PurchaseDTO, error)
	FindByFileID(ctx context.Context, request model.FileIDPurchaseRequest) ([]model.PurchaseDTO, error)
}

// Comment is an interface for CommentService repository methods.
type Comment interface {
	Create(ctx context.Context, request model.CreateCommentRequest) (string, error)
	Update(ctx context.Context, request model.UpdateCommentRequest) (string, error)
	Delete(ctx context.Context, request model.DeleteCommentRequest) (string, error)
	FindByID(ctx context.Context, request model.IDCommentRequest) (*model.CommentDTO, error)
	FindAllByUserID(ctx context.Context, request model.UserIDCommentRequest) ([]model.CommentDTO, error)
	FindByPurchaseID(ctx context.Context, request model.PurchaseIDCommentRequest) ([]model.CommentDTO, error)
	FindByUserIDAndPurchaseID(ctx context.Context, request model.UserPurchaseIDCommentRequest) ([]model.CommentDTO, error)
	FindAll(ctx context.Context) ([]model.CommentDTO, error)
	FindByText(ctx context.Context, request model.TextCommentRequest) ([]model.CommentDTO, error)
	FindByPeriod(ctx context.Context, request model.PeriodCommentRequest) ([]model.CommentDTO, error)
}

// File is an interface for FileService repository methods.
type File interface {
	Create(ctx context.Context, request model.CreateFileRequest) (string, error)
	Update(ctx context.Context, request model.UpdateFileRequest) (string, error)
	Delete(ctx context.Context, request model.DeleteFileRequest) (string, error)
	FindByID(ctx context.Context, request model.IDFileRequest) (*model.FileDTO, error)
	FindByName(ctx context.Context, request model.NameFileRequest) ([]model.FileDTO, error)
	FindAll(ctx context.Context) ([]model.FileDTO, error)
	FindByAuthorID(ctx context.Context, request model.AuthorIDFileRequest) ([]model.FileDTO, error)
	FindNotActual(ctx context.Context) ([]model.FileDTO, error)
	FindActual(ctx context.Context) ([]model.FileDTO, error)
	FindAddedByPeriod(ctx context.Context, request model.AddedPeriodFileRequest) ([]model.FileDTO, error)
	FindUpdatedByPeriod(ctx context.Context, request model.UpdatedPeriodFileRequest) ([]model.FileDTO, error)
}

// Services collects all service interfaces.
type Services struct {
	Purchase Purchase
	Comment  Comment
	File     File
}

// Deps represents dependencies for services.
type Deps struct {
	Repos        *repository.Repositories
	TokenManager auth.TokenManager
}

// NewServices is a Services constructor.
func NewServices(deps Deps) *Services {
	return &Services{
		Purchase: NewPurchaseService(deps.Repos.Purchase),
		Comment:  NewCommentService(deps.Repos.Comment),
		File:     NewFileService(deps.Repos.File),
	}
}
