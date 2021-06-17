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

func Connect2CommentMongo() (context.Context, *CommentRepo, error) {
	ctx := context.Background()
	cfg, err := config.Init()
	if err != nil {
		return nil, nil, err
	}

	db, err := mongo.NewMongo(ctx, cfg.Mongo)
	if err != nil {
		return nil, nil, err
	}

	return ctx, NewCommentRepo(db), nil
}
func TestCommentRepo_Create(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2CommentMongo()
	require.NoError(t, err)
	type test struct {
		name    string
		comment model.CommentDTO
	}
	tt := []test{
		{
			name: "all ok",
			comment: model.CommentDTO{
				UserID:     1,
				PurchaseID: primitive.NewObjectID().Hex(),
				Date:       time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
				Text:       "some",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var commentID string
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
			commentID, err = repo.Create(ctx, tc.comment)

			assert.NotEmpty(commentID)
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestCommentRepo_Update(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2CommentMongo()
	require.NoError(t, err)
	type test struct {
		name    string
		isOk    bool
		id      string
		comment model.CommentDTO
		expErr  error
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
			comment: model.CommentDTO{
				UserID:     1,
				PurchaseID: primitive.NewObjectID().Hex(),
				Date:       time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
				Text:       "some",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			commentID := tc.id
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
			if tc.isOk {
				commentID, err = repo.Create(ctx, tc.comment)
			}
			id, err := repo.Update(ctx, commentID, tc.comment)
			assert.Equal(tc.expErr, err)
			if tc.isOk {
				assert.Equal(commentID, id)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestCommentRepo_Delete(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2CommentMongo()
	require.NoError(t, err)
	type test struct {
		name    string
		isOk    bool
		id      string
		comment model.CommentDTO
		expErr  error
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
			comment: model.CommentDTO{
				UserID:     1,
				PurchaseID: primitive.NewObjectID().Hex(),
				Date:       time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
				Text:       "some",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			commentID := tc.id
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
			if tc.isOk {
				commentID, err = repo.Create(ctx, tc.comment)
			}
			id, err := repo.Delete(ctx, commentID)
			assert.Equal(tc.expErr, err)
			if tc.isOk {
				assert.Equal(commentID, id)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestCommentRepo_DeleteByPurchaseID(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2CommentMongo()
	require.NoError(t, err)
	type test struct {
		name    string
		isOk    bool
		id      string
		fn      func(data *test)
		comment model.CommentDTO
		expErr  error
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
				data.comment.PurchaseID = data.id
			},
			comment: model.CommentDTO{
				UserID: 1,
				Date:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
				Text:   "some",
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
				_, err = repo.Create(ctx, tc.comment)
			}
			id, err := repo.DeleteByPurchaseID(ctx, purchaseID)
			assert.Equal(tc.expErr, err)
			if tc.isOk {
				assert.Equal(purchaseID, id)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestCommentRepo_FindByID(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2CommentMongo()
	require.NoError(t, err)
	type test struct {
		name    string
		isOk    bool
		id      string
		fn      func(data *test)
		comment model.CommentDTO
		exp     *model.CommentDTO
		expErr  error
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
				data.exp = &data.comment
			},
			comment: model.CommentDTO{
				UserID:     1,
				PurchaseID: primitive.NewObjectID().Hex(),
				Date:       time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
				Text:       "some",
			},
			exp: &model.CommentDTO{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			commentID := tc.id
			if tc.fn != nil {
				tc.fn(&tc)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
			if tc.isOk {
				commentID, err = repo.Create(ctx, tc.comment)
			}
			comment, err := repo.FindByID(ctx, commentID)
			assert.Equal(tc.expErr, err)
			if tc.isOk {
				tc.exp.ID = comment.ID
				assert.Equal(tc.exp, comment)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestCommentRepo_FindByPurchaseID(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2CommentMongo()
	require.NoError(t, err)
	type test struct {
		name     string
		isOk     bool
		id       string
		fn       func(data *test)
		comments []model.CommentDTO
		exp      []model.CommentDTO
		expErr   error
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
				for i := range data.comments {
					data.comments[i].PurchaseID = data.id
				}
				data.exp = data.comments
			},
			id: primitive.NewObjectID().Hex(),
			comments: []model.CommentDTO{
				{
					UserID: 1,
					Date:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
					Text:   "some",
				},
			},
			exp: []model.CommentDTO{},
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
				for _, c := range tc.comments {
					_, err = repo.Create(ctx, c)
				}
			}
			comments, err := repo.FindByPurchaseID(ctx, purchaseID)
			assert.Equal(tc.expErr, err)

			if tc.isOk {
				for i := range tc.exp {
					tc.exp[i].ID = comments[i].ID
				}
				assert.Equal(tc.exp, comments)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestCommentRepo_FindAllByUserID(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2CommentMongo()
	require.NoError(t, err)
	type test struct {
		name     string
		isOk     bool
		id       int
		fn       func(data *test)
		comments []model.CommentDTO
		exp      []model.CommentDTO
		expErr   error
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
				for i := range data.comments {
					data.comments[i].UserID = data.id
				}
				data.exp = data.comments
			},
			id: 1,
			comments: []model.CommentDTO{
				{
					UserID:     1,
					PurchaseID: primitive.NewObjectID().Hex(),
					Date:       time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
					Text:       "some",
				},
			},
			exp: []model.CommentDTO{},
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
				for _, c := range tc.comments {
					_, err = repo.Create(ctx, c)
				}
			}
			comments, err := repo.FindAllByUserID(ctx, userID)
			assert.Equal(tc.expErr, err)

			if tc.isOk {
				for i := range tc.exp {
					tc.exp[i].ID = comments[i].ID
				}
				assert.Equal(tc.exp, comments)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestCommentRepo_FindAllByUserIDAndPurchaseID(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2CommentMongo()
	require.NoError(t, err)
	type test struct {
		name       string
		isOk       bool
		userID     int
		purchaseID string
		fn         func(data *test)
		comments   []model.CommentDTO
		exp        []model.CommentDTO
		expErr     error
	}
	tt := []test{
		{
			name:   "not correct userID",
			expErr: errors.New("the provided hex string is not a valid ObjectID"),
		},
		{
			name:       "not found",
			userID:     1,
			purchaseID: primitive.NewObjectID().Hex(),
		},
		{
			name: "all ok",
			isOk: true,
			fn: func(data *test) {
				for i := range data.comments {
					data.comments[i].UserID = data.userID
					data.comments[i].PurchaseID = data.purchaseID
				}
				data.exp = data.comments
			},
			userID:     1,
			purchaseID: primitive.NewObjectID().Hex(),
			comments: []model.CommentDTO{
				{
					Date: time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
					Text: "some",
				},
			},
			exp: []model.CommentDTO{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			userID := tc.userID
			purchaseID := tc.purchaseID
			if tc.fn != nil {
				tc.fn(&tc)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)

			if tc.isOk {
				for _, c := range tc.comments {
					_, err = repo.Create(ctx, c)
				}
			}
			comments, err := repo.FindByUserIDAndPurchaseID(ctx, userID, purchaseID)
			assert.Equal(tc.expErr, err)

			if tc.isOk {
				for i := range tc.exp {
					tc.exp[i].ID = comments[i].ID
				}
				assert.Equal(tc.exp, comments)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestCommentRepo_FindAll(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2CommentMongo()
	require.NoError(t, err)
	type test struct {
		name     string
		isOk     bool
		id       string
		fn       func(data *test)
		comments []model.CommentDTO
		exp      []model.CommentDTO
		expErr   error
	}
	tt := []test{
		{
			name: "not found",
		},
		{
			name: "all ok",
			isOk: true,
			fn: func(data *test) {
				data.exp = data.comments
			},
			comments: []model.CommentDTO{
				{
					UserID:     1,
					PurchaseID: primitive.NewObjectID().Hex(),
					Date:       time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
					Text:       "some",
				},
			},
			exp: []model.CommentDTO{},
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
				for _, c := range tc.comments {
					_, err = repo.Create(ctx, c)
				}
			}
			comments, err := repo.FindAll(ctx)
			assert.Equal(tc.expErr, err)

			if tc.isOk {
				for i := range tc.exp {
					tc.exp[i].ID = comments[i].ID
				}
				assert.Equal(tc.exp, comments)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestCommentRepo_FindByText(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2CommentMongo()
	require.NoError(t, err)
	type test struct {
		name     string
		isOk     bool
		text     string
		fn       func(data *test)
		comments []model.CommentDTO
		exp      []model.CommentDTO
		expErr   error
	}
	tt := []test{
		{
			name: "not found",
			text: "not correct",
		},
		{
			name: "all ok",
			isOk: true,
			fn: func(data *test) {
				data.exp = data.comments
			},
			text: "some",
			comments: []model.CommentDTO{
				{
					UserID:     1,
					PurchaseID: primitive.NewObjectID().Hex(),
					Date:       time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
					Text:       "some thing",
				},
			},
			exp: []model.CommentDTO{},
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
				for _, c := range tc.comments {
					_, err = repo.Create(ctx, c)
				}
			}
			comments, err := repo.FindByText(ctx, tc.text)
			assert.Equal(tc.expErr, err)

			if tc.isOk {
				for i := range tc.exp {
					tc.exp[i].ID = comments[i].ID
				}
				assert.Equal(tc.exp, comments)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}

func TestCommentRepo_FindAllByPeriod(t *testing.T) {
	assert := assertTest.New(t)
	ctx, repo, err := Connect2CommentMongo()
	require.NoError(t, err)
	type test struct {
		name     string
		isOk     bool
		start    time.Time
		end      time.Time
		fn       func(data *test)
		comments []model.CommentDTO
		exp      []model.CommentDTO
		expErr   error
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
				data.exp = data.comments
			},
			start: time.Date(2020, time.November, 10, 23, 10, 34, 0, time.UTC),
			end:   time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
			comments: []model.CommentDTO{
				{
					UserID:     1,
					PurchaseID: primitive.NewObjectID().Hex(),
					Date:       time.Date(2020, time.December, 10, 23, 10, 34, 0, time.UTC),
					Text:       "some",
				},
			},
			exp: []model.CommentDTO{},
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
				for _, c := range tc.comments {
					_, err = repo.Create(ctx, c)
				}
			}
			comments, err := repo.FindByPeriod(ctx, tc.start, tc.end)
			assert.Equal(tc.expErr, err)

			if tc.isOk {
				for i := range tc.exp {
					tc.exp[i].ID = comments[i].ID
				}
				assert.Equal(tc.exp, comments)
			}
			_, err = repo.collection.DeleteMany(ctx, bson.M{})
			assert.NoError(err)
		})
	}
}
