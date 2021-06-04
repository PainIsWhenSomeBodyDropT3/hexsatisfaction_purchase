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

func TestPurchaseService_Create(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name   string
		req    model.CreatePurchaseRequest
		fn     func(purchase *m.Purchase, data test)
		expID  string
		expErr error
	}
	tt := []test{
		{
			name: "Create errors",
			req: model.CreatePurchaseRequest{
				UserID: primitive.NewObjectID().Hex(),
				Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
				FileID: primitive.NewObjectID().Hex(),
			},
			fn: func(purchase *m.Purchase, data test) {
				purchase.On("Create", mock.Anything, model.PurchaseDTO{
					UserID: data.req.UserID,
					Date:   data.req.Date,
					FileID: data.req.FileID,
				}).
					Return(data.expID, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't create purchase"),
		},
		{
			name: "All ok",
			req: model.CreatePurchaseRequest{
				UserID: primitive.NewObjectID().Hex(),
				Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
				FileID: primitive.NewObjectID().Hex(),
			},
			fn: func(purchase *m.Purchase, data test) {
				purchase.On("Create", mock.Anything, model.PurchaseDTO{
					UserID: data.req.UserID,
					Date:   data.req.Date,
					FileID: data.req.FileID,
				}).
					Return(data.expID, nil)
			},
			expID: primitive.NewObjectID().Hex(),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			purchase := new(m.Purchase)
			ctx := context.Background()
			service := NewPurchaseService(purchase)

			if tc.fn != nil {
				tc.fn(purchase, tc)
			}
			id, err := service.Create(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expID, id)
		})
	}
}

