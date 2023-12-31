// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mock

import (
	context "context"
	time "time"

	model "github.com/JesusG2000/hexsatisfaction_purchase/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// Comment is an autogenerated mock type for the Comment type
type Comment struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, comment
func (_m *Comment) Create(ctx context.Context, comment model.CommentDTO) (string, error) {
	ret := _m.Called(ctx, comment)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, model.CommentDTO) string); ok {
		r0 = rf(ctx, comment)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.CommentDTO) error); ok {
		r1 = rf(ctx, comment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *Comment) Delete(ctx context.Context, id string) (string, error) {
	ret := _m.Called(ctx, id)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteByPurchaseID provides a mock function with given fields: ctx, id
func (_m *Comment) DeleteByPurchaseID(ctx context.Context, id string) (string, error) {
	ret := _m.Called(ctx, id)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAll provides a mock function with given fields: ctx
func (_m *Comment) FindAll(ctx context.Context) ([]model.CommentDTO, error) {
	ret := _m.Called(ctx)

	var r0 []model.CommentDTO
	if rf, ok := ret.Get(0).(func(context.Context) []model.CommentDTO); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.CommentDTO)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAllByUserID provides a mock function with given fields: ctx, id
func (_m *Comment) FindAllByUserID(ctx context.Context, id int) ([]model.CommentDTO, error) {
	ret := _m.Called(ctx, id)

	var r0 []model.CommentDTO
	if rf, ok := ret.Get(0).(func(context.Context, int) []model.CommentDTO); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.CommentDTO)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByID provides a mock function with given fields: ctx, id
func (_m *Comment) FindByID(ctx context.Context, id string) (*model.CommentDTO, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.CommentDTO
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.CommentDTO); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.CommentDTO)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByPeriod provides a mock function with given fields: ctx, start, end
func (_m *Comment) FindByPeriod(ctx context.Context, start time.Time, end time.Time) ([]model.CommentDTO, error) {
	ret := _m.Called(ctx, start, end)

	var r0 []model.CommentDTO
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, time.Time) []model.CommentDTO); ok {
		r0 = rf(ctx, start, end)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.CommentDTO)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, time.Time, time.Time) error); ok {
		r1 = rf(ctx, start, end)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByPurchaseID provides a mock function with given fields: ctx, id
func (_m *Comment) FindByPurchaseID(ctx context.Context, id string) ([]model.CommentDTO, error) {
	ret := _m.Called(ctx, id)

	var r0 []model.CommentDTO
	if rf, ok := ret.Get(0).(func(context.Context, string) []model.CommentDTO); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.CommentDTO)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByText provides a mock function with given fields: ctx, text
func (_m *Comment) FindByText(ctx context.Context, text string) ([]model.CommentDTO, error) {
	ret := _m.Called(ctx, text)

	var r0 []model.CommentDTO
	if rf, ok := ret.Get(0).(func(context.Context, string) []model.CommentDTO); ok {
		r0 = rf(ctx, text)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.CommentDTO)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, text)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByUserIDAndPurchaseID provides a mock function with given fields: ctx, userID, purchaseID
func (_m *Comment) FindByUserIDAndPurchaseID(ctx context.Context, userID int, purchaseID string) ([]model.CommentDTO, error) {
	ret := _m.Called(ctx, userID, purchaseID)

	var r0 []model.CommentDTO
	if rf, ok := ret.Get(0).(func(context.Context, int, string) []model.CommentDTO); ok {
		r0 = rf(ctx, userID, purchaseID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.CommentDTO)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, string) error); ok {
		r1 = rf(ctx, userID, purchaseID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, id, comment
func (_m *Comment) Update(ctx context.Context, id string, comment model.CommentDTO) (string, error) {
	ret := _m.Called(ctx, id, comment)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, model.CommentDTO) string); ok {
		r0 = rf(ctx, id, comment)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.CommentDTO) error); ok {
		r1 = rf(ctx, id, comment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
