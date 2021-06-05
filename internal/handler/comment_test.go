package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	m "github.com/JesusG2000/hexsatisfaction_purchase/internal/handler/mock"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/model"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/service"
	"github.com/pkg/errors"
	testAssert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	comment             = "comment"
	text                = "text"
	slash               = "/"
	api                 = "api"
	user                = "user"
	purchase            = "purchase"
	period              = "period"
	authorizationHeader = "Authorization"
)

func TestComment_Create(t *testing.T) {
	assert := testAssert.New(t)
	id := primitive.NewObjectID().Hex()
	testAPI, err := service.InitTest4Mock()
	require.NoError(t, err)
	token, err := testAPI.TokenManager.NewJWT(mock.Anything)
	require.NoError(t, err)

	type test struct {
		name    string
		path    string
		method  string
		req     model.CreateCommentRequest
		fn      func(commentService *m.Comment, data test)
		expCode int
		expBody string
	}

	tt := []test{
		{
			name:   "invalid user id",
			path:   slash + comment + slash + api + slash,
			method: http.MethodPost,
			req: model.CreateCommentRequest{
				Date: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				Text: "some text",
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("Create", mock.Anything, data.req).
					Return("", nil)
			},
			expCode: http.StatusBadRequest,
			expBody: "not correct user id",
		},
		{
			name:   "create err",
			path:   slash + comment + slash + api + slash,
			method: http.MethodPost,
			req: model.CreateCommentRequest{
				UserID:     id,
				PurchaseID: id,
				Date:       time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				Text:       "some text",
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("Create", mock.Anything, data.req).
					Return("", errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "all ok",
			path:   slash + comment + slash + api + slash,
			method: http.MethodPost,
			req: model.CreateCommentRequest{
				UserID:     id,
				PurchaseID: id,
				Date:       time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				Text:       "some text",
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("Create", mock.Anything, data.req).
					Return(id, nil)
			},
			expCode: http.StatusOK,
			expBody: id,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			comment := new(m.Comment)
			testAPI.Services.Comment = comment
			router := newComment(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(comment, tc)
			}

			body := new(bytes.Buffer)
			err := json.NewEncoder(body).Encode(&tc.req)
			assert.Nil(err)

			req, err := http.NewRequest(tc.method, tc.path, body)
			assert.Nil(err)

			req.Header.Set(authorizationHeader, "Bearer "+token)

			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			assert.Equal(tc.expCode, res.Code)

			err = json.NewDecoder(res.Body).Decode(&r)
			assert.Nil(err)
			assert.Equal(tc.expBody, r)
		})
	}
}

func TestComment_Update(t *testing.T) {
	assert := testAssert.New(t)
	id := primitive.NewObjectID().Hex()
	testAPI, err := service.InitTest4Mock()
	require.NoError(t, err)
	token, err := testAPI.TokenManager.NewJWT(mock.Anything)
	require.NoError(t, err)

	type test struct {
		name    string
		path    string
		method  string
		isOkRes bool
		req     model.UpdateCommentRequest
		fn      func(commentService *m.Comment, data test)
		expCode int
		expBody string
	}

	tt := []test{
		{
			name:    "invalid id",
			path:    slash + comment + slash + api + slash,
			method:  http.MethodPut,
			isOkRes: true,
			req: model.UpdateCommentRequest{
				ID:   "some",
				Date: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				Text: "some text",
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("Update", mock.Anything, data.req).
					Return("", nil)
			},
			expCode: http.StatusBadRequest,
			expBody: "not correct id",
		},
		{
			name:    "update err",
			path:    slash + comment + slash + api + slash,
			method:  http.MethodPut,
			isOkRes: true,
			req: model.UpdateCommentRequest{
				ID:         id,
				UserID:     id,
				PurchaseID: id,
				Date:       time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				Text:       "some text",
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("Update", mock.Anything, data.req).
					Return("", errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + comment + slash + api + slash,
			method: http.MethodPut,
			req: model.UpdateCommentRequest{
				ID:         id,
				UserID:     id,
				PurchaseID: id,
				Date:       time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				Text:       "some text",
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("Update", mock.Anything, data.req).
					Return("", nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + comment + slash + api + slash,
			method:  http.MethodPut,
			isOkRes: true,
			req: model.UpdateCommentRequest{
				ID:         id,
				UserID:     id,
				PurchaseID: id,
				Date:       time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				Text:       "some text",
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("Update", mock.Anything, data.req).
					Return(data.expBody, nil)
			},
			expCode: http.StatusOK,
			expBody: id,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			comment := new(m.Comment)
			testAPI.Services.Comment = comment
			router := newComment(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(comment, tc)
			}

			body := new(bytes.Buffer)
			err := json.NewEncoder(body).Encode(&tc.req)
			assert.Nil(err)

			req, err := http.NewRequest(tc.method, tc.path+tc.req.ID, body)
			assert.Nil(err)

			req.Header.Set(authorizationHeader, "Bearer "+token)

			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			assert.Equal(tc.expCode, res.Code)

			if tc.isOkRes {
				err = json.NewDecoder(res.Body).Decode(&r)
				assert.Nil(err)
			}
			assert.Equal(tc.expBody, r)
		})
	}
}

func TestComment_Delete(t *testing.T) {
	assert := testAssert.New(t)
	id := primitive.NewObjectID().Hex()
	testAPI, err := service.InitTest4Mock()
	require.NoError(t, err)
	token, err := testAPI.TokenManager.NewJWT(mock.Anything)
	require.NoError(t, err)

	type test struct {
		name    string
		path    string
		method  string
		isOkRes bool
		req     model.DeleteCommentRequest
		fn      func(commentService *m.Comment, data test)
		expCode int
		expBody string
	}

	tt := []test{
		{
			name:    "invalid id",
			path:    slash + comment + slash + api + slash,
			method:  http.MethodDelete,
			isOkRes: true,
			req:     model.DeleteCommentRequest{ID: "some"},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("Delete", mock.Anything, data.req).
					Return("", nil)
			},
			expCode: http.StatusBadRequest,
			expBody: "not correct id",
		},
		{
			name:    "delete err",
			path:    slash + comment + slash + api + slash,
			method:  http.MethodDelete,
			isOkRes: true,
			req: model.DeleteCommentRequest{
				ID: id,
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("Delete", mock.Anything, data.req).
					Return("", errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + comment + slash + api + slash,
			method: http.MethodDelete,
			req: model.DeleteCommentRequest{
				ID: id,
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("Delete", mock.Anything, data.req).
					Return("", nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + comment + slash + api + slash,
			method:  http.MethodDelete,
			isOkRes: true,
			req: model.DeleteCommentRequest{
				ID: id,
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("Delete", mock.Anything, data.req).
					Return(data.expBody, nil)
			},
			expCode: http.StatusOK,
			expBody: id,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			comment := new(m.Comment)
			testAPI.Services.Comment = comment
			router := newComment(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(comment, tc)
			}

			req, err := http.NewRequest(tc.method, tc.path+tc.req.ID, nil)
			assert.Nil(err)

			req.Header.Set(authorizationHeader, "Bearer "+token)

			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			assert.Equal(tc.expCode, res.Code)

			if tc.isOkRes {
				err = json.NewDecoder(res.Body).Decode(&r)
				assert.Nil(err)
			}
			assert.Equal(tc.expBody, r)
		})
	}
}

func TestComment_FindByID(t *testing.T) {
	assert := testAssert.New(t)
	id := primitive.NewObjectID().Hex()
	testAPI, err := service.InitTest4Mock()
	require.NoError(t, err)
	token, err := testAPI.TokenManager.NewJWT(mock.Anything)
	require.NoError(t, err)

	type test struct {
		name        string
		path        string
		method      string
		isOkRes     bool
		isOkMessage bool
		req         model.IDCommentRequest
		fn          func(commentService *m.Comment, data test)
		expCode     int
		expRes      model.CommentDTO
		message     string
	}

	tt := []test{
		{
			name:        "invalid id",
			path:        slash + comment + slash + api + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			req:         model.IDCommentRequest{ID: "some"},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByID", mock.Anything, data.req).
					Return(&data.expRes, nil)
			},
			expCode: http.StatusBadRequest,
			message: "not correct id",
		},
		{
			name:        "find err",
			path:        slash + comment + slash + api + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			req: model.IDCommentRequest{
				ID: id,
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByID", mock.Anything, data.req).
					Return(&data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + comment + slash + api + slash,
			method: http.MethodGet,
			req: model.IDCommentRequest{
				ID: id,
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByID", mock.Anything, data.req).
					Return(&data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + comment + slash + api + slash,
			method:  http.MethodGet,
			isOkRes: true,
			req: model.IDCommentRequest{
				ID: id,
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByID", mock.Anything, data.req).
					Return(&data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: model.CommentDTO{
				ID:         id,
				UserID:     id,
				PurchaseID: id,
				Date:       time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				Text:       "some text",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var c model.CommentDTO
			comment := new(m.Comment)
			testAPI.Services.Comment = comment
			router := newComment(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(comment, tc)
			}

			req, err := http.NewRequest(tc.method, tc.path+tc.req.ID, nil)
			assert.Nil(err)

			req.Header.Set(authorizationHeader, "Bearer "+token)

			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			assert.Equal(tc.expCode, res.Code)

			switch {
			case tc.isOkMessage:
				err = json.NewDecoder(res.Body).Decode(&r)
				assert.Nil(err)
				assert.Equal(tc.message, r)
			case tc.isOkRes:
				err = json.NewDecoder(res.Body).Decode(&c)
				assert.Nil(err)
				assert.Equal(tc.expRes, c)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}

func TestComment_FindAll(t *testing.T) {
	assert := testAssert.New(t)
	id := primitive.NewObjectID().Hex()
	testAPI, err := service.InitTest4Mock()
	require.NoError(t, err)
	token, err := testAPI.TokenManager.NewJWT(mock.Anything)
	require.NoError(t, err)

	type test struct {
		name        string
		path        string
		method      string
		isOkRes     bool
		isOkMessage bool
		fn          func(commentService *m.Comment, data test)
		expCode     int
		expRes      []model.CommentDTO
		message     string
	}

	tt := []test{
		{
			name:        "find err",
			path:        slash + comment + slash + api + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindAll", mock.Anything).
					Return(data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + comment + slash + api + slash,
			method: http.MethodGet,
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindAll", mock.Anything).
					Return(data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + comment + slash + api + slash,
			method:  http.MethodGet,
			isOkRes: true,
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindAll", mock.Anything).
					Return(data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: []model.CommentDTO{
				{
					ID:         id,
					UserID:     id,
					PurchaseID: id,
					Date:       time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					Text:       "some text",
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var c []model.CommentDTO
			comment := new(m.Comment)
			testAPI.Services.Comment = comment
			router := newComment(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(comment, tc)
			}

			req, err := http.NewRequest(tc.method, tc.path, nil)
			assert.Nil(err)

			req.Header.Set(authorizationHeader, "Bearer "+token)

			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			assert.Equal(tc.expCode, res.Code)

			switch {
			case tc.isOkMessage:
				err = json.NewDecoder(res.Body).Decode(&r)
				assert.Nil(err)
				assert.Equal(tc.message, r)
			case tc.isOkRes:
				err = json.NewDecoder(res.Body).Decode(&c)
				assert.Nil(err)
				assert.Equal(tc.expRes, c)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}

func TestComment_FindAllByUserID(t *testing.T) {
	assert := testAssert.New(t)
	id := primitive.NewObjectID().Hex()
	testAPI, err := service.InitTest4Mock()
	require.NoError(t, err)

	type test struct {
		name        string
		path        string
		method      string
		isOkRes     bool
		isOkMessage bool
		req         model.UserIDCommentRequest
		fn          func(commentService *m.Comment, data test)
		expCode     int
		expRes      []model.CommentDTO
		message     string
	}

	tt := []test{
		{
			name:        "invalid id",
			path:        slash + comment + slash + user + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			req:         model.UserIDCommentRequest{ID: "some"},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindAllByUserID", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusBadRequest,
			message: "not correct id",
		},
		{
			name:        "find err",
			path:        slash + comment + slash + user + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			req: model.UserIDCommentRequest{
				ID: id,
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindAllByUserID", mock.Anything, data.req).
					Return(data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + comment + slash + user + slash,
			method: http.MethodGet,
			req: model.UserIDCommentRequest{
				ID: id,
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindAllByUserID", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + comment + slash + user + slash,
			method:  http.MethodGet,
			isOkRes: true,
			req: model.UserIDCommentRequest{
				ID: id,
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindAllByUserID", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: []model.CommentDTO{
				{
					ID:         id,
					UserID:     id,
					PurchaseID: id,
					Date:       time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					Text:       "some text",
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var c []model.CommentDTO
			comment := new(m.Comment)
			testAPI.Services.Comment = comment
			router := newComment(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(comment, tc)
			}

			req, err := http.NewRequest(tc.method, tc.path+tc.req.ID, nil)
			assert.Nil(err)

			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			assert.Equal(tc.expCode, res.Code)

			switch {
			case tc.isOkMessage:
				err = json.NewDecoder(res.Body).Decode(&r)
				assert.Nil(err)
				assert.Equal(tc.message, r)
			case tc.isOkRes:
				err = json.NewDecoder(res.Body).Decode(&c)
				assert.Nil(err)
				assert.Equal(tc.expRes, c)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}

func TestComment_FindByPurchaseID(t *testing.T) {
	assert := testAssert.New(t)
	id := primitive.NewObjectID().Hex()
	testAPI, err := service.InitTest4Mock()
	require.NoError(t, err)

	type test struct {
		name        string
		path        string
		method      string
		isOkRes     bool
		isOkMessage bool
		req         model.PurchaseIDCommentRequest
		fn          func(commentService *m.Comment, data test)
		expCode     int
		expRes      []model.CommentDTO
		message     string
	}

	tt := []test{
		{
			name:        "invalid id",
			path:        slash + comment + slash + purchase + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			req:         model.PurchaseIDCommentRequest{ID: "some"},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByPurchaseID", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusBadRequest,
			message: "not correct id",
		},
		{
			name:        "find err",
			path:        slash + comment + slash + purchase + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			req: model.PurchaseIDCommentRequest{
				ID: id,
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByPurchaseID", mock.Anything, data.req).
					Return(data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + comment + slash + purchase + slash,
			method: http.MethodGet,
			req: model.PurchaseIDCommentRequest{
				ID: id,
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByPurchaseID", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + comment + slash + purchase + slash,
			method:  http.MethodGet,
			isOkRes: true,
			req: model.PurchaseIDCommentRequest{
				ID: id,
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByPurchaseID", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: []model.CommentDTO{
				{
					ID:         id,
					UserID:     id,
					PurchaseID: id,
					Date:       time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					Text:       "some text",
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var c []model.CommentDTO
			comment := new(m.Comment)
			testAPI.Services.Comment = comment
			router := newComment(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(comment, tc)
			}

			req, err := http.NewRequest(tc.method, tc.path+tc.req.ID, nil)
			assert.Nil(err)

			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			assert.Equal(tc.expCode, res.Code)

			switch {
			case tc.isOkMessage:
				err = json.NewDecoder(res.Body).Decode(&r)
				assert.Nil(err)
				assert.Equal(tc.message, r)
			case tc.isOkRes:
				err = json.NewDecoder(res.Body).Decode(&c)
				assert.Nil(err)
				assert.Equal(tc.expRes, c)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}

func TestComment_FindByUserIDAndPurchaseID(t *testing.T) {
	assert := testAssert.New(t)
	id := primitive.NewObjectID().Hex()
	testAPI, err := service.InitTest4Mock()
	require.NoError(t, err)

	type test struct {
		name        string
		path        string
		method      string
		isOkRes     bool
		isOkMessage bool
		req         model.UserPurchaseIDCommentRequest
		fn          func(commentService *m.Comment, data test)
		expCode     int
		expRes      []model.CommentDTO
		message     string
	}

	tt := []test{
		{
			name:        "invalid id",
			path:        slash + comment + slash + user + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			req:         model.UserPurchaseIDCommentRequest{UserID: "some", PurchaseID: "some"},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByUserIDAndPurchaseID", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusBadRequest,
			message: "not correct user id",
		},
		{
			name:        "find err",
			path:        slash + comment + slash + user + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			req: model.UserPurchaseIDCommentRequest{
				UserID:     id,
				PurchaseID: id,
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByUserIDAndPurchaseID", mock.Anything, data.req).
					Return(data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + comment + slash + user + slash,
			method: http.MethodGet,
			req: model.UserPurchaseIDCommentRequest{
				UserID:     id,
				PurchaseID: id,
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByUserIDAndPurchaseID", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + comment + slash + user + slash,
			method:  http.MethodGet,
			isOkRes: true,
			req: model.UserPurchaseIDCommentRequest{
				UserID:     id,
				PurchaseID: id,
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByUserIDAndPurchaseID", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: []model.CommentDTO{
				{
					ID:         id,
					UserID:     id,
					PurchaseID: id,
					Date:       time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					Text:       "some text",
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var c []model.CommentDTO
			comment := new(m.Comment)
			testAPI.Services.Comment = comment
			router := newComment(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(comment, tc)
			}
			fullPath := tc.path + tc.req.UserID + slash + purchase + slash + tc.req.PurchaseID
			req, err := http.NewRequest(tc.method, fullPath, nil)
			assert.Nil(err)

			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			assert.Equal(tc.expCode, res.Code)

			switch {
			case tc.isOkMessage:
				err = json.NewDecoder(res.Body).Decode(&r)
				assert.Nil(err)
				assert.Equal(tc.message, r)
			case tc.isOkRes:
				err = json.NewDecoder(res.Body).Decode(&c)
				assert.Nil(err)
				assert.Equal(tc.expRes, c)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}

func TestComment_FindByText(t *testing.T) {
	assert := testAssert.New(t)
	id := primitive.NewObjectID().Hex()
	testAPI, err := service.InitTest4Mock()
	require.NoError(t, err)

	type test struct {
		name        string
		path        string
		method      string
		isOkRes     bool
		isOkMessage bool
		req         model.TextCommentRequest
		fn          func(commentService *m.Comment, data test)
		expCode     int
		expRes      []model.CommentDTO
		message     string
	}

	tt := []test{
		{
			name:        "invalid text",
			path:        slash + comment + slash + text,
			method:      http.MethodPost,
			isOkMessage: true,
			req: model.TextCommentRequest{
				Text: "",
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByText", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusBadRequest,
			message: "text is required",
		},
		{
			name:        "find err",
			path:        slash + comment + slash + text,
			method:      http.MethodPost,
			isOkMessage: true,
			req: model.TextCommentRequest{
				Text: "some",
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByText", mock.Anything, data.req).
					Return(data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + comment + slash + text,
			method: http.MethodPost,
			req: model.TextCommentRequest{
				Text: "some",
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByText", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + comment + slash + text,
			method:  http.MethodPost,
			isOkRes: true,
			req: model.TextCommentRequest{
				Text: "some",
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByText", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: []model.CommentDTO{
				{
					ID:         id,
					UserID:     id,
					PurchaseID: id,
					Date:       time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					Text:       "some text",
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var c []model.CommentDTO
			comment := new(m.Comment)
			testAPI.Services.Comment = comment
			router := newComment(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(comment, tc)
			}

			body := new(bytes.Buffer)
			err := json.NewEncoder(body).Encode(&tc.req)
			assert.Nil(err)

			req, err := http.NewRequest(tc.method, tc.path, body)
			assert.Nil(err)

			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			assert.Equal(tc.expCode, res.Code)

			switch {
			case tc.isOkMessage:
				err = json.NewDecoder(res.Body).Decode(&r)
				assert.Nil(err)
				assert.Equal(tc.message, r)
			case tc.isOkRes:
				err = json.NewDecoder(res.Body).Decode(&c)
				assert.Nil(err)
				assert.Equal(tc.expRes, c)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}

func TestComment_FindByPeriod(t *testing.T) {
	assert := testAssert.New(t)
	id := primitive.NewObjectID().Hex()
	testAPI, err := service.InitTest4Mock()
	require.NoError(t, err)

	type test struct {
		name        string
		path        string
		method      string
		isOkRes     bool
		isOkMessage bool
		req         model.PeriodCommentRequest
		fn          func(commentService *m.Comment, data test)
		expCode     int
		expRes      []model.CommentDTO
		message     string
	}

	tt := []test{
		{
			name:        "invalid period",
			path:        slash + comment + slash + period,
			method:      http.MethodPost,
			isOkMessage: true,
			req: model.PeriodCommentRequest{
				Start: time.Time{},
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByPeriod", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusBadRequest,
			message: "invalid start",
		},
		{
			name:        "find err",
			path:        slash + comment + slash + period,
			method:      http.MethodPost,
			isOkMessage: true,
			req: model.PeriodCommentRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByPeriod", mock.Anything, data.req).
					Return(data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + comment + slash + period,
			method: http.MethodPost,
			req: model.PeriodCommentRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByPeriod", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + comment + slash + period,
			method:  http.MethodPost,
			isOkRes: true,
			req: model.PeriodCommentRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(commentService *m.Comment, data test) {
				commentService.On("FindByPeriod", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: []model.CommentDTO{
				{
					ID:         id,
					UserID:     id,
					PurchaseID: id,
					Date:       time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					Text:       "some text",
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var c []model.CommentDTO
			comment := new(m.Comment)
			testAPI.Services.Comment = comment
			router := newComment(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(comment, tc)
			}

			body := new(bytes.Buffer)
			err := json.NewEncoder(body).Encode(&tc.req)
			assert.Nil(err)

			req, err := http.NewRequest(tc.method, tc.path, body)
			assert.Nil(err)

			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			assert.Equal(tc.expCode, res.Code)

			switch {
			case tc.isOkMessage:
				err = json.NewDecoder(res.Body).Decode(&r)
				assert.Nil(err)
				assert.Equal(tc.message, r)
			case tc.isOkRes:
				err = json.NewDecoder(res.Body).Decode(&c)
				assert.Nil(err)
				assert.Equal(tc.expRes, c)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}
