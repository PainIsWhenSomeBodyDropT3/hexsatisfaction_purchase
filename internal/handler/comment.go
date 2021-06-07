package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/JesusG2000/hexsatisfaction_purchase/internal/model"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/service"
	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/auth"
	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/middleware"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type commentRouter struct {
	*mux.Router
	services     *service.Services
	tokenManager auth.TokenManager
}

func newComment(services *service.Services, tokenManager auth.TokenManager) commentRouter {
	router := mux.NewRouter().PathPrefix(commentPath).Subrouter()
	handler := commentRouter{
		router,
		services,
		tokenManager,
	}

	router.Path("/user/{id}").
		Methods(http.MethodGet).
		HandlerFunc(handler.findByUserIDComment)

	router.Path("/purchase/{id}").
		Methods(http.MethodGet).
		HandlerFunc(handler.findByPurchaseIDComment)

	router.Path("/user/{userID}/purchase/{purchaseID}").
		Methods(http.MethodGet).
		HandlerFunc(handler.findByUserIDAndPurchaseIDComment)

	router.Path("/text").
		Methods(http.MethodPost).
		HandlerFunc(handler.findByTextComment)

	router.Path("/period").
		Methods(http.MethodPost).
		HandlerFunc(handler.findByPeriodComment)

	secure := router.PathPrefix("/api").Subrouter()
	secure.Use(handler.tokenManager.UserIdentity)

	secure.Path("/").
		Methods(http.MethodGet).
		HandlerFunc(handler.findAllComment)

	secure.Path("/").
		Methods(http.MethodPost).
		HandlerFunc(handler.createComment)

	secure.Path("/{id}").
		Methods(http.MethodPut).
		HandlerFunc(handler.updateComment)

	secure.Path("/{id}").
		Methods(http.MethodDelete).
		HandlerFunc(handler.deleteComment)

	secure.Path("/{id}").
		Methods(http.MethodGet).
		HandlerFunc(handler.findByIDComment)

	return handler
}

type createCommentRequest struct {
	model.CreateCommentRequest
}

// Build builds request for create comment.
func (req *createCommentRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&req.CreateCommentRequest)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("%v", err)
		}
	}(r.Body)

	return nil
}

// Validate validates request for create comment.
func (req *createCommentRequest) Validate() error {
	switch {
	case req.UserID == 0:
		return fmt.Errorf("not correct user id")
	case req.PurchaseID == "":
		return fmt.Errorf("not correct purchase id")
	case req.Date == time.Time{}:
		return fmt.Errorf("date is required")
	case req.Text == "":
		return fmt.Errorf("text is required")
	default:
		return nil
	}
}

// @Summary Create
// @Security ApiKeyAuth
// @Tags comment
// @Description Create comment
// @Accept  json
// @Produce  json
// @Param comment body model.CreateCommentRequest true "Comment"
// @Success 200 {string} string id
// @Failure 400 {object} middleware.SwagError
// @Failure 500 {object} middleware.SwagError
// @Router /comment/api/ [post]
func (c *commentRouter) createComment(w http.ResponseWriter, r *http.Request) {
	var req createCommentRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	id, err := c.services.Comment.Create(r.Context(), req.CreateCommentRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, id)
}

type updateCommentRequest struct {
	model.UpdateCommentRequest
}

// Build builds request for update comment.
func (req *updateCommentRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&req.UpdateCommentRequest)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("%v", err)
		}
	}(r.Body)

	vID, ok := mux.Vars(r)["id"]
	if !ok {
		return fmt.Errorf("no id")
	}

	req.ID = vID

	return nil
}

// Validate validates request for update comment.
func (req *updateCommentRequest) Validate() error {
	switch {
	case !primitive.IsValidObjectID(req.ID):
		return fmt.Errorf("not correct id")
	case req.UserID == 0:
		return fmt.Errorf("not correct user id")
	case !primitive.IsValidObjectID(req.PurchaseID):
		return fmt.Errorf("not correct purchase id")
	case req.Date == time.Time{}:
		return fmt.Errorf("date is required")
	case req.Text == "":
		return fmt.Errorf("text is required")
	default:
		return nil
	}
}

