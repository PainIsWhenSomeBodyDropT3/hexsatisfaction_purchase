package service

import (
	"context"

	"github.com/JesusG2000/hexsatisfaction/pkg/grpc/api"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/model"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/repository"
	"github.com/pkg/errors"
)

// CommentService is a purchase service.
type CommentService struct {
	repository.Comment
	client api.ExistanceClient
}

// NewCommentService is a CommentService service constructor.
func NewCommentService(comment repository.Comment, client api.ExistanceClient) *CommentService {
	return &CommentService{comment, client}
}

// Create creates comments and returns id.
func (c CommentService) Create(ctx context.Context, request model.CreateCommentRequest) (string, error) {
	var id string
	res, err := c.client.User(ctx, &api.IsUserExistRequest{Id: int32(request.UserID)})
	if err != nil {
		return "", errors.Wrap(err, "couldn't check user existence")
	}

	if res.Exist {
		comment := model.CommentDTO{
			UserID:     request.UserID,
			PurchaseID: request.PurchaseID,
			Date:       request.Date,
			Text:       request.Text,
		}
		id, err = c.Comment.Create(ctx, comment)
		if err != nil {
			return "", errors.Wrap(err, "couldn't create comment")
		}
	}

	return id, nil
}

// Update updates comments and returns id.
func (c CommentService) Update(ctx context.Context, request model.UpdateCommentRequest) (string, error) {
	var id string
	res, err := c.client.User(ctx, &api.IsUserExistRequest{Id: int32(request.UserID)})
	if err != nil {
		return "", errors.Wrap(err, "couldn't check user existence")
	}

	if res.Exist {
		comment := model.CommentDTO{
			UserID:     request.UserID,
			PurchaseID: request.PurchaseID,
			Date:       request.Date,
			Text:       request.Text,
		}
		id, err = c.Comment.Update(ctx, request.ID, comment)
		if err != nil {
			return "", errors.Wrap(err, "couldn't update comment")
		}
	}
	return id, nil
}

// Delete deletes comments and returns id.
func (c CommentService) Delete(ctx context.Context, request model.DeleteCommentRequest) (string, error) {
	id, err := c.Comment.Delete(ctx, request.ID)
	if err != nil {
		return "", errors.Wrap(err, "couldn't delete comment")
	}

	return id, nil
}

// FindByID finds comments by id.
func (c CommentService) FindByID(ctx context.Context, request model.IDCommentRequest) (*model.CommentDTO, error) {
	comment, err := c.Comment.FindByID(ctx, request.ID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find comment")
	}

	return comment, nil
}

// FindAllByUserID finds comments by user id.
func (c CommentService) FindAllByUserID(ctx context.Context, request model.UserIDCommentRequest) ([]model.CommentDTO, error) {
	var comments []model.CommentDTO
	res, err := c.client.User(ctx, &api.IsUserExistRequest{Id: int32(request.ID)})
	if err != nil {
		return nil, errors.Wrap(err, "couldn't check user existence")
	}

	if res.Exist {
		comments, err = c.Comment.FindAllByUserID(ctx, request.ID)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't find comments")
		}

	}
	return comments, nil
}

// FindByPurchaseID finds comments by purchase id.
func (c CommentService) FindByPurchaseID(ctx context.Context, request model.PurchaseIDCommentRequest) ([]model.CommentDTO, error) {
	comments, err := c.Comment.FindByPurchaseID(ctx, request.ID)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find comments")
	}

	return comments, nil
}

// FindByUserIDAndPurchaseID finds comments by purchase and user id.
func (c CommentService) FindByUserIDAndPurchaseID(ctx context.Context, request model.UserPurchaseIDCommentRequest) ([]model.CommentDTO, error) {
	var comments []model.CommentDTO
	res, err := c.client.User(ctx, &api.IsUserExistRequest{Id: int32(request.UserID)})
	if err != nil {
		return nil, errors.Wrap(err, "couldn't check user existence")
	}

	if res.Exist {
		comments, err = c.Comment.FindByUserIDAndPurchaseID(ctx, request.UserID, request.PurchaseID)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't find comments")
		}
	}
	return comments, nil
}

// FindAll finds all comments.
func (c CommentService) FindAll(ctx context.Context) ([]model.CommentDTO, error) {
	comments, err := c.Comment.FindAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find comments")
	}

	return comments, nil
}

// FindByText finds comments by text.
func (c CommentService) FindByText(ctx context.Context, request model.TextCommentRequest) ([]model.CommentDTO, error) {
	comments, err := c.Comment.FindByText(ctx, request.Text)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find comments")
	}

	return comments, nil
}

// FindByPeriod finds comments by date period.
func (c CommentService) FindByPeriod(ctx context.Context, request model.PeriodCommentRequest) ([]model.CommentDTO, error) {
	comments, err := c.Comment.FindByPeriod(ctx, request.Start, request.End)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find comments")
	}

	return comments, nil
}
