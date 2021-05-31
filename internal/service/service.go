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
	FindByID(ctx context.Context, request model.IDPurchaseRequest) (*model.Purchase, error)
	FindLastByUserID(ctx context.Context, request model.UserIDPurchaseRequest) (*model.Purchase, error)
	FindAllByUserID(ctx context.Context, request model.UserIDPurchaseRequest) ([]model.Purchase, error)
	FindByUserIDAndPeriod(ctx context.Context, request model.UserIDPeriodPurchaseRequest) ([]model.Purchase, error)
	FindByUserIDAfterDate(ctx context.Context, request model.UserIDAfterDatePurchaseRequest) ([]model.Purchase, error)
	FindByUserIDBeforeDate(ctx context.Context, request model.UserIDBeforeDatePurchaseRequest) ([]model.Purchase, error)
	FindByUserIDAndFileID(ctx context.Context, request model.UserIDFileIDPurchaseRequest) ([]model.Purchase, error)
	FindLast(ctx context.Context) (*model.Purchase, error)
	FindAll(ctx context.Context) ([]model.Purchase, error)
	FindByPeriod(ctx context.Context, request model.PeriodPurchaseRequest) ([]model.Purchase, error)
	FindAfterDate(ctx context.Context, request model.AfterDatePurchaseRequest) ([]model.Purchase, error)
	FindBeforeDate(ctx context.Context, request model.BeforeDatePurchaseRequest) ([]model.Purchase, error)
	FindByFileID(ctx context.Context, request model.FileIDPurchaseRequest) ([]model.Purchase, error)
}

// Comment is an interface for CommentService repository methods.
type Comment interface {
	Create(ctx context.Context, request model.CreateCommentRequest) (string, error)
	Update(ctx context.Context, request model.UpdateCommentRequest) (string, error)
	Delete(ctx context.Context, request model.DeleteCommentRequest) (string, error)
	FindByID(ctx context.Context, request model.IDCommentRequest) (*model.Comment, error)
	FindAllByUserID(ctx context.Context, request model.UserIDCommentRequest) ([]model.Comment, error)
	FindByPurchaseID(ctx context.Context, request model.PurchaseIDCommentRequest) ([]model.Comment, error)
	FindByUserIDAndPurchaseID(ctx context.Context, request model.UserPurchaseIDCommentRequest) ([]model.Comment, error)
	FindAll(ctx context.Context) ([]model.Comment, error)
	FindByText(ctx context.Context, request model.TextCommentRequest) ([]model.Comment, error)
	FindByPeriod(ctx context.Context, request model.PeriodCommentRequest) ([]model.Comment, error)
}

// File is an interface for FileService repository methods.
type File interface {
	Create(ctx context.Context, request model.CreateFileRequest) (string, error)
	Update(ctx context.Context, request model.UpdateFileRequest) (string, error)
	Delete(ctx context.Context, request model.DeleteFileRequest) (string, error)
	FindByID(ctx context.Context, request model.IDFileRequest) (*model.File, error)
	FindByName(ctx context.Context, request model.NameFileRequest) ([]model.File, error)
	FindAll(ctx context.Context) ([]model.File, error)
	FindByAuthorID(ctx context.Context, request model.AuthorIDFileRequest) ([]model.File, error)
	FindNotActual(ctx context.Context) ([]model.File, error)
	FindActual(ctx context.Context) ([]model.File, error)
	FindAddedByPeriod(ctx context.Context, request model.AddedPeriodFileRequest) ([]model.File, error)
	FindUpdatedByPeriod(ctx context.Context, request model.UpdatedPeriodFileRequest) ([]model.File, error)
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
