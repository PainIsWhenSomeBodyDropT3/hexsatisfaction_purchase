package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/JesusG2000/hexsatisfaction_purchase/internal/config"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/model"
	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/database/mongo"
	assertTest "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Connect2FileMongo() (context.Context, *FileRepo, error) {
	ctx := context.Background()
	cfg, err := config.Init(configPath)
	if err != nil {
		return nil, nil, err
	}

	db, err := mongo.NewMongo(ctx, cfg.Mongo)
	if err != nil {
		return nil, nil, err
	}

	return ctx, NewFileRepo(db), nil
}

func TestFileRepo_Create(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2FileMongo()
	require.NoError(t, err)
	type test struct {
		name string
		file model.FileDTO
	}
	tt := []test{
		{
			name: "all ok",
			file: model.FileDTO{
				Name:        "some",
				Description: "some",
				Size:        1,
				Path:        "some",
				AddDate:     time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
				UpdateDate:  time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
				Actual:      false,
				AuthorID:    primitive.NewObjectID().Hex(),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var fileID string
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
			fileID, err = repo.Create(ctx, tc.file)

			assert.NotEmpty(fileID)
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestFileRepo_Update(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2FileMongo()
	require.NoError(t, err)
	type test struct {
		name   string
		isOk   bool
		id     string
		file   model.FileDTO
		expErr error
	}
	tt := []test{
		{
			name:   "not correct userID",
			expErr: errors.New("the provided hex string is not a valid ObjectID"),
		},
		{
			name:   "not found",
			id:     primitive.NewObjectID().Hex(),
			expErr: errors.New("mongo: no documents in result"),
		},
		{
			name: "all ok",
			isOk: true,
			file: model.FileDTO{
				Name:        "some",
				Description: "some",
				Size:        1,
				Path:        "some",
				AddDate:     time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
				UpdateDate:  time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
				Actual:      false,
				AuthorID:    primitive.NewObjectID().Hex(),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			fileID := tc.id
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
			if tc.isOk {
				fileID, err = repo.Create(ctx, tc.file)
			}
			id, err := repo.Update(ctx, fileID, tc.file)
			assert.Equal(tc.expErr, err)
			if tc.isOk {
				assert.Equal(fileID, id)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestFileRepo_Delete(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2FileMongo()
	require.NoError(t, err)
	type test struct {
		name   string
		isOk   bool
		id     string
		file   model.FileDTO
		expErr error
	}
	tt := []test{
		{
			name:   "not correct userID",
			expErr: errors.New("the provided hex string is not a valid ObjectID"),
		},
		{
			name:   "not found",
			id:     primitive.NewObjectID().Hex(),
			expErr: errors.New("mongo: no documents in result"),
		},
		{
			name: "all ok",
			isOk: true,
			file: model.FileDTO{
				Name:        "some",
				Description: "some",
				Size:        1,
				Path:        "some",
				AddDate:     time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
				UpdateDate:  time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
				Actual:      false,
				AuthorID:    primitive.NewObjectID().Hex(),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			fileID := tc.id
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
			if tc.isOk {
				fileID, err = repo.Create(ctx, tc.file)
			}
			id, err := repo.Delete(ctx, fileID)
			assert.Equal(tc.expErr, err)
			if tc.isOk {
				assert.Equal(fileID, id)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestFileRepo_DeleteByAuthorID(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2FileMongo()
	require.NoError(t, err)
	type test struct {
		name   string
		isOk   bool
		id     string
		fn     func(data *test)
		file   model.FileDTO
		expErr error
	}
	tt := []test{
		{
			name:   "not correct userID",
			expErr: errors.New("the provided hex string is not a valid ObjectID"),
		},
		{
			name:   "not found",
			id:     primitive.NewObjectID().Hex(),
			expErr: errors.New("mongo: no documents in result"),
		},
		{
			name: "all ok",
			isOk: true,
			id:   primitive.NewObjectID().Hex(),
			fn: func(data *test) {
				data.file.AuthorID = data.id
			},
			file: model.FileDTO{
				Name:        "some",
				Description: "some",
				Size:        1,
				Path:        "some",
				AddDate:     time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
				UpdateDate:  time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
				Actual:      false,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			purchaseID := tc.id
			if tc.fn != nil {
				tc.fn(&tc)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
			if tc.isOk {
				_, err = repo.Create(ctx, tc.file)
			}
			id, err := repo.DeleteByAuthorID(ctx, purchaseID)
			assert.Equal(tc.expErr, err)
			if tc.isOk {
				assert.Equal(purchaseID, id)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestFileRepo_FindByID(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2FileMongo()
	require.NoError(t, err)
	type test struct {
		name   string
		isOk   bool
		id     string
		fn     func(data *test)
		file   model.FileDTO
		exp    *model.FileDTO
		expErr error
	}
	tt := []test{
		{
			name:   "not correct userID",
			expErr: errors.New("the provided hex string is not a valid ObjectID"),
		},
		{
			name:   "not found",
			id:     primitive.NewObjectID().Hex(),
			expErr: errors.New("mongo: no documents in result"),
		},
		{
			name: "all ok",
			isOk: true,
			fn: func(data *test) {
				data.exp = &data.file
			},
			file: model.FileDTO{
				Name:        "some",
				Description: "some",
				Size:        1,
				Path:        "some",
				AddDate:     time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
				UpdateDate:  time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
				AuthorID:    primitive.NewObjectID().Hex(),
				Actual:      false,
			},
			exp: &model.FileDTO{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			fileID := tc.id
			if tc.fn != nil {
				tc.fn(&tc)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
			if tc.isOk {
				fileID, err = repo.Create(ctx, tc.file)
			}
			file, err := repo.FindByID(ctx, fileID)
			assert.Equal(tc.expErr, err)
			if tc.isOk {
				tc.exp.ID = file.ID
				assert.Equal(tc.exp, file)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestFileRepo_FindByName(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2FileMongo()
	require.NoError(t, err)
	type test struct {
		name     string
		isOk     bool
		textName string
		fn       func(data *test)
		files    []model.FileDTO
		exp      []model.FileDTO
		expErr   error
	}
	tt := []test{
		{
			name:     "not found",
			textName: "some",
		},
		{
			name: "all ok",
			isOk: true,
			fn: func(data *test) {
				data.exp = data.files
			},
			textName: "some",
			files: []model.FileDTO{
				{
					Name:        "some",
					Description: "some",
					Size:        1,
					Path:        "some",
					AddDate:     time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
					UpdateDate:  time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
					AuthorID:    primitive.NewObjectID().Hex(),
					Actual:      false,
				},
			},
			exp: []model.FileDTO{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			purchaseID := tc.textName
			if tc.fn != nil {
				tc.fn(&tc)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)

			if tc.isOk {
				for _, c := range tc.files {
					_, err = repo.Create(ctx, c)
				}
			}
			files, err := repo.FindByName(ctx, purchaseID)
			assert.Equal(tc.expErr, err)

			if tc.isOk {
				for i := range tc.exp {
					tc.exp[i].ID = files[i].ID
				}
				assert.Equal(tc.exp, files)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestFileRepo_FindAll(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2FileMongo()
	require.NoError(t, err)
	type test struct {
		name   string
		isOk   bool
		fn     func(data *test)
		files  []model.FileDTO
		exp    []model.FileDTO
		expErr error
	}
	tt := []test{
		{
			name: "not found",
		},
		{
			name: "all ok",
			isOk: true,
			fn: func(data *test) {
				data.exp = data.files
			},
			files: []model.FileDTO{
				{
					Name:        "some",
					Description: "some",
					Size:        1,
					Path:        "some",
					AddDate:     time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
					UpdateDate:  time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
					AuthorID:    primitive.NewObjectID().Hex(),
					Actual:      false,
				},
			},
			exp: []model.FileDTO{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.fn != nil {
				tc.fn(&tc)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)

			if tc.isOk {
				for _, c := range tc.files {
					_, err = repo.Create(ctx, c)
				}
			}
			files, err := repo.FindAll(ctx)
			assert.Equal(tc.expErr, err)

			if tc.isOk {
				for i := range tc.exp {
					tc.exp[i].ID = files[i].ID
				}
				assert.Equal(tc.exp, files)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestFileRepo_FindByAuthorID(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2FileMongo()
	require.NoError(t, err)
	type test struct {
		name   string
		isOk   bool
		id     string
		fn     func(data *test)
		files  []model.FileDTO
		exp    []model.FileDTO
		expErr error
	}
	tt := []test{
		{
			name:   "not correct userID",
			expErr: errors.New("the provided hex string is not a valid ObjectID"),
		},
		{
			name: "not found",
			id:   primitive.NewObjectID().Hex(),
		},
		{
			name: "all ok",
			isOk: true,
			fn: func(data *test) {
				for i := range data.files {
					data.files[i].AuthorID = data.id
				}
				data.exp = data.files
			},
			id: primitive.NewObjectID().Hex(),
			files: []model.FileDTO{
				{
					Name:        "some",
					Description: "some",
					Size:        1,
					Path:        "some",
					AddDate:     time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
					UpdateDate:  time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
					Actual:      false,
				},
			},
			exp: []model.FileDTO{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			authorID := tc.id
			if tc.fn != nil {
				tc.fn(&tc)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)

			if tc.isOk {
				for _, c := range tc.files {
					_, err = repo.Create(ctx, c)
				}
			}
			files, err := repo.FindByAuthorID(ctx, authorID)
			assert.Equal(tc.expErr, err)

			if tc.isOk {
				for i := range tc.exp {
					tc.exp[i].ID = files[i].ID
				}
				assert.Equal(tc.exp, files)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestFileRepo_FindActual(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2FileMongo()
	require.NoError(t, err)
	type test struct {
		name   string
		isOk   bool
		fn     func(data *test)
		files  []model.FileDTO
		exp    []model.FileDTO
		expErr error
	}
	tt := []test{
		{
			name: "not found",
		},
		{
			name: "all ok",
			isOk: true,
			fn: func(data *test) {
				data.exp = data.files
			},
			files: []model.FileDTO{
				{
					Name:        "some",
					Description: "some",
					Size:        1,
					Path:        "some",
					AddDate:     time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
					UpdateDate:  time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
					AuthorID:    primitive.NewObjectID().Hex(),
					Actual:      true,
				},
			},
			exp: []model.FileDTO{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.fn != nil {
				tc.fn(&tc)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)

			if tc.isOk {
				for _, c := range tc.files {
					_, err = repo.Create(ctx, c)
				}
			}
			files, err := repo.FindActual(ctx)
			assert.Equal(tc.expErr, err)

			if tc.isOk {
				for i := range tc.exp {
					tc.exp[i].ID = files[i].ID
				}
				assert.Equal(tc.exp, files)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestFileRepo_FindNotActual(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2FileMongo()
	require.NoError(t, err)
	type test struct {
		name   string
		isOk   bool
		fn     func(data *test)
		files  []model.FileDTO
		exp    []model.FileDTO
		expErr error
	}
	tt := []test{
		{
			name: "not found",
		},
		{
			name: "all ok",
			isOk: true,
			fn: func(data *test) {
				data.exp = data.files
			},
			files: []model.FileDTO{
				{
					Name:        "some",
					Description: "some",
					Size:        1,
					Path:        "some",
					AddDate:     time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
					UpdateDate:  time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
					AuthorID:    primitive.NewObjectID().Hex(),
					Actual:      false,
				},
			},
			exp: []model.FileDTO{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.fn != nil {
				tc.fn(&tc)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)

			if tc.isOk {
				for _, c := range tc.files {
					_, err = repo.Create(ctx, c)
				}
			}
			files, err := repo.FindNotActual(ctx)
			assert.Equal(tc.expErr, err)

			if tc.isOk {
				for i := range tc.exp {
					tc.exp[i].ID = files[i].ID
				}
				assert.Equal(tc.exp, files)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestFileRepo_FindAddedByPeriod(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2FileMongo()
	require.NoError(t, err)
	type test struct {
		name   string
		isOk   bool
		start  time.Time
		end    time.Time
		fn     func(data *test)
		files  []model.FileDTO
		exp    []model.FileDTO
		expErr error
	}
	tt := []test{
		{
			name:  "not found",
			start: time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
			end:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
		},
		{
			name: "all ok",
			isOk: true,
			fn: func(data *test) {
				data.exp = data.files
			},
			start: time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
			end:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
			files: []model.FileDTO{
				{
					Name:        "some",
					Description: "some",
					Size:        1,
					Path:        "some",
					AddDate:     time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
					UpdateDate:  time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
					AuthorID:    primitive.NewObjectID().Hex(),
					Actual:      false,
				},
			},
			exp: []model.FileDTO{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.fn != nil {
				tc.fn(&tc)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)

			if tc.isOk {
				for _, c := range tc.files {
					_, err = repo.Create(ctx, c)
				}
			}
			files, err := repo.FindAddedByPeriod(ctx, tc.start, tc.end)
			assert.Equal(tc.expErr, err)

			if tc.isOk {
				for i := range tc.exp {
					tc.exp[i].ID = files[i].ID
				}
				assert.Equal(tc.exp, files)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestFileRepo_FindUpdatedByPeriod(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2FileMongo()
	require.NoError(t, err)
	type test struct {
		name   string
		isOk   bool
		start  time.Time
		end    time.Time
		fn     func(data *test)
		files  []model.FileDTO
		exp    []model.FileDTO
		expErr error
	}
	tt := []test{
		{
			name:  "not found",
			start: time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
			end:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
		},
		{
			name: "all ok",
			isOk: true,
			fn: func(data *test) {
				data.exp = data.files
			},
			start: time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
			end:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
			files: []model.FileDTO{
				{
					Name:        "some",
					Description: "some",
					Size:        1,
					Path:        "some",
					AddDate:     time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
					UpdateDate:  time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
					AuthorID:    primitive.NewObjectID().Hex(),
					Actual:      false,
				},
			},
			exp: []model.FileDTO{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.fn != nil {
				tc.fn(&tc)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)

			if tc.isOk {
				for _, c := range tc.files {
					_, err = repo.Create(ctx, c)
				}
			}
			files, err := repo.FindUpdatedByPeriod(ctx, tc.start, tc.end)
			assert.Equal(tc.expErr, err)

			if tc.isOk {
				for i := range tc.exp {
					tc.exp[i].ID = files[i].ID
				}
				assert.Equal(tc.exp, files)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}
