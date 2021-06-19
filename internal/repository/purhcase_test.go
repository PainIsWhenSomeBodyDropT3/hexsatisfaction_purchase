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

func Connect2PurchaseMongo() (context.Context, *PurchaseRepo, error) {
	ctx := context.Background()
	cfg, err := config.Init()
	if err != nil {
		return nil, nil, err
	}

	db, err := mongo.NewMongo(ctx, cfg.Mongo)
	if err != nil {
		return nil, nil, err
	}

	return ctx, NewPurchaseRepo(db), nil
}

func TestPurchaseRepo_Create(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2PurchaseMongo()
	require.NoError(t, err)
	type test struct {
		name     string
		purchase model.PurchaseDTO
	}
	tt := []test{
		{
			name: "all ok",
			purchase: model.PurchaseDTO{
				UserID: 1,
				Date:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
				FileID: primitive.NewObjectID().Hex(),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var fileID string
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
			fileID, err = repo.Create(ctx, tc.purchase)

			assert.NotEmpty(fileID)
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestPurchaseRepo_Delete(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2PurchaseMongo()
	require.NoError(t, err)
	type test struct {
		name     string
		isOk     bool
		id       string
		purchase model.PurchaseDTO
		expErr   error
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
			purchase: model.PurchaseDTO{
				UserID: 1,
				Date:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
				FileID: primitive.NewObjectID().Hex(),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			purchaseID := tc.id
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
			if tc.isOk {
				purchaseID, err = repo.Create(ctx, tc.purchase)
			}
			id, err := repo.Delete(ctx, purchaseID)
			assert.Equal(tc.expErr, err)
			if tc.isOk {
				assert.Equal(purchaseID, id)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestPurchaseRepo_DeleteByFileID(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2PurchaseMongo()
	require.NoError(t, err)
	type test struct {
		name     string
		isOk     bool
		id       string
		fn       func(data *test)
		purchase model.PurchaseDTO
		expErr   error
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
				data.purchase.FileID = data.id
			},
			purchase: model.PurchaseDTO{
				UserID: 1,
				Date:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
			},
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
				_, err = repo.Create(ctx, tc.purchase)
			}
			id, err := repo.DeleteByFileID(ctx, fileID)
			assert.Equal(tc.expErr, err)
			if tc.isOk {
				assert.Equal(fileID, id)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestPurchaseRepo_FindByID(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2PurchaseMongo()
	require.NoError(t, err)
	type test struct {
		name     string
		isOk     bool
		id       string
		fn       func(data *test)
		purchase model.PurchaseDTO
		exp      *model.PurchaseDTO
		expErr   error
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
				data.exp = &data.purchase
			},
			purchase: model.PurchaseDTO{
				UserID: 1,
				Date:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
				FileID: primitive.NewObjectID().Hex(),
			},
			exp: &model.PurchaseDTO{},
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
				fileID, err = repo.Create(ctx, tc.purchase)
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

func TestPurchaseRepo_FindLastByUserID(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2PurchaseMongo()
	require.NoError(t, err)
	type test struct {
		name     string
		isOk     bool
		id       int
		fn       func(data *test)
		purchase model.PurchaseDTO
		exp      *model.PurchaseDTO
		expErr   error
	}
	tt := []test{
		{
			name:   "not found",
			id:     1,
			expErr: errors.New("mongo: no documents in result"),
		},
		{
			name: "all ok",
			isOk: true,
			fn: func(data *test) {
				data.purchase.UserID = data.id
				data.exp = &data.purchase
			},
			id: 1,
			purchase: model.PurchaseDTO{
				Date:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
				FileID: primitive.NewObjectID().Hex(),
			},
			exp: &model.PurchaseDTO{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			userID := tc.id
			if tc.fn != nil {
				tc.fn(&tc)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
			if tc.isOk {
				_, err = repo.Create(ctx, tc.purchase)
			}
			file, err := repo.FindLastByUserID(ctx, userID)
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

func TestPurchaseRepo_FindAllByUserID(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2PurchaseMongo()
	require.NoError(t, err)
	type test struct {
		name      string
		isOk      bool
		id        int
		fn        func(data *test)
		purchases []model.PurchaseDTO
		exp       []model.PurchaseDTO
		expErr    error
	}
	tt := []test{
		{
			name: "not found",
			id:   1,
		},
		{
			name: "all ok",
			isOk: true,
			fn: func(data *test) {
				for i := range data.purchases {
					data.purchases[i].UserID = data.id
				}
				data.exp = data.purchases
			},
			id: 1,
			purchases: []model.PurchaseDTO{
				{
					Date:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
					FileID: primitive.NewObjectID().Hex(),
				},
			},
			exp: []model.PurchaseDTO{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			userID := tc.id
			if tc.fn != nil {
				tc.fn(&tc)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)

			if tc.isOk {
				for _, c := range tc.purchases {
					_, err = repo.Create(ctx, c)
				}
			}
			files, err := repo.FindAllByUserID(ctx, userID)
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

func TestPurchaseRepo_FindByUserIDAndPeriod(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2PurchaseMongo()
	require.NoError(t, err)
	type test struct {
		name      string
		isOk      bool
		id        int
		start     time.Time
		end       time.Time
		fn        func(data *test)
		purchases []model.PurchaseDTO
		exp       []model.PurchaseDTO
		expErr    error
	}
	tt := []test{
		{
			name: "not found",
			id:   1,
		},
		{
			name: "all ok",
			isOk: true,
			id:   1,
			fn: func(data *test) {
				for i := range data.purchases {
					data.purchases[i].UserID = data.id
				}
				data.exp = data.purchases
			},

			start: time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
			end:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
			purchases: []model.PurchaseDTO{
				{
					Date:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
					FileID: primitive.NewObjectID().Hex(),
				},
			},
			exp: []model.PurchaseDTO{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			userID := tc.id
			if tc.fn != nil {
				tc.fn(&tc)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)

			if tc.isOk {
				for _, c := range tc.purchases {
					_, err = repo.Create(ctx, c)
				}
			}
			files, err := repo.FindByUserIDAndPeriod(ctx, userID, tc.start, tc.end)
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

func TestPurchaseRepo_FindByUserIDAfterDate(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2PurchaseMongo()
	require.NoError(t, err)
	type test struct {
		name      string
		isOk      bool
		id        int
		start     time.Time
		fn        func(data *test)
		purchases []model.PurchaseDTO
		exp       []model.PurchaseDTO
		expErr    error
	}
	tt := []test{
		{
			name: "not found",
			id:   1,
		},
		{
			name: "all ok",
			isOk: true,
			id:   1,
			fn: func(data *test) {
				for i := range data.purchases {
					data.purchases[i].UserID = data.id
				}
				data.exp = data.purchases
			},

			start: time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
			purchases: []model.PurchaseDTO{
				{
					Date:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
					FileID: primitive.NewObjectID().Hex(),
				},
			},
			exp: []model.PurchaseDTO{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			userID := tc.id
			if tc.fn != nil {
				tc.fn(&tc)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)

			if tc.isOk {
				for _, c := range tc.purchases {
					_, err = repo.Create(ctx, c)
				}
			}
			files, err := repo.FindByUserIDAfterDate(ctx, userID, tc.start)
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

func TestPurchaseRepo_FindByUserIDBeforeDate(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2PurchaseMongo()
	require.NoError(t, err)
	type test struct {
		name      string
		isOk      bool
		id        int
		end       time.Time
		fn        func(data *test)
		purchases []model.PurchaseDTO
		exp       []model.PurchaseDTO
		expErr    error
	}
	tt := []test{
		{
			name: "not found",
			id:   1,
		},
		{
			name: "all ok",
			isOk: true,
			id:   1,
			fn: func(data *test) {
				for i := range data.purchases {
					data.purchases[i].UserID = data.id
				}
				data.exp = data.purchases
			},

			end: time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
			purchases: []model.PurchaseDTO{
				{
					Date:   time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
					FileID: primitive.NewObjectID().Hex(),
				},
			},
			exp: []model.PurchaseDTO{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			userID := tc.id
			if tc.fn != nil {
				tc.fn(&tc)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)

			if tc.isOk {
				for _, c := range tc.purchases {
					_, err = repo.Create(ctx, c)
				}
			}
			files, err := repo.FindByUserIDBeforeDate(ctx, userID, tc.end)
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

func TestPurchaseRepo_FindByUserIDAndFileID(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2PurchaseMongo()
	require.NoError(t, err)
	type test struct {
		name      string
		isOk      bool
		userID    int
		fileID    string
		fn        func(data *test)
		purchases []model.PurchaseDTO
		exp       []model.PurchaseDTO
		expErr    error
	}
	tt := []test{
		{
			name:   "not found",
			userID: 1,
			fileID: primitive.NewObjectID().Hex(),
		},
		{
			name:   "all ok",
			isOk:   true,
			userID: 1,
			fileID: primitive.NewObjectID().Hex(),
			fn: func(data *test) {
				for i := range data.purchases {
					data.purchases[i].UserID = data.userID
					data.purchases[i].FileID = data.fileID
				}
				data.exp = data.purchases
			},
			purchases: []model.PurchaseDTO{
				{
					Date: time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
				},
			},
			exp: []model.PurchaseDTO{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			userID := tc.userID
			fileID := tc.fileID
			if tc.fn != nil {
				tc.fn(&tc)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)

			if tc.isOk {
				for _, c := range tc.purchases {
					_, err = repo.Create(ctx, c)
				}
			}
			files, err := repo.FindByUserIDAndFileID(ctx, userID, fileID)
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

func TestPurchaseRepo_FindLast(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2PurchaseMongo()
	require.NoError(t, err)
	type test struct {
		name     string
		isOk     bool
		fn       func(data *test)
		purchase model.PurchaseDTO
		exp      *model.PurchaseDTO
		expErr   error
	}
	tt := []test{
		{
			name:   "not found",
			expErr: errors.New("mongo: no documents in result"),
		},
		{
			name: "all ok",
			isOk: true,
			fn: func(data *test) {
				data.exp = &data.purchase
			},
			purchase: model.PurchaseDTO{
				UserID: 1,
				Date:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
				FileID: primitive.NewObjectID().Hex(),
			},
			exp: &model.PurchaseDTO{},
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
				_, err = repo.Create(ctx, tc.purchase)
			}
			file, err := repo.FindLast(ctx)
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

func TestPurchaseRepo_FindAll(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2PurchaseMongo()
	require.NoError(t, err)
	type test struct {
		name      string
		isOk      bool
		fn        func(data *test)
		purchases []model.PurchaseDTO
		exp       []model.PurchaseDTO
		expErr    error
	}
	tt := []test{
		{
			name: "not found",
		},
		{
			name: "all ok",
			isOk: true,
			fn: func(data *test) {
				data.exp = data.purchases
			},
			purchases: []model.PurchaseDTO{
				{
					UserID: 1,
					Date:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
					FileID: primitive.NewObjectID().Hex(),
				},
			},
			exp: []model.PurchaseDTO{},
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
				for _, c := range tc.purchases {
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

func TestPurchaseRepo_FindByPeriod(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2PurchaseMongo()
	require.NoError(t, err)
	type test struct {
		name      string
		isOk      bool
		start     time.Time
		end       time.Time
		fn        func(data *test)
		purchases []model.PurchaseDTO
		exp       []model.PurchaseDTO
		expErr    error
	}
	tt := []test{
		{
			name: "not found",
		},
		{
			name: "all ok",
			isOk: true,
			fn: func(data *test) {
				data.exp = data.purchases
			},

			start: time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
			end:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
			purchases: []model.PurchaseDTO{
				{
					UserID: 1,
					Date:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
					FileID: primitive.NewObjectID().Hex(),
				},
			},
			exp: []model.PurchaseDTO{},
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
				for _, c := range tc.purchases {
					_, err = repo.Create(ctx, c)
				}
			}
			files, err := repo.FindByPeriod(ctx, tc.start, tc.end)
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

func TestPurchaseRepo_FindAfterDate(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2PurchaseMongo()
	require.NoError(t, err)
	type test struct {
		name      string
		isOk      bool
		start     time.Time
		fn        func(data *test)
		purchases []model.PurchaseDTO
		exp       []model.PurchaseDTO
		expErr    error
	}
	tt := []test{
		{
			name: "not found",
		},
		{
			name: "all ok",
			isOk: true,
			fn: func(data *test) {
				data.exp = data.purchases
			},

			start: time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
			purchases: []model.PurchaseDTO{
				{
					UserID: 1,
					Date:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
					FileID: primitive.NewObjectID().Hex(),
				},
			},
			exp: []model.PurchaseDTO{},
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
				for _, c := range tc.purchases {
					_, err = repo.Create(ctx, c)
				}
			}
			files, err := repo.FindAfterDate(ctx, tc.start)
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

func TestPurchaseRepo_FindBeforeDate(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2PurchaseMongo()
	require.NoError(t, err)
	type test struct {
		name      string
		isOk      bool
		end       time.Time
		fn        func(data *test)
		purchases []model.PurchaseDTO
		exp       []model.PurchaseDTO
		expErr    error
	}
	tt := []test{
		{
			name: "not found",
		},
		{
			name: "all ok",
			isOk: true,
			fn: func(data *test) {
				data.exp = data.purchases
			},

			end: time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
			purchases: []model.PurchaseDTO{
				{
					UserID: 1,
					Date:   time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
					FileID: primitive.NewObjectID().Hex(),
				},
			},
			exp: []model.PurchaseDTO{},
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
				for _, c := range tc.purchases {
					_, err = repo.Create(ctx, c)
				}
			}
			files, err := repo.FindBeforeDate(ctx, tc.end)
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

func TestPurchaseRepo_FindByFileID(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2PurchaseMongo()
	require.NoError(t, err)
	type test struct {
		name      string
		isOk      bool
		fileID    string
		fn        func(data *test)
		purchases []model.PurchaseDTO
		exp       []model.PurchaseDTO
		expErr    error
	}
	tt := []test{
		{
			name:   "not correct fileID",
			expErr: errors.New("the provided hex string is not a valid ObjectID"),
		},
		{
			name:   "not found",
			fileID: primitive.NewObjectID().Hex(),
		},
		{
			name:   "all ok",
			isOk:   true,
			fileID: primitive.NewObjectID().Hex(),
			fn: func(data *test) {
				for i := range data.purchases {
					data.purchases[i].FileID = data.fileID
				}
				data.exp = data.purchases
			},
			purchases: []model.PurchaseDTO{
				{
					UserID: 1,
					Date:   time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
				},
			},
			exp: []model.PurchaseDTO{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			fileID := tc.fileID
			if tc.fn != nil {
				tc.fn(&tc)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)

			if tc.isOk {
				for _, c := range tc.purchases {
					_, err = repo.Create(ctx, c)
				}
			}
			files, err := repo.FindByFileID(ctx, fileID)
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