func TestPurchaseService_Delete(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name   string
		req    model.DeletePurchaseRequest
		fn     func(purchase *m.Purchase, data *test)
		expID  string
		expErr error
	}
	tt := []test{

		{
			name: "Delete errors",
			req: model.DeletePurchaseRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(purchase *m.Purchase, data *test) {
				purchase.On("Delete", mock.Anything, data.req.ID).
					Return(data.expID, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't delete purchase"),
		},
		{
			name: "All ok",
			req: model.DeletePurchaseRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(purchase *m.Purchase, data *test) {
				data.expID = data.req.ID
				purchase.On("Delete", mock.Anything, data.req.ID).
					Return(data.expID, nil)
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			purchase := new(m.Purchase)
			ctx := context.Background()
			service := NewPurchaseService(purchase)

			if tc.fn != nil {
				tc.fn(purchase, &tc)
			}
			id, err := service.Delete(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expID, id)
		})
	}
}

func TestPurchaseService_FindById(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name        string
		req         model.IDPurchaseRequest
		fn          func(purchase *m.Purchase, data *test)
		expPurchase *model.PurchaseDTO
		expErr      error
	}
	tt := []test{
		{
			name: "FindByID errors",
			req: model.IDPurchaseRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(purchase *m.Purchase, data *test) {
				purchase.On("FindByID", mock.Anything, data.req.ID).
					Return(nil, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find purchase"),
		},
		{
			name: "All ok",
			req: model.IDPurchaseRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(purchase *m.Purchase, data *test) {
				data.expPurchase.ID = data.req.ID
				purchase.On("FindByID", mock.Anything, data.req.ID).
					Return(data.expPurchase, nil)
			},
			expPurchase: &model.PurchaseDTO{
				UserID: primitive.NewObjectID().Hex(),
				Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
				FileID: primitive.NewObjectID().Hex(),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			purchase := new(m.Purchase)
			ctx := context.Background()
			service := NewPurchaseService(purchase)

			if tc.fn != nil {
				tc.fn(purchase, &tc)
			}
			p, err := service.FindByID(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expPurchase, p)
		})
	}
}

func TestPurchaseService_FindLastByUserId(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name        string
		req         model.UserIDPurchaseRequest
		fn          func(purchase *m.Purchase, data *test)
		expPurchase *model.PurchaseDTO
		expErr      error
	}
	tt := []test{
		{
			name: "FindLastByUserID errors",
			req: model.UserIDPurchaseRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(purchase *m.Purchase, data *test) {
				purchase.On("FindLastByUserID", mock.Anything, data.req.ID).
					Return(data.expPurchase, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find purchase"),
		},
		{
			name: "All ok",
			req: model.UserIDPurchaseRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(purchase *m.Purchase, data *test) {
				data.expPurchase.UserID = data.req.ID
				purchase.On("FindLastByUserID", mock.Anything, data.req.ID).
					Return(data.expPurchase, nil)
			},
			expPurchase: &model.PurchaseDTO{
				ID:     primitive.NewObjectID().Hex(),
				Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
				FileID: primitive.NewObjectID().Hex(),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			purchase := new(m.Purchase)
			ctx := context.Background()
			service := NewPurchaseService(purchase)

			if tc.fn != nil {
				tc.fn(purchase, &tc)
			}
			p, err := service.FindLastByUserID(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expPurchase, p)
		})
	}
}

func TestPurchaseService_FindAllByUserId(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name        string
		req         model.UserIDPurchaseRequest
		fn          func(purchase *m.Purchase, data *test)
		expPurchase []model.PurchaseDTO
		expErr      error
	}
	tt := []test{
		{
			name: "FindAllByUserID errors",
			req: model.UserIDPurchaseRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(purchase *m.Purchase, data *test) {
				purchase.On("FindAllByUserID", mock.Anything, data.req.ID).
					Return(data.expPurchase, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find purchases"),
		},
		{
			name: "All ok",
			req: model.UserIDPurchaseRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(purchase *m.Purchase, data *test) {
				for i := range data.expPurchase {
					data.expPurchase[i].UserID = data.req.ID
				}
				purchase.On("FindAllByUserID", mock.Anything, data.req.ID).
					Return(data.expPurchase, nil)
			},
			expPurchase: []model.PurchaseDTO{
				{
					ID:     primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					FileID: primitive.NewObjectID().Hex(),
				},
				{
					ID:     primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					FileID: primitive.NewObjectID().Hex(),
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			purchase := new(m.Purchase)
			ctx := context.Background()
			service := NewPurchaseService(purchase)

			if tc.fn != nil {
				tc.fn(purchase, &tc)
			}
			p, err := service.FindAllByUserID(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expPurchase, p)
		})
	}
}

func TestPurchaseService_FindByUserIdAndPeriod(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name        string
		req         model.UserIDPeriodPurchaseRequest
		fn          func(purchase *m.Purchase, data *test)
		expPurchase []model.PurchaseDTO
		expErr      error
	}
	tt := []test{
		{
			name: "FindByUserIDAndPeriod errors",
			req: model.UserIDPeriodPurchaseRequest{
				ID:    primitive.NewObjectID().Hex(),
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchase *m.Purchase, data *test) {
				purchase.On("FindByUserIDAndPeriod", mock.Anything, data.req.ID, data.req.Start, data.req.End).
					Return(data.expPurchase, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find purchases"),
		},
		{
			name: "All ok",
			req: model.UserIDPeriodPurchaseRequest{
				ID:    primitive.NewObjectID().Hex(),
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchase *m.Purchase, data *test) {
				for i := range data.expPurchase {
					data.expPurchase[i].UserID = data.req.ID
				}
				purchase.On("FindByUserIDAndPeriod", mock.Anything, data.req.ID, data.req.Start, data.req.End).
					Return(data.expPurchase, nil)
			},
			expPurchase: []model.PurchaseDTO{
				{
					ID:     primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.November, 15, 23, 0, 0, 0, time.Local),
					FileID: primitive.NewObjectID().Hex(),
				},
				{
					ID:     primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.December, 3, 23, 0, 0, 0, time.Local),
					FileID: primitive.NewObjectID().Hex(),
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			purchase := new(m.Purchase)
			ctx := context.Background()
			service := NewPurchaseService(purchase)

			if tc.fn != nil {
				tc.fn(purchase, &tc)
			}
			p, err := service.FindByUserIDAndPeriod(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expPurchase, p)
		})
	}
}

func TestPurchaseService_FindByUserIdAfterDate(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name        string
		req         model.UserIDAfterDatePurchaseRequest
		fn          func(purchase *m.Purchase, data *test)
		expPurchase []model.PurchaseDTO
		expErr      error
	}
	tt := []test{
		{
			name: "FindByUserIDAfterDate errors",
			req: model.UserIDAfterDatePurchaseRequest{
				ID:    primitive.NewObjectID().Hex(),
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchase *m.Purchase, data *test) {
				purchase.On("FindByUserIDAfterDate", mock.Anything, data.req.ID, data.req.Start).
					Return(data.expPurchase, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find purchases"),
		},
		{
			name: "All ok",
			req: model.UserIDAfterDatePurchaseRequest{
				ID:    primitive.NewObjectID().Hex(),
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchase *m.Purchase, data *test) {
				for i := range data.expPurchase {
					data.expPurchase[i].UserID = data.req.ID
				}
				purchase.On("FindByUserIDAfterDate", mock.Anything, data.req.ID, data.req.Start).
					Return(data.expPurchase, nil)
			},
			expPurchase: []model.PurchaseDTO{
				{
					ID:     primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					FileID: primitive.NewObjectID().Hex(),
				},
				{
					ID:     primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					FileID: primitive.NewObjectID().Hex(),
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			purchase := new(m.Purchase)
			ctx := context.Background()
			service := NewPurchaseService(purchase)

			if tc.fn != nil {
				tc.fn(purchase, &tc)
			}
			p, err := service.FindByUserIDAfterDate(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expPurchase, p)
		})
	}
}

func TestPurchaseService_FindByUserIdBeforeDate(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name        string
		req         model.UserIDBeforeDatePurchaseRequest
		fn          func(purchase *m.Purchase, data *test)
		expPurchase []model.PurchaseDTO
		expErr      error
	}
	tt := []test{
		{
			name: "FindByUserIDBeforeDate errors",
			req: model.UserIDBeforeDatePurchaseRequest{
				ID:  primitive.NewObjectID().Hex(),
				End: time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchase *m.Purchase, data *test) {
				purchase.On("FindByUserIDBeforeDate", mock.Anything, data.req.ID, data.req.End).
					Return(data.expPurchase, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find purchases"),
		},
		{
			name: "All ok",
			req: model.UserIDBeforeDatePurchaseRequest{
				ID:  primitive.NewObjectID().Hex(),
				End: time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchase *m.Purchase, data *test) {
				for i := range data.expPurchase {
					data.expPurchase[i].UserID = data.req.ID
				}
				purchase.On("FindByUserIDBeforeDate", mock.Anything, data.req.ID, data.req.End).
					Return(data.expPurchase, nil)
			},
			expPurchase: []model.PurchaseDTO{
				{
					ID:     primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					FileID: primitive.NewObjectID().Hex(),
				},
				{
					ID:     primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					FileID: primitive.NewObjectID().Hex(),
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			purchase := new(m.Purchase)
			ctx := context.Background()
			service := NewPurchaseService(purchase)

			if tc.fn != nil {
				tc.fn(purchase, &tc)
			}
			p, err := service.FindByUserIDBeforeDate(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expPurchase, p)
		})
	}
}

func TestPurchaseService_FindByUserIdAndFileID(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name        string
		req         model.UserIDFileIDPurchaseRequest
		fn          func(purchase *m.Purchase, data *test)
		expPurchase []model.PurchaseDTO
		expErr      error
	}
	tt := []test{
		{
			name: "FindByUserIDAndFileID errors",
			req: model.UserIDFileIDPurchaseRequest{
				UserID: primitive.NewObjectID().Hex(),
				FileID: primitive.NewObjectID().Hex(),
			},
			fn: func(purchase *m.Purchase, data *test) {
				purchase.On("FindByUserIDAndFileID", mock.Anything, data.req.UserID, data.req.FileID).
					Return(data.expPurchase, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find purchases"),
		},
		{
			name: "All ok",
			req: model.UserIDFileIDPurchaseRequest{
				UserID: primitive.NewObjectID().Hex(),
				FileID: primitive.NewObjectID().Hex(),
			},
			fn: func(purchase *m.Purchase, data *test) {
				for i := range data.expPurchase {
					data.expPurchase[i].UserID = data.req.UserID
					data.expPurchase[i].FileID = data.req.FileID
				}
				purchase.On("FindByUserIDAndFileID", mock.Anything, data.req.UserID, data.req.FileID).
					Return(data.expPurchase, nil)
			},
			expPurchase: []model.PurchaseDTO{
				{
					ID:   primitive.NewObjectID().Hex(),
					Date: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				},
				{
					ID:   primitive.NewObjectID().Hex(),
					Date: time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			purchase := new(m.Purchase)
			ctx := context.Background()
			service := NewPurchaseService(purchase)

			if tc.fn != nil {
				tc.fn(purchase, &tc)
			}
			p, err := service.FindByUserIDAndFileID(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expPurchase, p)
		})
	}
}

func TestPurchaseService_FindLast(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name        string
		fn          func(purchase *m.Purchase, data test)
		expPurchase *model.PurchaseDTO
		expErr      error
	}
	tt := []test{
		{
			name: "FindLast errors",
			fn: func(purchase *m.Purchase, data test) {
				purchase.On("FindLast", mock.Anything).
					Return(data.expPurchase, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find purchase"),
		},
		{
			name: "All ok",
			fn: func(purchase *m.Purchase, data test) {
				purchase.On("FindLast", mock.Anything).
					Return(data.expPurchase, nil)
			},
			expPurchase: &model.PurchaseDTO{
				ID:     primitive.NewObjectID().Hex(),
				UserID: primitive.NewObjectID().Hex(),
				Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
				FileID: primitive.NewObjectID().Hex(),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			purchase := new(m.Purchase)
			ctx := context.Background()
			service := NewPurchaseService(purchase)

			if tc.fn != nil {
				tc.fn(purchase, tc)
			}
			p, err := service.FindLast(ctx)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expPurchase, p)
		})
	}
}

func TestPurchaseService_FindAll(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name        string
		fn          func(purchase *m.Purchase, data test)
		expPurchase []model.PurchaseDTO
		expErr      error
	}
	tt := []test{
		{
			name: "FindAll errors",
			fn: func(purchase *m.Purchase, data test) {
				purchase.On("FindAll", mock.Anything).
					Return(data.expPurchase, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find purchases"),
		},
		{
			name: "All ok",
			fn: func(purchase *m.Purchase, data test) {
				purchase.On("FindAll", mock.Anything).
					Return(data.expPurchase, nil)
			},
			expPurchase: []model.PurchaseDTO{
				{
					ID:     primitive.NewObjectID().Hex(),
					UserID: primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					FileID: primitive.NewObjectID().Hex(),
				},
				{
					ID:     primitive.NewObjectID().Hex(),
					UserID: primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					FileID: primitive.NewObjectID().Hex(),
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			purchase := new(m.Purchase)
			ctx := context.Background()
			service := NewPurchaseService(purchase)

			if tc.fn != nil {
				tc.fn(purchase, tc)
			}
			p, err := service.FindAll(ctx)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expPurchase, p)
		})
	}
}

func TestPurchaseService_FindByPeriod(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name        string
		req         model.PeriodPurchaseRequest
		fn          func(purchase *m.Purchase, data test)
		expPurchase []model.PurchaseDTO
		expErr      error
	}
	tt := []test{
		{
			name: "FindByPeriod errors",
			req: model.PeriodPurchaseRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchase *m.Purchase, data test) {
				purchase.On("FindByPeriod", mock.Anything, data.req.Start, data.req.End).
					Return(data.expPurchase, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find purchases"),
		},
		{
			name: "All ok",
			req: model.PeriodPurchaseRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchase *m.Purchase, data test) {
				purchase.On("FindByPeriod", mock.Anything, data.req.Start, data.req.End).
					Return(data.expPurchase, nil)
			},
			expPurchase: []model.PurchaseDTO{
				{
					ID:     primitive.NewObjectID().Hex(),
					UserID: primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					FileID: primitive.NewObjectID().Hex(),
				},
				{
					ID:     primitive.NewObjectID().Hex(),
					UserID: primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					FileID: primitive.NewObjectID().Hex(),
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			purchase := new(m.Purchase)
			ctx := context.Background()
			service := NewPurchaseService(purchase)

			if tc.fn != nil {
				tc.fn(purchase, tc)
			}
			p, err := service.FindByPeriod(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expPurchase, p)
		})
	}
}

func TestPurchaseService_FindAfterDate(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name        string
		req         model.AfterDatePurchaseRequest
		fn          func(purchase *m.Purchase, data test)
		expPurchase []model.PurchaseDTO
		expErr      error
	}
	tt := []test{
		{
			name: "FindAfterDate errors",
			req: model.AfterDatePurchaseRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchase *m.Purchase, data test) {
				purchase.On("FindAfterDate", mock.Anything, data.req.Start).
					Return(data.expPurchase, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find purchases"),
		},
		{
			name: "All ok",
			req: model.AfterDatePurchaseRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchase *m.Purchase, data test) {
				purchase.On("FindAfterDate", mock.Anything, data.req.Start).
					Return(data.expPurchase, nil)
			},
			expPurchase: []model.PurchaseDTO{
				{
					ID:     primitive.NewObjectID().Hex(),
					UserID: primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					FileID: primitive.NewObjectID().Hex(),
				},
				{
					ID:     primitive.NewObjectID().Hex(),
					UserID: primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					FileID: primitive.NewObjectID().Hex(),
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			purchase := new(m.Purchase)
			ctx := context.Background()
			service := NewPurchaseService(purchase)

			if tc.fn != nil {
				tc.fn(purchase, tc)
			}
			p, err := service.FindAfterDate(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expPurchase, p)
		})
	}
}

func TestPurchaseService_FindBeforeDate(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name        string
		req         model.BeforeDatePurchaseRequest
		fn          func(purchase *m.Purchase, data test)
		expPurchase []model.PurchaseDTO
		expErr      error
	}
	tt := []test{
		{
			name: "FindBeforeDate errors",
			fn: func(purchase *m.Purchase, data test) {
				purchase.On("FindBeforeDate", mock.Anything, data.req.End).
					Return(data.expPurchase, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find purchases"),
		},
		{
			name: "All ok",
			req: model.BeforeDatePurchaseRequest{
				End: time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchase *m.Purchase, data test) {
				purchase.On("FindBeforeDate", mock.Anything, data.req.End).
					Return(data.expPurchase, nil)
			},
			expPurchase: []model.PurchaseDTO{
				{
					ID:     primitive.NewObjectID().Hex(),
					UserID: primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					FileID: primitive.NewObjectID().Hex(),
				},
				{
					ID:     primitive.NewObjectID().Hex(),
					UserID: primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					FileID: primitive.NewObjectID().Hex(),
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			purchase := new(m.Purchase)
			ctx := context.Background()
			service := NewPurchaseService(purchase)

			if tc.fn != nil {
				tc.fn(purchase, tc)
			}
			p, err := service.FindBeforeDate(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expPurchase, p)
		})
	}
}

func TestPurchaseService_FindByFileID(t *testing.T) {
	assert := testAssert.New(t)
	type test struct {
		name        string
		req         model.FileIDPurchaseRequest
		fn          func(purchase *m.Purchase, data test)
		expPurchase []model.PurchaseDTO
		expErr      error
	}
	tt := []test{
		{
			name: "FindByFileID errors",
			req: model.FileIDPurchaseRequest{
				FileID: primitive.NewObjectID().Hex(),
			},
			fn: func(purchase *m.Purchase, data test) {
				purchase.On("FindByFileID", mock.Anything, data.req.FileID).
					Return(data.expPurchase, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find purchases"),
		},
		{
			name: "All ok",
			req: model.FileIDPurchaseRequest{
				FileID: primitive.NewObjectID().Hex(),
			},
			fn: func(purchase *m.Purchase, data test) {
				for i := range data.expPurchase {
					data.expPurchase[i].FileID = data.req.FileID
				}
				purchase.On("FindByFileID", mock.Anything, data.req.FileID).
					Return(data.expPurchase, nil)
			},
			expPurchase: []model.PurchaseDTO{
				{
					ID:     primitive.NewObjectID().Hex(),
					UserID: primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				},
				{
					ID:     primitive.NewObjectID().Hex(),
					UserID: primitive.NewObjectID().Hex(),
					Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			purchase := new(m.Purchase)
			ctx := context.Background()
			service := NewPurchaseService(purchase)

			if tc.fn != nil {
				tc.fn(purchase, tc)
			}
			p, err := service.FindByFileID(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expPurchase, p)
		})
	}
}