// @Summary Update
// @Security ApiKeyAuth
// @Tags comment
// @Description Update comment
// @Accept  json
// @Produce  json
// @Param id path string true "Comment id"
// @Param comment body model.UpdateCommentRequest true "Comment"
// @Success 200 {string} string id
// @Failure 400 {object} middleware.SwagError
// @Failure 404 {object} middleware.SwagEmptyError "No comment"
// @Failure 500 {object} middleware.SwagError
// @Router /comment/api/{id} [put]
func (c *commentRouter) updateComment(w http.ResponseWriter, r *http.Request) {
	var req updateCommentRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	id, err := c.services.Comment.Update(r.Context(), req.UpdateCommentRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if id == "" {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, id)
}

type deleteCommentRequest struct {
	model.DeleteCommentRequest
}

// Build builds request for delete comment.
func (req *deleteCommentRequest) Build(r *http.Request) error {
	vID, ok := mux.Vars(r)["id"]
	if !ok {
		return fmt.Errorf("no id")
	}

	req.ID = vID

	return nil
}

// Validate validates request for delete comment.
func (req *deleteCommentRequest) Validate() error {
	switch {
	case !primitive.IsValidObjectID(req.ID):
		return fmt.Errorf("not correct id")
	default:
		return nil
	}
}

// @Summary Delete
// @Security ApiKeyAuth
// @Tags comment
// @Description Delete comment
// @Accept  json
// @Produce  json
// @Param id path string true "Comment id"
// @Success 200 {string} string id
// @Failure 400 {object} middleware.SwagError
// @Failure 404 {object} middleware.SwagEmptyError "No comment"
// @Failure 500 {object} middleware.SwagError
// @Router /comment/api/{id} [delete]
func (c *commentRouter) deleteComment(w http.ResponseWriter, r *http.Request) {
	var req deleteCommentRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	id, err := c.services.Comment.Delete(r.Context(), req.DeleteCommentRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if id == "" {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, id)
}

type idCommentRequest struct {
	model.IDCommentRequest
}

// Build builds request to find comment by id.
func (req *idCommentRequest) Build(r *http.Request) error {
	vID, ok := mux.Vars(r)["id"]
	if !ok {
		return fmt.Errorf("no id")
	}

	req.ID = vID

	return nil
}

// Validate validates request to find comment by id.
func (req *idCommentRequest) Validate() error {
	switch {
	case !primitive.IsValidObjectID(req.ID):
		return fmt.Errorf("not correct id")
	default:
		return nil
	}
}

// @Summary FindByID
// @Security ApiKeyAuth
// @Tags comment
// @Description Find comment by id
// @Accept  json
// @Produce  json
// @Param id path string true "Comment id"
// @Success 200 {object} model.Comment
// @Failure 400 {object} middleware.SwagError
// @Failure 404 {object} middleware.SwagEmptyError "No comment"
// @Failure 500 {object} middleware.SwagError
// @Router /comment/api/{id} [get]
func (c *commentRouter) findByIDComment(w http.ResponseWriter, r *http.Request) {
	var req idCommentRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	comment, err := c.services.Comment.FindByID(r.Context(), req.IDCommentRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if comment.ID == "" {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, comment)
}

type userIDCommentRequest struct {
	model.UserIDCommentRequest
}

// Build builds request to find comment by user id.
func (req *userIDCommentRequest) Build(r *http.Request) error {
	vID, ok := mux.Vars(r)["id"]
	if !ok {
		return fmt.Errorf("no id")
	}

	id, err := strconv.Atoi(vID)
	if err != nil {
		return errors.Wrap(err, "conversation error")
	}
	req.ID = id

	return nil
}

// Validate validates request to find comment by user id.
func (req *userIDCommentRequest) Validate() error {
	switch {
	case req.ID == 0:
		return fmt.Errorf("not correct id")
	default:
		return nil
	}
}

// @Summary FindAllByUserID
// @Tags comment
// @Description Find comments by user id
// @Accept  json
// @Produce  json
// @Param id path string true "User id"
// @Success 200 {array} model.Comment
// @Failure 400 {object} middleware.SwagError
// @Failure 404 {object} middleware.SwagEmptyError "No comments"
// @Failure 500 {object} middleware.SwagError
// @Router /comment/user/{id} [get]
func (c *commentRouter) findByUserIDComment(w http.ResponseWriter, r *http.Request) {
	var req userIDCommentRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	comments, err := c.services.Comment.FindAllByUserID(r.Context(), req.UserIDCommentRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(comments) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, comments)
}

type purchaseIDCommentRequest struct {
	model.PurchaseIDCommentRequest
}

// Build builds request to find comment by user id.
func (req *purchaseIDCommentRequest) Build(r *http.Request) error {
	vID, ok := mux.Vars(r)["id"]
	if !ok {
		return fmt.Errorf("no id")
	}

	req.ID = vID

	return nil
}

// Validate validates request to find comment by user id.
func (req *purchaseIDCommentRequest) Validate() error {
	switch {
	case !primitive.IsValidObjectID(req.ID):
		return fmt.Errorf("not correct id")
	default:
		return nil
	}
}

// @Summary FindByPurchaseID
// @Tags comment
// @Description Find comments by purchase id
// @Accept  json
// @Produce  json
// @Param id path string true "Purchase id"
// @Success 200 {array} model.Comment
// @Failure 400 {object} middleware.SwagError
// @Failure 404 {object} middleware.SwagEmptyError "No comments"
// @Failure 500 {object} middleware.SwagError
// @Router /comment/purchase/{id} [get]
func (c *commentRouter) findByPurchaseIDComment(w http.ResponseWriter, r *http.Request) {
	var req purchaseIDCommentRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	comments, err := c.services.Comment.FindByPurchaseID(r.Context(), req.PurchaseIDCommentRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(comments) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, comments)
}

type userPurchaseIDCommentRequest struct {
	model.UserPurchaseIDCommentRequest
}

// Build builds request to find comment by user id and purchase id.
func (req *userPurchaseIDCommentRequest) Build(r *http.Request) error {
	vars := mux.Vars(r)
	vUserID, ok := vars["userID"]
	if !ok {
		return fmt.Errorf("no user id")
	}
	purchaseID, ok := vars["purchaseID"]
	if !ok {
		return fmt.Errorf("no purchase id")
	}

	userID, err := strconv.Atoi(vUserID)
	if err != nil {
		return errors.Wrap(err, "conversation error")
	}

	req.UserID = userID
	req.PurchaseID = purchaseID

	return nil
}

// Validate validates request to find comment by user id and purchase id.
func (req *userPurchaseIDCommentRequest) Validate() error {
	switch {
	case req.UserID == 0:
		return fmt.Errorf("not correct user id")
	case !primitive.IsValidObjectID(req.PurchaseID):
		return fmt.Errorf("not correct purchase id")
	default:
		return nil
	}
}

// @Summary FindByUserIDAndPurchaseID
// @Tags comment
// @Description Find comments by purchase and user ids
// @Accept  json
// @Produce  json
// @Param userID path string true "User id"
// @Param purchaseID path string true "Purchase id"
// @Success 200 {array} model.Comment
// @Failure 400 {object} middleware.SwagError
// @Failure 404 {object} middleware.SwagEmptyError "No comments"
// @Failure 500 {object} middleware.SwagError
// @Router /comment/user/{userID}/purchase/{purchaseID} [get]
func (c *commentRouter) findByUserIDAndPurchaseIDComment(w http.ResponseWriter, r *http.Request) {
	var req userPurchaseIDCommentRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	comments, err := c.services.Comment.FindByUserIDAndPurchaseID(r.Context(), req.UserPurchaseIDCommentRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(comments) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, comments)
}

// @Summary FindAll
// @Security ApiKeyAuth
// @Tags comment
// @Description Find all comments
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Comment
// @Failure 400 {object} middleware.SwagError
// @Failure 404 {object} middleware.SwagEmptyError "No comments"
// @Failure 500 {object} middleware.SwagError
// @Router /comment/api/ [get]
func (c *commentRouter) findAllComment(w http.ResponseWriter, r *http.Request) {
	comments, err := c.services.Comment.FindAll(r.Context())
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(comments) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, comments)
}

type textCommentRequest struct {
	model.TextCommentRequest
}

// Build builds request to find comment by text.
func (req *textCommentRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&req.TextCommentRequest)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("%v", err)
		}
	}(r.Body)

	return nil
}

