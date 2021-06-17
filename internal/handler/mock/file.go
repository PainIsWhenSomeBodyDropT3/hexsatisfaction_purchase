// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mock

import (
	context "context"

	model "github.com/JesusG2000/hexsatisfaction_purchase/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// File is an autogenerated mock type for the File type
type File struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, request
func (_m *File) Create(ctx context.Context, request model.CreateFileRequest) (string, error) {
	ret := _m.Called(ctx, request)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, model.CreateFileRequest) string); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.CreateFileRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, request
func (_m *File) Delete(ctx context.Context, request model.DeleteFileRequest) (string, error) {
	ret := _m.Called(ctx, request)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, model.DeleteFileRequest) string); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.DeleteFileRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindActual provides a mock function with given fields: ctx
func (_m *File) FindActual(ctx context.Context) ([]model.FileDTO, error) {
	ret := _m.Called(ctx)

	var r0 []model.FileDTO
	if rf, ok := ret.Get(0).(func(context.Context) []model.FileDTO); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.FileDTO)
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

// FindAddedByPeriod provides a mock function with given fields: ctx, request
func (_m *File) FindAddedByPeriod(ctx context.Context, request model.AddedPeriodFileRequest) ([]model.FileDTO, error) {
	ret := _m.Called(ctx, request)

	var r0 []model.FileDTO
	if rf, ok := ret.Get(0).(func(context.Context, model.AddedPeriodFileRequest) []model.FileDTO); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.FileDTO)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.AddedPeriodFileRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAll provides a mock function with given fields: ctx
func (_m *File) FindAll(ctx context.Context) ([]model.FileDTO, error) {
	ret := _m.Called(ctx)

	var r0 []model.FileDTO
	if rf, ok := ret.Get(0).(func(context.Context) []model.FileDTO); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.FileDTO)
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

// FindByAuthorID provides a mock function with given fields: ctx, request
func (_m *File) FindByAuthorID(ctx context.Context, request model.AuthorIDFileRequest) ([]model.FileDTO, error) {
	ret := _m.Called(ctx, request)

	var r0 []model.FileDTO
	if rf, ok := ret.Get(0).(func(context.Context, model.AuthorIDFileRequest) []model.FileDTO); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.FileDTO)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.AuthorIDFileRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByID provides a mock function with given fields: ctx, request
func (_m *File) FindByID(ctx context.Context, request model.IDFileRequest) (*model.FileDTO, error) {
	ret := _m.Called(ctx, request)

	var r0 *model.FileDTO
	if rf, ok := ret.Get(0).(func(context.Context, model.IDFileRequest) *model.FileDTO); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FileDTO)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.IDFileRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByName provides a mock function with given fields: ctx, request
func (_m *File) FindByName(ctx context.Context, request model.NameFileRequest) ([]model.FileDTO, error) {
	ret := _m.Called(ctx, request)

	var r0 []model.FileDTO
	if rf, ok := ret.Get(0).(func(context.Context, model.NameFileRequest) []model.FileDTO); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.FileDTO)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.NameFileRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindNotActual provides a mock function with given fields: ctx
func (_m *File) FindNotActual(ctx context.Context) ([]model.FileDTO, error) {
	ret := _m.Called(ctx)

	var r0 []model.FileDTO
	if rf, ok := ret.Get(0).(func(context.Context) []model.FileDTO); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.FileDTO)
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

// FindUpdatedByPeriod provides a mock function with given fields: ctx, request
func (_m *File) FindUpdatedByPeriod(ctx context.Context, request model.UpdatedPeriodFileRequest) ([]model.FileDTO, error) {
	ret := _m.Called(ctx, request)

	var r0 []model.FileDTO
	if rf, ok := ret.Get(0).(func(context.Context, model.UpdatedPeriodFileRequest) []model.FileDTO); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.FileDTO)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.UpdatedPeriodFileRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, request
func (_m *File) Update(ctx context.Context, request model.UpdateFileRequest) (string, error) {
	ret := _m.Called(ctx, request)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, model.UpdateFileRequest) string); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.UpdateFileRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}