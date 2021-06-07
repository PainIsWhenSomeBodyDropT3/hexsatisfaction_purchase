package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
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
	last   = "last"
	after  = "after"
	before = "before"
)

func TestPurchase_Create(t *testing.T) {
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
		req     model.CreatePurchaseRequest
		fn      func(purchaseService *m.Purchase, data test)
		expCode int
		expBody string
	}

	tt := []test{
		{
			name:   "invalid user id",
			path:   slash + purchase + slash + api + slash,
			method: http.MethodPost,
			req: model.CreatePurchaseRequest{
				Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				FileID: id,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("Create", mock.Anything, data.req).
					Return("", nil)
			},
			expCode: http.StatusBadRequest,
			expBody: "not correct user id",
		},
		{
			name:   "create err",
			path:   slash + purchase + slash + api + slash,
			method: http.MethodPost,
			req: model.CreatePurchaseRequest{
				UserID: 1,
				Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				FileID: id,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("Create", mock.Anything, data.req).
					Return("", errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "all ok",
			path:   slash + purchase + slash + api + slash,
			method: http.MethodPost,
			req: model.CreatePurchaseRequest{
				UserID: 1,
				Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				FileID: id,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("Create", mock.Anything, data.req).
					Return(data.expBody, nil)
			},
			expCode: http.StatusOK,
			expBody: id,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			purchaseService := new(m.Purchase)
			testAPI.Services.Purchase = purchaseService
			router := newPurchase(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(purchaseService, tc)
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

func TestPurchase_Delete(t *testing.T) {
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
		req     model.DeletePurchaseRequest
		fn      func(purchaseService *m.Purchase, data test)
		expCode int
		expBody string
	}

	tt := []test{
		{
			name:    "invalid  id",
			path:    slash + purchase + slash + api + slash,
			method:  http.MethodDelete,
			isOkRes: true,
			req: model.DeletePurchaseRequest{
				ID: "some",
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("Delete", mock.Anything, data.req).
					Return("", nil)
			},
			expCode: http.StatusBadRequest,
			expBody: "not correct id",
		},
		{
			name:    "delete err",
			path:    slash + purchase + slash + api + slash,
			method:  http.MethodDelete,
			isOkRes: true,
			req: model.DeletePurchaseRequest{
				ID: id,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("Delete", mock.Anything, data.req).
					Return("", errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + purchase + slash + api + slash,
			method: http.MethodDelete,
			req: model.DeletePurchaseRequest{
				ID: id,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("Delete", mock.Anything, data.req).
					Return("", nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + purchase + slash + api + slash,
			method:  http.MethodDelete,
			isOkRes: true,
			req: model.DeletePurchaseRequest{
				ID: id,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("Delete", mock.Anything, data.req).
					Return(data.expBody, nil)
			},
			expCode: http.StatusOK,
			expBody: id,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			purchaseService := new(m.Purchase)
			testAPI.Services.Purchase = purchaseService
			router := newPurchase(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(purchaseService, tc)
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

func TestPurchase_FindById(t *testing.T) {
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
		isOkMessage bool
		isOkRes     bool
		req         model.IDPurchaseRequest
		fn          func(purchaseService *m.Purchase, data test)
		expCode     int
		expRes      model.PurchaseDTO
		message     string
	}

	tt := []test{
		{
			name:        "invalid  id",
			path:        slash + purchase + slash + api + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			req: model.IDPurchaseRequest{
				ID: "some",
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByID", mock.Anything, data.req).
					Return(&data.expRes, nil)
			},
			expCode: http.StatusBadRequest,
			message: "not correct id",
		},
		{
			name:        "find err",
			path:        slash + purchase + slash + api + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			req: model.IDPurchaseRequest{
				ID: id,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByID", mock.Anything, data.req).
					Return(&data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + purchase + slash + api + slash,
			method: http.MethodGet,
			req: model.IDPurchaseRequest{
				ID: id,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByID", mock.Anything, data.req).
					Return(&data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + purchase + slash + api + slash,
			method:  http.MethodGet,
			isOkRes: true,
			req: model.IDPurchaseRequest{
				ID: id,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByID", mock.Anything, data.req).
					Return(&data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: model.PurchaseDTO{
				ID:     id,
				UserID: 1,
				Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				FileID: id,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var p model.PurchaseDTO
			purchaseService := new(m.Purchase)
			testAPI.Services.Purchase = purchaseService
			router := newPurchase(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(purchaseService, tc)
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
				err = json.NewDecoder(res.Body).Decode(&p)
				assert.Nil(err)
				assert.Equal(tc.expRes, p)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}

func TestPurchase_FindLastByUserId(t *testing.T) {
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
		isOkMessage bool
		isOkRes     bool
		req         model.UserIDPurchaseRequest
		fn          func(purchaseService *m.Purchase, data test)
		expCode     int
		expRes      model.PurchaseDTO
		message     string
	}

	tt := []test{
		{
			name:        "invalid  id",
			path:        slash + purchase + slash + api + slash + last + slash + user + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			req:         model.UserIDPurchaseRequest{},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindLastByUserID", mock.Anything, data.req).
					Return(&data.expRes, nil)
			},
			expCode: http.StatusBadRequest,
			message: "not correct id",
		},
		{
			name:        "find err",
			path:        slash + purchase + slash + api + slash + last + slash + user + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			req:         model.UserIDPurchaseRequest{ID: 1},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindLastByUserID", mock.Anything, data.req).
					Return(&data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + purchase + slash + api + slash + last + slash + user + slash,
			method: http.MethodGet,
			req: model.UserIDPurchaseRequest{
				ID: 1,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindLastByUserID", mock.Anything, data.req).
					Return(&data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + purchase + slash + api + slash + last + slash + user + slash,
			method:  http.MethodGet,
			isOkRes: true,
			req: model.UserIDPurchaseRequest{
				ID: 1,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindLastByUserID", mock.Anything, data.req).
					Return(&data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: model.PurchaseDTO{
				ID:     id,
				UserID: 1,
				Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				FileID: id,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var p model.PurchaseDTO
			purchaseService := new(m.Purchase)
			testAPI.Services.Purchase = purchaseService
			router := newPurchase(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(purchaseService, tc)
			}

			req, err := http.NewRequest(tc.method, tc.path+strconv.Itoa(tc.req.ID), nil)
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
				err = json.NewDecoder(res.Body).Decode(&p)
				assert.Nil(err)
				assert.Equal(tc.expRes, p)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}

func TestPurchase_FindAllByUserId(t *testing.T) {
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
		isOkMessage bool
		isOkRes     bool
		req         model.UserIDPurchaseRequest
		fn          func(purchaseService *m.Purchase, data test)
		expCode     int
		expRes      []model.PurchaseDTO
		message     string
	}

	tt := []test{
		{
			name:        "invalid  id",
			path:        slash + purchase + slash + api + slash + user + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			req:         model.UserIDPurchaseRequest{},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindAllByUserID", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusBadRequest,
			message: "not correct id",
		},
		{
			name:        "find err",
			path:        slash + purchase + slash + api + slash + user + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			req:         model.UserIDPurchaseRequest{ID: 1},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindAllByUserID", mock.Anything, data.req).
					Return(data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + purchase + slash + api + slash + user + slash,
			method: http.MethodGet,
			req: model.UserIDPurchaseRequest{
				ID: 1,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindAllByUserID", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + purchase + slash + api + slash + user + slash,
			method:  http.MethodGet,
			isOkRes: true,
			req: model.UserIDPurchaseRequest{
				ID: 1,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindAllByUserID", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: []model.PurchaseDTO{
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var p []model.PurchaseDTO
			purchaseService := new(m.Purchase)
			testAPI.Services.Purchase = purchaseService
			router := newPurchase(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(purchaseService, tc)
			}

			req, err := http.NewRequest(tc.method, tc.path+strconv.Itoa(tc.req.ID), nil)
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
				err = json.NewDecoder(res.Body).Decode(&p)
				assert.Nil(err)
				assert.Equal(tc.expRes, p)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}

func TestPurchase_FindByUserIdAndPeriod(t *testing.T) {
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
		isOkMessage bool
		isOkRes     bool
		req         model.UserIDPeriodPurchaseRequest
		fn          func(purchaseService *m.Purchase, data test)
		expCode     int
		expRes      []model.PurchaseDTO
		message     string
	}

	tt := []test{
		{
			name:        "invalid  id",
			path:        slash + purchase + slash + api + slash + period + slash + user + slash,
			method:      http.MethodPost,
			isOkMessage: true,
			req: model.UserIDPeriodPurchaseRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByUserIDAndPeriod", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusBadRequest,
			message: "not correct id",
		},
		{
			name:        "find err",
			path:        slash + purchase + slash + api + slash + period + slash + user + slash,
			method:      http.MethodPost,
			isOkMessage: true,
			req: model.UserIDPeriodPurchaseRequest{
				ID:    1,
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByUserIDAndPeriod", mock.Anything, data.req).
					Return(data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + purchase + slash + api + slash + period + slash + user + slash,
			method: http.MethodPost,
			req: model.UserIDPeriodPurchaseRequest{
				ID:    1,
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByUserIDAndPeriod", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + purchase + slash + api + slash + period + slash + user + slash,
			method:  http.MethodPost,
			isOkRes: true,
			req: model.UserIDPeriodPurchaseRequest{
				ID:    1,
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByUserIDAndPeriod", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: []model.PurchaseDTO{
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var p []model.PurchaseDTO
			purchaseService := new(m.Purchase)
			testAPI.Services.Purchase = purchaseService
			router := newPurchase(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(purchaseService, tc)
			}

			body := new(bytes.Buffer)
			err := json.NewEncoder(body).Encode(&tc.req)
			assert.Nil(err)

			req, err := http.NewRequest(tc.method, tc.path+strconv.Itoa(tc.req.ID), body)
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
				err = json.NewDecoder(res.Body).Decode(&p)
				assert.Nil(err)
				assert.Equal(tc.expRes, p)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}

func TestPurchase_FindByUserIdAfterDate(t *testing.T) {
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
		isOkMessage bool
		isOkRes     bool
		req         model.UserIDAfterDatePurchaseRequest
		fn          func(purchaseService *m.Purchase, data test)
		expCode     int
		expRes      []model.PurchaseDTO
		message     string
	}

	tt := []test{
		{
			name:        "invalid  id",
			path:        slash + purchase + slash + api + slash + after + slash + user + slash,
			method:      http.MethodPost,
			isOkMessage: true,
			req: model.UserIDAfterDatePurchaseRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByUserIDAfterDate", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusBadRequest,
			message: "not correct id",
		},
		{
			name:        "find err",
			path:        slash + purchase + slash + api + slash + after + slash + user + slash,
			method:      http.MethodPost,
			isOkMessage: true,
			req: model.UserIDAfterDatePurchaseRequest{
				ID:    1,
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByUserIDAfterDate", mock.Anything, data.req).
					Return(data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + purchase + slash + api + slash + after + slash + user + slash,
			method: http.MethodPost,
			req: model.UserIDAfterDatePurchaseRequest{
				ID:    1,
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByUserIDAfterDate", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + purchase + slash + api + slash + after + slash + user + slash,
			method:  http.MethodPost,
			isOkRes: true,
			req: model.UserIDAfterDatePurchaseRequest{
				ID:    1,
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByUserIDAfterDate", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: []model.PurchaseDTO{
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var p []model.PurchaseDTO
			purchaseService := new(m.Purchase)
			testAPI.Services.Purchase = purchaseService
			router := newPurchase(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(purchaseService, tc)
			}

			body := new(bytes.Buffer)
			err := json.NewEncoder(body).Encode(&tc.req)
			assert.Nil(err)

			req, err := http.NewRequest(tc.method, tc.path+strconv.Itoa(tc.req.ID), body)
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
				err = json.NewDecoder(res.Body).Decode(&p)
				assert.Nil(err)
				assert.Equal(tc.expRes, p)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}

func TestPurchase_FindByUserIdBeforeDate(t *testing.T) {
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
		isOkMessage bool
		isOkRes     bool
		req         model.UserIDBeforeDatePurchaseRequest
		fn          func(purchaseService *m.Purchase, data test)
		expCode     int
		expRes      []model.PurchaseDTO
		message     string
	}

	tt := []test{
		{
			name:        "invalid  id",
			path:        slash + purchase + slash + api + slash + before + slash + user + slash,
			method:      http.MethodPost,
			isOkMessage: true,
			req: model.UserIDBeforeDatePurchaseRequest{
				End: time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByUserIDBeforeDate", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusBadRequest,
			message: "not correct id",
		},
		{
			name:        "find err",
			path:        slash + purchase + slash + api + slash + before + slash + user + slash,
			method:      http.MethodPost,
			isOkMessage: true,
			req: model.UserIDBeforeDatePurchaseRequest{
				ID:  1,
				End: time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByUserIDBeforeDate", mock.Anything, data.req).
					Return(data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + purchase + slash + api + slash + before + slash + user + slash,
			method: http.MethodPost,
			req: model.UserIDBeforeDatePurchaseRequest{
				ID:  1,
				End: time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByUserIDBeforeDate", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + purchase + slash + api + slash + before + slash + user + slash,
			method:  http.MethodPost,
			isOkRes: true,
			req: model.UserIDBeforeDatePurchaseRequest{
				ID:  1,
				End: time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByUserIDBeforeDate", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: []model.PurchaseDTO{
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var p []model.PurchaseDTO
			purchaseService := new(m.Purchase)
			testAPI.Services.Purchase = purchaseService
			router := newPurchase(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(purchaseService, tc)
			}

			body := new(bytes.Buffer)
			err := json.NewEncoder(body).Encode(&tc.req)
			assert.Nil(err)

			req, err := http.NewRequest(tc.method, tc.path+strconv.Itoa(tc.req.ID), body)
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
				err = json.NewDecoder(res.Body).Decode(&p)
				assert.Nil(err)
				assert.Equal(tc.expRes, p)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}

func TestPurchase_FindByUserIdAndFileID(t *testing.T) {
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
		isOkMessage bool
		isOkRes     bool
		req         model.UserIDFileIDPurchaseRequest
		fn          func(purchaseService *m.Purchase, data test)
		expCode     int
		expRes      []model.PurchaseDTO
		message     string
	}

	tt := []test{
		{
			name:        "invalid  id",
			path:        slash + purchase + slash + api + slash + user + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			req: model.UserIDFileIDPurchaseRequest{
				FileID: id,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByUserIDAndFileID", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusBadRequest,
			message: "not correct user id",
		},
		{
			name:        "find err",
			path:        slash + purchase + slash + api + slash + user + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			req: model.UserIDFileIDPurchaseRequest{
				UserID: 1,
				FileID: id,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByUserIDAndFileID", mock.Anything, data.req).
					Return(data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + purchase + slash + api + slash + user + slash,
			method: http.MethodGet,
			req: model.UserIDFileIDPurchaseRequest{
				UserID: 1,
				FileID: id,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByUserIDAndFileID", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + purchase + slash + api + slash + user + slash,
			method:  http.MethodGet,
			isOkRes: true,
			req: model.UserIDFileIDPurchaseRequest{
				UserID: 1,
				FileID: id,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByUserIDAndFileID", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: []model.PurchaseDTO{
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var p []model.PurchaseDTO
			purchaseService := new(m.Purchase)
			testAPI.Services.Purchase = purchaseService
			router := newPurchase(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(purchaseService, tc)
			}

			fullPath := tc.path + strconv.Itoa(tc.req.UserID) + slash + file + slash + tc.req.FileID
			req, err := http.NewRequest(tc.method, fullPath, nil)
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
				err = json.NewDecoder(res.Body).Decode(&p)
				assert.Nil(err)
				assert.Equal(tc.expRes, p)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}

func TestPurchase_FindLast(t *testing.T) {
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
		isOkMessage bool
		isOkRes     bool
		fn          func(purchaseService *m.Purchase, data test)
		expCode     int
		expRes      model.PurchaseDTO
		message     string
	}

	tt := []test{
		{
			name:        "find err",
			path:        slash + purchase + slash + api + slash + last + slash,
			method:      http.MethodGet,
			isOkMessage: true,

			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindLast", mock.Anything).
					Return(&data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + purchase + slash + api + slash + last + slash,
			method: http.MethodGet,

			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindLast", mock.Anything).
					Return(&data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + purchase + slash + api + slash + last + slash,
			method:  http.MethodGet,
			isOkRes: true,

			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindLast", mock.Anything).
					Return(&data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: model.PurchaseDTO{
				ID:     id,
				UserID: 1,
				Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				FileID: id,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var p model.PurchaseDTO
			purchaseService := new(m.Purchase)
			testAPI.Services.Purchase = purchaseService
			router := newPurchase(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(purchaseService, tc)
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
				err = json.NewDecoder(res.Body).Decode(&p)
				assert.Nil(err)
				assert.Equal(tc.expRes, p)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}

func TestPurchase_FindAll(t *testing.T) {
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
		isOkMessage bool
		isOkRes     bool
		fn          func(purchaseService *m.Purchase, data test)
		expCode     int
		expRes      []model.PurchaseDTO
		message     string
	}

	tt := []test{
		{
			name:        "find err",
			path:        slash + purchase + slash + api + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindAll", mock.Anything).
					Return(data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + purchase + slash + api + slash,
			method: http.MethodGet,
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindAll", mock.Anything).
					Return(data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + purchase + slash + api + slash,
			method:  http.MethodGet,
			isOkRes: true,
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindAll", mock.Anything).
					Return(data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: []model.PurchaseDTO{
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var p []model.PurchaseDTO
			purchaseService := new(m.Purchase)
			testAPI.Services.Purchase = purchaseService
			router := newPurchase(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(purchaseService, tc)
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
				err = json.NewDecoder(res.Body).Decode(&p)
				assert.Nil(err)
				assert.Equal(tc.expRes, p)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}

func TestPurchase_FindByPeriod(t *testing.T) {
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
		isOkMessage bool
		isOkRes     bool
		req         model.PeriodPurchaseRequest
		fn          func(purchaseService *m.Purchase, data test)
		expCode     int
		expRes      []model.PurchaseDTO
		message     string
	}

	tt := []test{
		{
			name:        "invalid  start date",
			path:        slash + purchase + slash + api + slash + period,
			method:      http.MethodPost,
			isOkMessage: true,
			req: model.PeriodPurchaseRequest{
				Start: time.Time{},
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByPeriod", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusBadRequest,
			message: "start date is required",
		},
		{
			name:   "find err",
			path:   slash + purchase + slash + api + slash + period,
			method: http.MethodPost,
			req: model.PeriodPurchaseRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByPeriod", mock.Anything, data.req).
					Return(data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + purchase + slash + api + slash + period,
			method: http.MethodPost,
			req: model.PeriodPurchaseRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByPeriod", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + purchase + slash + api + slash + period,
			method:  http.MethodPost,
			isOkRes: true,
			req: model.PeriodPurchaseRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
				End:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByPeriod", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: []model.PurchaseDTO{
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var p []model.PurchaseDTO
			purchaseService := new(m.Purchase)
			testAPI.Services.Purchase = purchaseService
			router := newPurchase(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(purchaseService, tc)
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

			switch {
			case tc.isOkMessage:
				err = json.NewDecoder(res.Body).Decode(&r)
				assert.Nil(err)
				assert.Equal(tc.message, r)
			case tc.isOkRes:
				err = json.NewDecoder(res.Body).Decode(&p)
				assert.Nil(err)
				assert.Equal(tc.expRes, p)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}

func TestPurchase_FindAfterDate(t *testing.T) {
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
		isOkMessage bool
		isOkRes     bool
		req         model.AfterDatePurchaseRequest
		fn          func(purchaseService *m.Purchase, data test)
		expCode     int
		expRes      []model.PurchaseDTO
		message     string
	}

	tt := []test{
		{
			name:        "invalid start date",
			path:        slash + purchase + slash + api + slash + after,
			method:      http.MethodPost,
			isOkMessage: true,
			req: model.AfterDatePurchaseRequest{
				Start: time.Time{},
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindAfterDate", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusBadRequest,
			message: "start date is required",
		},
		{
			name:        "find err",
			path:        slash + purchase + slash + api + slash + after,
			method:      http.MethodPost,
			isOkMessage: true,
			req: model.AfterDatePurchaseRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindAfterDate", mock.Anything, data.req).
					Return(data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + purchase + slash + api + slash + after,
			method: http.MethodPost,
			req: model.AfterDatePurchaseRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindAfterDate", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + purchase + slash + api + slash + after,
			method:  http.MethodPost,
			isOkRes: true,
			req: model.AfterDatePurchaseRequest{
				Start: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindAfterDate", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: []model.PurchaseDTO{
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var p []model.PurchaseDTO
			purchaseService := new(m.Purchase)
			testAPI.Services.Purchase = purchaseService
			router := newPurchase(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(purchaseService, tc)
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

			switch {
			case tc.isOkMessage:
				err = json.NewDecoder(res.Body).Decode(&r)
				assert.Nil(err)
				assert.Equal(tc.message, r)
			case tc.isOkRes:
				err = json.NewDecoder(res.Body).Decode(&p)
				assert.Nil(err)
				assert.Equal(tc.expRes, p)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}

func TestPurchase_FindBeforeDate(t *testing.T) {
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
		isOkMessage bool
		isOkRes     bool
		req         model.BeforeDatePurchaseRequest
		fn          func(purchaseService *m.Purchase, data test)
		expCode     int
		expRes      []model.PurchaseDTO
		message     string
	}

	tt := []test{
		{
			name:        "invalid  end date",
			path:        slash + purchase + slash + api + slash + before,
			method:      http.MethodPost,
			isOkMessage: true,
			req: model.BeforeDatePurchaseRequest{
				End: time.Time{},
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindBeforeDate", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusBadRequest,
			message: "end date is required",
		},
		{
			name:        "find err",
			path:        slash + purchase + slash + api + slash + before,
			method:      http.MethodPost,
			isOkMessage: true,
			req: model.BeforeDatePurchaseRequest{
				End: time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindBeforeDate", mock.Anything, data.req).
					Return(data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + purchase + slash + api + slash + before,
			method: http.MethodPost,
			req: model.BeforeDatePurchaseRequest{
				End: time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindBeforeDate", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + purchase + slash + api + slash + before,
			method:  http.MethodPost,
			isOkRes: true,
			req: model.BeforeDatePurchaseRequest{
				End: time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindBeforeDate", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: []model.PurchaseDTO{
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var p []model.PurchaseDTO
			purchaseService := new(m.Purchase)
			testAPI.Services.Purchase = purchaseService
			router := newPurchase(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(purchaseService, tc)
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

			switch {
			case tc.isOkMessage:
				err = json.NewDecoder(res.Body).Decode(&r)
				assert.Nil(err)
				assert.Equal(tc.message, r)
			case tc.isOkRes:
				err = json.NewDecoder(res.Body).Decode(&p)
				assert.Nil(err)
				assert.Equal(tc.expRes, p)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}

func TestPurchase_FindByFileID(t *testing.T) {
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
		isOkMessage bool
		isOkRes     bool
		req         model.FileIDPurchaseRequest
		fn          func(purchaseService *m.Purchase, data test)
		expCode     int
		expRes      []model.PurchaseDTO
		message     string
	}

	tt := []test{
		{
			name:        "invalid  id",
			path:        slash + purchase + slash + api + slash + file + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			req: model.FileIDPurchaseRequest{
				FileID: "some",
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByFileID", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusBadRequest,
			message: "not correct file id",
		},
		{
			name:        "find err",
			path:        slash + purchase + slash + api + slash + file + slash,
			method:      http.MethodGet,
			isOkMessage: true,
			req: model.FileIDPurchaseRequest{
				FileID: id,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByFileID", mock.Anything, data.req).
					Return(data.expRes, errors.New(""))
			},
			expCode: http.StatusInternalServerError,
		},
		{
			name:   "not found",
			path:   slash + purchase + slash + api + slash + file + slash,
			method: http.MethodGet,
			req: model.FileIDPurchaseRequest{
				FileID: id,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByFileID", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusNotFound,
		},
		{
			name:    "all ok",
			path:    slash + purchase + slash + api + slash + file + slash,
			method:  http.MethodGet,
			isOkRes: true,
			req: model.FileIDPurchaseRequest{
				FileID: id,
			},
			fn: func(purchaseService *m.Purchase, data test) {
				purchaseService.On("FindByFileID", mock.Anything, data.req).
					Return(data.expRes, nil)
			},
			expCode: http.StatusOK,
			expRes: []model.PurchaseDTO{
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
				{
					ID:     id,
					UserID: 1,
					Date:   time.Date(2009, time.December, 10, 23, 0, 0, 0, time.Local),
					FileID: id,
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var r string
			var p []model.PurchaseDTO
			purchaseService := new(m.Purchase)
			testAPI.Services.Purchase = purchaseService
			router := newPurchase(testAPI.Services, testAPI.TokenManager)
			if tc.fn != nil {
				tc.fn(purchaseService, tc)
			}

			req, err := http.NewRequest(tc.method, tc.path+tc.req.FileID, nil)
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
				err = json.NewDecoder(res.Body).Decode(&p)
				assert.Nil(err)
				assert.Equal(tc.expRes, p)
			default:
				assert.Equal(tc.message, r)
			}
		})
	}
}