// Validate validates request to find comment by text.
func (req *textCommentRequest) Validate() error {
	switch {
	case req.Text == "":
		return fmt.Errorf("text is required")
	default:
		return nil
	}
}

// @Summary FindByText
// @Tags comment
// @Description Find comments by text
// @Accept  json
// @Produce  json
// @Param text body model.TextCommentRequest true "Comment text"
// @Success 200 {array} model.Comment
// @Failure 400 {object} middleware.SwagError
// @Failure 404 {object} middleware.SwagEmptyError "No comments"
// @Failure 500 {object} middleware.SwagError
// @Router /comment/text [post]
func (c *commentRouter) findByTextComment(w http.ResponseWriter, r *http.Request) {
	var req textCommentRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	comments, err := c.services.Comment.FindByText(r.Context(), req.TextCommentRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(comments) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, comments)
}

type periodCommentRequest struct {
	model.PeriodCommentRequest
}

// Build builds request to find comment by date period.
func (req *periodCommentRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&req.PeriodCommentRequest)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("%v", err)
		}
	}(r.Body)

	return nil
}

// Validate validates request to find comment by date period.
func (req *periodCommentRequest) Validate() error {
	switch {
	case req.Start == time.Time{}:
		return fmt.Errorf("invalid start")
	case req.End == time.Time{}:
		return fmt.Errorf("invalid end")
	default:
		return nil
	}
}

// @Summary FindByPeriod
// @Tags comment
// @Description Find comments by period
// @Accept  json
// @Produce  json
// @Param period body model.PeriodCommentRequest true "Comment period"
// @Success 200 {array} model.Comment
// @Failure 400 {object} middleware.SwagError
// @Failure 404 {object} middleware.SwagEmptyError "No comments"
// @Failure 500 {object} middleware.SwagError
// @Router /comment/period [post]
func (c *commentRouter) findByPeriodComment(w http.ResponseWriter, r *http.Request) {
	var req periodCommentRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	comments, err := c.services.Comment.FindByPeriod(r.Context(), req.PeriodCommentRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(comments) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, comments)
}
