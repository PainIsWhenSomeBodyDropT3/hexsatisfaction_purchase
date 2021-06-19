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
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFileService_Create(t *testing.T) {
	assert := testAssert.New(t)
	testApi, err := InitTest4Mock()
	require.NoError(t, err)
	type test struct {
		name   string
		req    model.CreateFileRequest
		fn     func(file *m.File, data test)
		expID  string
		expErr error
	}
	tt := []test{
		{
			name: "Create errors",
			req: model.CreateFileRequest{
				Name:        "some",
				Description: "some",
				Size:        1,
				Path:        "some",
				AddDate:     time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				UpdateDate:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				Actual:      true,
				AuthorID:    1,
			},
			fn: func(file *m.File, data test) {
				file.On("Create", mock.Anything, model.FileDTO{
					Name:        data.req.Name,
					Description: data.req.Description,
					Size:        data.req.Size,
					Path:        data.req.Path,
					AddDate:     data.req.AddDate,
					UpdateDate:  data.req.UpdateDate,
					Actual:      data.req.Actual,
					AuthorID:    data.req.AuthorID,
				}).
					Return(data.expID, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't create file"),
		},
		{
			name: "All ok",
			req: model.CreateFileRequest{
				Name:        "some",
				Description: "some",
				Size:        1,
				Path:        "some",
				AddDate:     time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				UpdateDate:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				Actual:      true,
				AuthorID:    1,
			},
			fn: func(file *m.File, data test) {
				file.On("Create", mock.Anything, model.FileDTO{
					Name:        data.req.Name,
					Description: data.req.Description,
					Size:        data.req.Size,
					Path:        data.req.Path,
					AddDate:     data.req.AddDate,
					UpdateDate:  data.req.UpdateDate,
					Actual:      data.req.Actual,
					AuthorID:    data.req.AuthorID,
				}).
					Return(data.expID, nil)
			},
			expID: primitive.NewObjectID().Hex(),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			file := new(m.File)
			ctx := context.Background()
			service := NewFileService(file, testApi.GRPCClient)
			if tc.fn != nil {
				tc.fn(file, tc)
			}
			id, err := service.Create(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expID, id)
		})
	}
}

func TestFileService_Update(t *testing.T) {
	assert := testAssert.New(t)
	testApi, err := InitTest4Mock()
	require.NoError(t, err)
	type test struct {
		name   string
		req    model.UpdateFileRequest
		fn     func(file *m.File, data *test)
		expID  string
		expErr error
	}
	tt := []test{
		{
			name: "Update errors",
			req: model.UpdateFileRequest{
				ID:          primitive.NewObjectID().Hex(),
				Name:        "some",
				Description: "some",
				Size:        1,
				Path:        "some",
				AddDate:     time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				UpdateDate:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				Actual:      true,
				AuthorID:    1,
			},
			fn: func(file *m.File, data *test) {
				file.On("Update", mock.Anything, data.req.ID, model.FileDTO{
					Name:        data.req.Name,
					Description: data.req.Description,
					Size:        data.req.Size,
					Path:        data.req.Path,
					AddDate:     data.req.AddDate,
					UpdateDate:  data.req.UpdateDate,
					Actual:      data.req.Actual,
					AuthorID:    data.req.AuthorID,
				}).
					Return(data.expID, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't update file"),
		},
		{
			name: "All ok",
			req: model.UpdateFileRequest{
				ID:          primitive.NewObjectID().Hex(),
				Name:        "some",
				Description: "some",
				Size:        1,
				Path:        "some",
				AddDate:     time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				UpdateDate:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				Actual:      true,
				AuthorID:    1,
			},
			fn: func(file *m.File, data *test) {
				data.expID = data.req.ID
				file.On("Update", mock.Anything, data.req.ID, model.FileDTO{
					Name:        data.req.Name,
					Description: data.req.Description,
					Size:        data.req.Size,
					Path:        data.req.Path,
					AddDate:     data.req.AddDate,
					UpdateDate:  data.req.UpdateDate,
					Actual:      data.req.Actual,
					AuthorID:    data.req.AuthorID,
				}).
					Return(data.expID, nil)
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			file := new(m.File)
			ctx := context.Background()
			service := NewFileService(file, testApi.GRPCClient)
			if tc.fn != nil {
				tc.fn(file, &tc)
			}
			id, err := service.Update(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expID, id)
		})
	}
}

func TestFileService_Delete(t *testing.T) {
	assert := testAssert.New(t)
	testApi, err := InitTest4Mock()
	require.NoError(t, err)
	type test struct {
		name   string
		req    model.DeleteFileRequest
		fn     func(file *m.File, data *test)
		expID  string
		expErr error
	}
	tt := []test{
		{
			name: "Delete file errors",

			req: model.DeleteFileRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(file *m.File, data *test) {
				file.On("Delete", mock.Anything, data.req.ID).
					Return(data.expID, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't delete file"),
		},
		{
			name: "All ok",
			req: model.DeleteFileRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(file *m.File, data *test) {
				data.expID = data.req.ID
				file.On("Delete", mock.Anything, data.req.ID).
					Return(data.expID, nil)
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			file := new(m.File)
			ctx := context.Background()
			service := NewFileService(file, testApi.GRPCClient)
			if tc.fn != nil {
				tc.fn(file, &tc)
			}
			id, err := service.Delete(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.expID, id)
		})
	}
}

func TestFileService_FindByID(t *testing.T) {
	assert := testAssert.New(t)
	testApi, err := InitTest4Mock()
	require.NoError(t, err)
	type test struct {
		name   string
		req    model.IDFileRequest
		fn     func(file *m.File, data *test)
		exp    *model.FileDTO
		expErr error
	}
	tt := []test{
		{
			name: "Find errors",
			req: model.IDFileRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(file *m.File, data *test) {
				file.On("FindByID", mock.Anything, data.req.ID).
					Return(data.exp, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find file"),
		},
		{
			name: "All ok",
			req: model.IDFileRequest{
				ID: primitive.NewObjectID().Hex(),
			},
			fn: func(file *m.File, data *test) {
				data.exp.ID = data.req.ID
				file.On("FindByID", mock.Anything, data.req.ID).
					Return(data.exp, nil)
			},
			exp: &model.FileDTO{
				Name:        "some",
				Description: "some",
				Size:        1,
				Path:        "some",
				AddDate:     time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				UpdateDate:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				Actual:      true,
				AuthorID:    1,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			file := new(m.File)
			ctx := context.Background()
			service := NewFileService(file, testApi.GRPCClient)
			if tc.fn != nil {
				tc.fn(file, &tc)
			}
			f, err := service.FindByID(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.exp, f)
		})
	}
}

func TestFileService_FindByName(t *testing.T) {
	assert := testAssert.New(t)
	testApi, err := InitTest4Mock()
	require.NoError(t, err)
	type test struct {
		name   string
		req    model.NameFileRequest
		fn     func(file *m.File, data test)
		exp    []model.FileDTO
		expErr error
	}
	tt := []test{
		{
			name: "Find errors",
			req: model.NameFileRequest{
				Name: "some",
			},
			fn: func(file *m.File, data test) {
				file.On("FindByName", mock.Anything, data.req.Name).
					Return(data.exp, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find files"),
		},
		{
			name: "All ok",
			req: model.NameFileRequest{
				Name: "some",
			},
			fn: func(file *m.File, data test) {
				file.On("FindByName", mock.Anything, data.req.Name).
					Return(data.exp, nil)
			},
			exp: []model.FileDTO{
				{
					ID:          primitive.NewObjectID().Hex(),
					Name:        "some",
					Description: "some",
					Size:        1,
					Path:        "some",
					AddDate:     time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					UpdateDate:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					Actual:      true,
					AuthorID:    1,
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			file := new(m.File)
			ctx := context.Background()
			service := NewFileService(file, testApi.GRPCClient)
			if tc.fn != nil {
				tc.fn(file, tc)
			}
			f, err := service.FindByName(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.exp, f)
		})
	}
}

func TestFileService_FindAll(t *testing.T) {
	assert := testAssert.New(t)
	testApi, err := InitTest4Mock()
	require.NoError(t, err)
	type test struct {
		name   string
		fn     func(file *m.File, data test)
		exp    []model.FileDTO
		expErr error
	}
	tt := []test{
		{
			name: "Find errors",

			fn: func(file *m.File, data test) {
				file.On("FindAll", mock.Anything).
					Return(data.exp, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find files"),
		},
		{
			name: "All ok",

			fn: func(file *m.File, data test) {
				file.On("FindAll", mock.Anything).
					Return(data.exp, nil)
			},
			exp: []model.FileDTO{
				{
					ID:          primitive.NewObjectID().Hex(),
					Name:        "some",
					Description: "some",
					Size:        1,
					Path:        "some",
					AddDate:     time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					UpdateDate:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					Actual:      true,
					AuthorID:    1,
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			file := new(m.File)
			ctx := context.Background()
			service := NewFileService(file, testApi.GRPCClient)
			if tc.fn != nil {
				tc.fn(file, tc)
			}
			f, err := service.FindAll(ctx)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.exp, f)
		})
	}
}

func TestFileService_FindByAuthorID(t *testing.T) {
	assert := testAssert.New(t)
	testApi, err := InitTest4Mock()
	require.NoError(t, err)
	type test struct {
		name   string
		req    model.AuthorIDFileRequest
		fn     func(file *m.File, data *test)
		exp    []model.FileDTO
		expErr error
	}
	tt := []test{
		{
			name: "Find errors",
			req: model.AuthorIDFileRequest{
				ID: 1,
			},
			fn: func(file *m.File, data *test) {
				file.On("FindByAuthorID", mock.Anything, data.req.ID).
					Return(data.exp, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find files"),
		},
		{
			name: "All ok",
			req: model.AuthorIDFileRequest{
				ID: 1,
			},
			fn: func(file *m.File, data *test) {
				for i := range data.exp {
					data.exp[i].AuthorID = data.req.ID
				}
				file.On("FindByAuthorID", mock.Anything, data.req.ID).
					Return(data.exp, nil)
			},
			exp: []model.FileDTO{
				{
					ID:          primitive.NewObjectID().Hex(),
					Name:        "some",
					Description: "some",
					Size:        1,
					Path:        "some",
					AddDate:     time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					UpdateDate:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					Actual:      true,
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			file := new(m.File)
			ctx := context.Background()
			service := NewFileService(file, testApi.GRPCClient)
			if tc.fn != nil {
				tc.fn(file, &tc)
			}
			f, err := service.FindByAuthorID(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.exp, f)
		})
	}
}

func TestFileService_FindNotActual(t *testing.T) {
	assert := testAssert.New(t)
	testApi, err := InitTest4Mock()
	require.NoError(t, err)
	type test struct {
		name   string
		fn     func(file *m.File, data test)
		exp    []model.FileDTO
		expErr error
	}
	tt := []test{
		{
			name: "Find errors",

			fn: func(file *m.File, data test) {
				file.On("FindNotActual", mock.Anything).
					Return(data.exp, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find files"),
		},
		{
			name: "All ok",

			fn: func(file *m.File, data test) {
				file.On("FindNotActual", mock.Anything).
					Return(data.exp, nil)
			},
			exp: []model.FileDTO{
				{
					ID:          primitive.NewObjectID().Hex(),
					Name:        "some",
					Description: "some",
					Size:        1,
					Path:        "some",
					AddDate:     time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					UpdateDate:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					Actual:      false,
					AuthorID:    1,
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			file := new(m.File)
			ctx := context.Background()
			service := NewFileService(file, testApi.GRPCClient)
			if tc.fn != nil {
				tc.fn(file, tc)
			}
			f, err := service.FindNotActual(ctx)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.exp, f)
		})
	}
}

func TestFileService_FindActual(t *testing.T) {
	assert := testAssert.New(t)
	testApi, err := InitTest4Mock()
	require.NoError(t, err)
	type test struct {
		name   string
		fn     func(file *m.File, data test)
		exp    []model.FileDTO
		expErr error
	}
	tt := []test{
		{
			name: "Find errors",

			fn: func(file *m.File, data test) {
				file.On("FindActual", mock.Anything).
					Return(data.exp, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find files"),
		},
		{
			name: "All ok",

			fn: func(file *m.File, data test) {
				file.On("FindActual", mock.Anything).
					Return(data.exp, nil)
			},
			exp: []model.FileDTO{
				{
					ID:          primitive.NewObjectID().Hex(),
					Name:        "some",
					Description: "some",
					Size:        1,
					Path:        "some",
					AddDate:     time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					UpdateDate:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					Actual:      true,
					AuthorID:    1,
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			file := new(m.File)
			ctx := context.Background()
			service := NewFileService(file, testApi.GRPCClient)
			if tc.fn != nil {
				tc.fn(file, tc)
			}
			f, err := service.FindActual(ctx)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.exp, f)
		})
	}
}

func TestFileService_FindAddedByPeriod(t *testing.T) {
	assert := testAssert.New(t)
	testApi, err := InitTest4Mock()
	require.NoError(t, err)
	type test struct {
		name   string
		req    model.AddedPeriodFileRequest
		fn     func(file *m.File, data test)
		exp    []model.FileDTO
		expErr error
	}
	tt := []test{
		{
			name: "Find errors",
			req: model.AddedPeriodFileRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(file *m.File, data test) {
				file.On("FindAddedByPeriod", mock.Anything, data.req.Start, data.req.End).
					Return(data.exp, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find files"),
		},
		{
			name: "All ok",
			req: model.AddedPeriodFileRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(file *m.File, data test) {
				file.On("FindAddedByPeriod", mock.Anything, data.req.Start, data.req.End).
					Return(data.exp, nil)
			},
			exp: []model.FileDTO{
				{
					ID:          primitive.NewObjectID().Hex(),
					Name:        "some",
					Description: "some",
					Size:        1,
					Path:        "some",
					AddDate:     time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					UpdateDate:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					Actual:      true,
					AuthorID:    1,
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			file := new(m.File)
			ctx := context.Background()
			service := NewFileService(file, testApi.GRPCClient)
			if tc.fn != nil {
				tc.fn(file, tc)
			}
			f, err := service.FindAddedByPeriod(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.exp, f)
		})
	}
}

func TestFileService_FindUpdatedByPeriod(t *testing.T) {
	assert := testAssert.New(t)
	testApi, err := InitTest4Mock()
	require.NoError(t, err)
	type test struct {
		name   string
		req    model.UpdatedPeriodFileRequest
		fn     func(file *m.File, data test)
		exp    []model.FileDTO
		expErr error
	}
	tt := []test{
		{
			name: "Find errors",
			req: model.UpdatedPeriodFileRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(file *m.File, data test) {
				file.On("FindUpdatedByPeriod", mock.Anything, data.req.Start, data.req.End).
					Return(data.exp, errors.New(""))
			},
			expErr: errors.Wrap(errors.New(""), "couldn't find files"),
		},
		{
			name: "All ok",
			req: model.UpdatedPeriodFileRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(file *m.File, data test) {
				file.On("FindUpdatedByPeriod", mock.Anything, data.req.Start, data.req.End).
					Return(data.exp, nil)
			},
			exp: []model.FileDTO{
				{
					ID:          primitive.NewObjectID().Hex(),
					Name:        "some",
					Description: "some",
					Size:        1,
					Path:        "some",
					AddDate:     time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					UpdateDate:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					Actual:      true,
					AuthorID:    1,
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			file := new(m.File)
			ctx := context.Background()
			service := NewFileService(file, testApi.GRPCClient)
			if tc.fn != nil {
				tc.fn(file, tc)
			}
			f, err := service.FindUpdatedByPeriod(ctx, tc.req)
			if err != nil {
				assert.Equal(tc.expErr.Error(), err.Error())
			}
			assert.Equal(tc.exp, f)
		})
	}
}
