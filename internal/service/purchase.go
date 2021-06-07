package service

import (
	"context"

	"github.com/JesusG2000/hexsatisfaction/pkg/grpc/api"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/model"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/repository"
	"github.com/pkg/errors"
)

// PurchaseService is a purchase service.
type PurchaseService struct {
	repository.Purchase
	client api.ExistanceClient
}

// NewPurchaseService is a PurchaseService service constructor.
func NewPurchaseService(purchase repository.Purchase, client api.ExistanceClient) *PurchaseService {
	return &PurchaseService{purchase, client}
}

// Create creates new purchase and returns id.
func (p PurchaseService) Create(ctx context.Context, request model.CreatePurchaseRequest) (string, error) {
	purchase := model.PurchaseDTO{
		UserID: request.UserID,
		Date:   request.Date,
		FileID: request.FileID,
	}
	id, err := p.Purchase.Create(ctx, purchase)
	if err != nil {
		return "", errors.Wrap(err, "couldn't create purchase")
	}

	return id, nil
}

// Delete deletes purchase and returns deleted id.
func (p PurchaseService) Delete(ctx context.Context, request model.DeletePurchaseRequest) (string, error) {
	id, err := p.Purchase.Delete(ctx, request.ID)
	if err != nil {
		return "", errors.Wrap(err, "couldn't delete purchase")
	}

	return id, nil
}

// FindByID finds purchase by id.
func (p PurchaseService) FindByID(ctx context.Context, request model.IDPurchaseRequest) (*model.PurchaseDTO, error) {
	purchase, err := p.Purchase.FindByID(ctx, request.ID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find purchase")
	}

	return purchase, nil
}

// FindLastByUserID finds last purchase by user id.
func (p PurchaseService) FindLastByUserID(ctx context.Context, request model.UserIDPurchaseRequest) (*model.PurchaseDTO, error) {
	purchase, err := p.Purchase.FindLastByUserID(ctx, request.ID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find purchase")
	}

	return purchase, nil
}

// FindAllByUserID finds purchases by user id.
func (p PurchaseService) FindAllByUserID(ctx context.Context, request model.UserIDPurchaseRequest) ([]model.PurchaseDTO, error) {
	purchases, err := p.Purchase.FindAllByUserID(ctx, request.ID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find purchases")
	}

	return purchases, nil
}

// FindByUserIDAndPeriod finds purchases by user id and date period.
func (p PurchaseService) FindByUserIDAndPeriod(ctx context.Context, request model.UserIDPeriodPurchaseRequest) ([]model.PurchaseDTO, error) {
	purchases, err := p.Purchase.FindByUserIDAndPeriod(ctx, request.ID, request.Start, request.End)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find purchases")
	}

	return purchases, nil
}

// FindByUserIDAfterDate finds purchases by user id and after date.
func (p PurchaseService) FindByUserIDAfterDate(ctx context.Context, request model.UserIDAfterDatePurchaseRequest) ([]model.PurchaseDTO, error) {
	purchases, err := p.Purchase.FindByUserIDAfterDate(ctx, request.ID, request.Start)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find purchases")
	}

	return purchases, nil
}

// FindByUserIDBeforeDate finds purchases by user id and before date.
func (p PurchaseService) FindByUserIDBeforeDate(ctx context.Context, request model.UserIDBeforeDatePurchaseRequest) ([]model.PurchaseDTO, error) {
	purchases, err := p.Purchase.FindByUserIDBeforeDate(ctx, request.ID, request.End)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find purchases")
	}

	return purchases, nil
}

// FindByUserIDAndFileID finds purchases by user id and file id.
func (p PurchaseService) FindByUserIDAndFileID(ctx context.Context, request model.UserIDFileIDPurchaseRequest) ([]model.PurchaseDTO, error) {
	purchases, err := p.Purchase.FindByUserIDAndFileID(ctx, request.UserID, request.FileID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find purchases")
	}

	return purchases, nil
}

// FindLast finds last purchase.
func (p PurchaseService) FindLast(ctx context.Context) (*model.PurchaseDTO, error) {
	purchase, err := p.Purchase.FindLast(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find purchase")
	}

	return purchase, nil
}

// FindAll finds purchases.
func (p PurchaseService) FindAll(ctx context.Context) ([]model.PurchaseDTO, error) {
	purchases, err := p.Purchase.FindAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find purchases")
	}

	return purchases, nil
}

// FindByPeriod finds purchases by date period.
func (p PurchaseService) FindByPeriod(ctx context.Context, request model.PeriodPurchaseRequest) ([]model.PurchaseDTO, error) {
	purchases, err := p.Purchase.FindByPeriod(ctx, request.Start, request.End)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find purchases")
	}

	return purchases, nil
}

// FindAfterDate finds purchases after date.
func (p PurchaseService) FindAfterDate(ctx context.Context, request model.AfterDatePurchaseRequest) ([]model.PurchaseDTO, error) {
	purchases, err := p.Purchase.FindAfterDate(ctx, request.Start)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find purchases")
	}

	return purchases, nil
}

// FindBeforeDate finds purchases before date.
func (p PurchaseService) FindBeforeDate(ctx context.Context, request model.BeforeDatePurchaseRequest) ([]model.PurchaseDTO, error) {
	purchases, err := p.Purchase.FindBeforeDate(ctx, request.End)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find purchases")
	}

	return purchases, nil
}

// FindByFileID finds purchases by file id.
func (p PurchaseService) FindByFileID(ctx context.Context, request model.FileIDPurchaseRequest) ([]model.PurchaseDTO, error) {
	purchases, err := p.Purchase.FindByFileID(ctx, request.FileID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find purchases")
	}

	return purchases, nil
}
