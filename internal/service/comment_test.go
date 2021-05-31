package service

import (
	"context"
	"testing"
	"time"

	"github.com/JesusG2000/hexsatisfaction_purchase/internal/model"
	m "github.com/JesusG2000/hexsatisfaction_purchase/internal/service/mock"
	"github.com/pkg/errors"
	testAssert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCommentService_Create(t *testing.T) {
	primitive.NewObjectID().Hex()
	assert := testAssert.New(t)
	type test struct {
		name   string
		req    model.CreateCommentRequest
		fn     func(comment *m.Comment, data test)
		expID  string
		expErr error
	}
	tt := []test{
		{
			name: "Create errors",
			req: model.CreateCommentRequest{
				UserID:     primitive.NewObjectID().Hex(),
				PurchaseID: primitive.NewObjectID().Hex(),
				Date:       time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
				Text:       "some text",
			},
			fn: func(comment *m.Comment, data test) {
				comment.On("Create", mock.Anything, model.CommentDTO{
					UserID:     data.req.UserID,
					PurchaseID: data.req.PurchaseID,
					Date:       data.req.Date,
					Text:       data.req.Text,
				}).
					Return(data.expID, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't create comment"),
		},
		{
			name: "All ok",
			req: model.CreateCommentRequest{
				UserID:     primitive.NewObjectID().Hex(),
				PurchaseID: primitive.NewObjectID().Hex(),
				Date:       time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
				Text:       "some text",
			},
			fn: func(comment *m.Comment, data test) {
				comment.On("Create", mock.Anything, model.CommentDTO{
					UserID:     data.req.UserID,
					PurchaseID: data.req.PurchaseID,
					Date:       data.req.Date,
					Text:       data.req.Text,
				}).
					Return(data.expID, nil)
			},
			expID: primitive.NewObjectID().Hex(),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			comment := new(m.Comment)
			ctx := context.Background()
			service := NewCommentService(comment)
			if tc.fn != nil {
				tc.fn(comment, tc)
			}
			id, err := service.Create(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expID, id)
		})
	}
}

func TestCommentService_Update(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name   string
		req    model.UpdateCommentRequest
		fn     func(comment *m.Comment, data test)
		expID  string
		expErr error
	}
	tt := []test{
		{
			name: "Update errors",
			req: model.UpdateCommentRequest{
				ID:         primitive.NewObjectID().Hex(),
				UserID:     primitive.NewObjectID().Hex(),
				PurchaseID: primitive.NewObjectID().Hex(),
				Date:       time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
				Text:       "some text",
			},
			fn: func(comment *m.Comment, data test) {
				comment.On("Update", mock.Anything, data.req.ID, model.CommentDTO{
					UserID:     data.req.UserID,
					PurchaseID: data.req.PurchaseID,
					Date:       data.req.Date,
					Text:       data.req.Text,
				}).
					Return(data.expID, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't update comment"),
		},
		{
			name: "All ok",
			req: model.UpdateCommentRequest{
				ID:         primitive.NewObjectID().Hex(),
				UserID:     primitive.NewObjectID().Hex(),
				PurchaseID: primitive.NewObjectID().Hex(),
				Date:       time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
				Text:       "some text",
			},
			fn: func(comment *m.Comment, data test) {
				comment.On("Update", mock.Anything, data.req.ID, model.CommentDTO{
					UserID:     data.req.UserID,
					PurchaseID: data.req.PurchaseID,
					Date:       data.req.Date,
					Text:       data.req.Text,
				}).
					Return(data.expID, nil)
			},
			expID: primitive.NewObjectID().Hex(),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			comment := new(m.Comment)
			ctx := context.Background()
			service := NewCommentService(comment)
			if tc.fn != nil {
				tc.fn(comment, tc)
			}
			id, err := service.Update(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expID, id)
		})
	}
}

func TestCommentService_Delete(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name   string
		req    model.DeleteCommentRequest
		fn     func(comment *m.Comment, data test)
		expID  string
		expErr error
	}
	tt := []test{
		{
			name: "Delete errors",
			req: model.DeleteCommentRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(comment *m.Comment, data test) {
				comment.On("Delete", mock.Anything, data.req.ID).
					Return(data.expID, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't delete comment"),
		},
		{
			name: "All ok",
			req: model.DeleteCommentRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(comment *m.Comment, data test) {
				comment.On("Delete", mock.Anything, data.req.ID).
					Return(data.expID, nil)
			},
			expID: primitive.NewObjectID().Hex(),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			comment := new(m.Comment)
			ctx := context.Background()
			service := NewCommentService(comment)
			if tc.fn != nil {
				tc.fn(comment, tc)
			}
			id, err := service.Delete(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expID, id)
		})
	}
}

func TestCommentService_FindByID(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name   string
		req    model.IDCommentRequest
		fn     func(comment *m.Comment, data test)
		exp    *model.CommentDTO
		expErr error
	}
	tt := []test{
		{
			name: "Find errors",
			req: model.IDCommentRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(comment *m.Comment, data test) {
				comment.On("FindByID", mock.Anything, data.req.ID).
					Return(data.exp, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find comment"),
		},
		{
			name: "All ok",
			req: model.IDCommentRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(comment *m.Comment, data test) {
				comment.On("FindByID", mock.Anything, data.req.ID).
					Return(data.exp, nil)
			},
			exp: &model.CommentDTO{
				ID:         primitive.NewObjectID().Hex(),
				UserID:     primitive.NewObjectID().Hex(),
				PurchaseID: primitive.NewObjectID().Hex(),
				Date:       time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
				Text:       "some text",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			comment := new(m.Comment)
			ctx := context.Background()
			service := NewCommentService(comment)
			if tc.fn != nil {
				tc.fn(comment, tc)
			}
			c, err := service.FindByID(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.exp, c)
		})
	}
}

func TestCommentService_FindAllByUserID(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name   string
		req    model.UserIDCommentRequest
		fn     func(comment *m.Comment, data *test)
		exp    []model.CommentDTO
		expErr error
	}
	tt := []test{
		{
			name: "Find errors",
			req: model.UserIDCommentRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(comment *m.Comment, data *test) {
				comment.On("FindAllByUserID", mock.Anything, data.req.ID).
					Return(data.exp, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find comments"),
		},
		{
			name: "All ok",
			req: model.UserIDCommentRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(comment *m.Comment, data *test) {
				for i := range data.exp {
					data.exp[i].UserID = data.req.ID
				}
				comment.On("FindAllByUserID", mock.Anything, data.req.ID).
					Return(data.exp, nil)
			},
			exp: []model.CommentDTO{
				{
					ID:         primitive.NewObjectID().Hex(),
					PurchaseID: primitive.NewObjectID().Hex(),
					Date:       time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					Text:       "some text1",
				},
				{
					ID:         primitive.NewObjectID().Hex(),
					PurchaseID: primitive.NewObjectID().Hex(),
					Date:       time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					Text:       "some text2",
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			comment := new(m.Comment)
			ctx := context.Background()
			service := NewCommentService(comment)
			if tc.fn != nil {
				tc.fn(comment, &tc)
			}
			c, err := service.FindAllByUserID(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.exp, c)
		})
	}
}

func TestCommentService_FindAllByPurchaseID(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name   string
		req    model.PurchaseIDCommentRequest
		fn     func(comment *m.Comment, data *test)
		exp    []model.CommentDTO
		expErr error
	}
	tt := []test{
		{
			name: "Find errors",
			req: model.PurchaseIDCommentRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(comment *m.Comment, data *test) {
				comment.On("FindByPurchaseID", mock.Anything, data.req.ID).
					Return(data.exp, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find comments"),
		},
		{
			name: "All ok",
			req: model.PurchaseIDCommentRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(comment *m.Comment, data *test) {
				for i := range data.exp {
					data.exp[i].PurchaseID = data.req.ID
				}
				comment.On("FindByPurchaseID", mock.Anything, data.req.ID).
					Return(data.exp, nil)
			},
			exp: []model.CommentDTO{
				{
					ID:     primitive.NewObjectID().Hex(),
					UserID: primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					Text:   "some text1",
				},
				{
					ID:     primitive.NewObjectID().Hex(),
					UserID: primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					Text:   "some text2",
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			comment := new(m.Comment)
			ctx := context.Background()
			service := NewCommentService(comment)
			if tc.fn != nil {
				tc.fn(comment, &tc)
			}
			c, err := service.FindByPurchaseID(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.exp, c)
		})
	}
}

func TestCommentService_FindByUserIDAndPurchaseID(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name   string
		req    model.UserPurchaseIDCommentRequest
		fn     func(comment *m.Comment, data *test)
		exp    []model.CommentDTO
		expErr error
	}
	tt := []test{
		{
			name: "Find errors",
			req: model.UserPurchaseIDCommentRequest{
				UserID:     primitive.NewObjectID().Hex(),
				PurchaseID: primitive.NewObjectID().Hex(),
			},
			fn: func(comment *m.Comment, data *test) {
				comment.On("FindByUserIDAndPurchaseID", mock.Anything, data.req.UserID, data.req.PurchaseID).
					Return(data.exp, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find comments"),
		},
		{
			name: "All ok",
			req: model.UserPurchaseIDCommentRequest{
				UserID:     primitive.NewObjectID().Hex(),
				PurchaseID: primitive.NewObjectID().Hex(),
			},
			fn: func(comment *m.Comment, data *test) {
				for i := range data.exp {
					data.exp[i].UserID = data.req.UserID
					data.exp[i].PurchaseID = data.req.PurchaseID
				}
				comment.On("FindByUserIDAndPurchaseID", mock.Anything, data.req.UserID, data.req.PurchaseID).
					Return(data.exp, nil)
			},
			exp: []model.CommentDTO{
				{
					ID:   primitive.NewObjectID().Hex(),
					Date: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					Text: "some text1",
				},
				{
					ID:   primitive.NewObjectID().Hex(),
					Date: time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					Text: "some text2",
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			comment := new(m.Comment)
			ctx := context.Background()
			service := NewCommentService(comment)
			if tc.fn != nil {
				tc.fn(comment, &tc)
			}
			c, err := service.FindByUserIDAndPurchaseID(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.exp, c)
		})
	}
}

func TestCommentService_FindAll(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name   string
		fn     func(comment *m.Comment, data test)
		exp    []model.CommentDTO
		expErr error
	}
	tt := []test{
		{
			name: "Find errors",

			fn: func(comment *m.Comment, data test) {
				comment.On("FindAll", mock.Anything).
					Return(data.exp, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find comments"),
		},
		{
			name: "All ok",
			fn: func(comment *m.Comment, data test) {
				comment.On("FindAll", mock.Anything).
					Return(data.exp, nil)
			},
			exp: []model.CommentDTO{
				{
					ID:         primitive.NewObjectID().Hex(),
					UserID:     primitive.NewObjectID().Hex(),
					PurchaseID: primitive.NewObjectID().Hex(),
					Date:       time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					Text:       "some text1",
				},
				{
					ID:         primitive.NewObjectID().Hex(),
					UserID:     primitive.NewObjectID().Hex(),
					PurchaseID: primitive.NewObjectID().Hex(),
					Date:       time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					Text:       "some text2",
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			comment := new(m.Comment)
			ctx := context.Background()
			service := NewCommentService(comment)
			if tc.fn != nil {
				tc.fn(comment, tc)
			}
			c, err := service.FindAll(ctx)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.exp, c)
		})
	}
}

func TestCommentService_FindByText(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name   string
		req    model.TextCommentRequest
		fn     func(comment *m.Comment, data test)
		exp    []model.CommentDTO
		expErr error
	}
	tt := []test{
		{
			name: "Find errors",
			req: model.TextCommentRequest{
				Text: "some",
			},
			fn: func(comment *m.Comment, data test) {
				comment.On("FindByText", mock.Anything, data.req.Text).
					Return(data.exp, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find comments"),
		},
		{
			name: "All ok",
			req: model.TextCommentRequest{
				Text: "some",
			},
			fn: func(comment *m.Comment, data test) {
				comment.On("FindByText", mock.Anything, data.req.Text).
					Return(data.exp, nil)
			},
			exp: []model.CommentDTO{
				{
					ID:         primitive.NewObjectID().Hex(),
					UserID:     primitive.NewObjectID().Hex(),
					PurchaseID: primitive.NewObjectID().Hex(),
					Date:       time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					Text:       "some text1",
				},
				{
					ID:         primitive.NewObjectID().Hex(),
					UserID:     primitive.NewObjectID().Hex(),
					PurchaseID: primitive.NewObjectID().Hex(),
					Date:       time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					Text:       "some text2",
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			comment := new(m.Comment)
			ctx := context.Background()
			service := NewCommentService(comment)
			if tc.fn != nil {
				tc.fn(comment, tc)
			}
			c, err := service.FindByText(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.exp, c)
		})
	}
}

func TestCommentService_FindByPeriod(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name   string
		req    model.PeriodCommentRequest
		fn     func(comment *m.Comment, data test)
		exp    []model.CommentDTO
		expErr error
	}
	tt := []test{
		{
			name: "Find errors",
			req: model.PeriodCommentRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(comment *m.Comment, data test) {
				comment.On("FindByPeriod", mock.Anything, data.req.Start, data.req.End).
					Return(data.exp, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find comments"),
		},
		{
			name: "All ok",
			req: model.PeriodCommentRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(comment *m.Comment, data test) {
				comment.On("FindByPeriod", mock.Anything, data.req.Start, data.req.End).
					Return(data.exp, nil)
			},
			exp: []model.CommentDTO{
				{
					ID:         primitive.NewObjectID().Hex(),
					UserID:     primitive.NewObjectID().Hex(),
					PurchaseID: primitive.NewObjectID().Hex(),
					Date:       time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					Text:       "some text1",
				},
				{
					ID:         primitive.NewObjectID().Hex(),
					UserID:     primitive.NewObjectID().Hex(),
					PurchaseID: primitive.NewObjectID().Hex(),
					Date:       time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					Text:       "some text2",
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			comment := new(m.Comment)
			ctx := context.Background()
			service := NewCommentService(comment)
			if tc.fn != nil {
				tc.fn(comment, tc)
			}
			c, err := service.FindByPeriod(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.exp, c)
		})
	}
}
