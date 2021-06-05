package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/JesusG2000/hexsatisfaction_purchase/internal/model"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/service"
	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/auth"
	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/middleware"
	"github.com/gorilla/mux"
)

type purchaseRouter struct {
	*mux.Router
	services     *service.Services
	tokenManager auth.TokenManager
}

func newPurchase(services *service.Services, tokenManager auth.TokenManager) purchaseRouter {
	router := mux.NewRouter().PathPrefix(purchasePath).Subrouter()
	handler := purchaseRouter{
		router,
		services,
		tokenManager,
	}

	secure := router.PathPrefix("/api").Subrouter()
	secure.Use(handler.tokenManager.UserIdentity)

	secure.Path("/{id}").
		Methods(http.MethodGet).
		HandlerFunc(handler.findByIDPurchase)

	secure.Path("/last/user/{id}").
		Methods(http.MethodGet).
		HandlerFunc(handler.findLastByUserIDPurchase)

	secure.Path("/user/{id}").
		Methods(http.MethodGet).
		HandlerFunc(handler.findAllByUserIDPurchase)

	secure.Path("/last/").
		Methods(http.MethodGet).
		HandlerFunc(handler.findLast)

	secure.Path("/").
		Methods(http.MethodGet).
		HandlerFunc(handler.findAll)

	secure.Path("/user/{userID}/file/{fileID}").
		Methods(http.MethodGet).
		HandlerFunc(handler.findByUserIDAndFileIDPurchase)

	secure.Path("/file/{fileID}").
		Methods(http.MethodGet).
		HandlerFunc(handler.findByFileIDPurchase)

	secure.Path("/").
		Methods(http.MethodPost).
		HandlerFunc(handler.createPurchase)

	secure.Path("/period/user/{id}").
		Methods(http.MethodPost).
		HandlerFunc(handler.findByUserIDAndPeriodPurchase)

	secure.Path("/after/user/{id}").
		Methods(http.MethodPost).
		HandlerFunc(handler.findByUserIDAfterDatePurchase)

	secure.Path("/before/user/{id}").
		Methods(http.MethodPost).
		HandlerFunc(handler.findByUserIDBeforeDatePurchase)

	secure.Path("/period").
		Methods(http.MethodPost).
		HandlerFunc(handler.findByPeriodPurchase)

	secure.Path("/after").
		Methods(http.MethodPost).
		HandlerFunc(handler.findAfterDatePurchase)

	secure.Path("/before").
		Methods(http.MethodPost).
		HandlerFunc(handler.findBeforeDatePurchase)

	secure.Path("/{id}").
		Methods(http.MethodDelete).
		HandlerFunc(handler.deletePurchase)

	return handler

	return handler
}

type createPurchaseRequest struct {
	model.CreatePurchaseRequest
}

// Build builds request for create purchase.
func (req *createPurchaseRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&req.CreatePurchaseRequest)
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

// Validate validates request for create purchase.
func (req *createPurchaseRequest) Validate() error {
	switch {
	case req.UserID == "":
		return fmt.Errorf("not correct user id")
	case req.Date == time.Time{}:
		return fmt.Errorf("date is required")
	case req.FileID == "":
		return fmt.Errorf("file id is required")
	default:
		return nil
	}
}

func (p *purchaseRouter) createPurchase(w http.ResponseWriter, r *http.Request) {
	var req createPurchaseRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	id, err := p.services.Purchase.Create(r.Context(), req.CreatePurchaseRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, id)
}

type deletePurchaseRequest struct {
	model.DeletePurchaseRequest
}

// Build builds request to delete purchase.
func (req *deletePurchaseRequest) Build(r *http.Request) error {
	vID, ok := mux.Vars(r)["id"]
	if !ok {
		return fmt.Errorf("no id")
	}

	req.ID = vID

	return nil
}

// Validate validates request to delete purchase.
func (req *deletePurchaseRequest) Validate() error {
	switch {
	case req.ID == "":
		return fmt.Errorf("not correct id")
	default:
		return nil
	}
}

func (p *purchaseRouter) deletePurchase(w http.ResponseWriter, r *http.Request) {
	var req deletePurchaseRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	id, err := p.services.Purchase.Delete(r.Context(), req.DeletePurchaseRequest)
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

type idPurchaseRequest struct {
	model.IDPurchaseRequest
}

// Build builds request to find purchase by id.
func (req *idPurchaseRequest) Build(r *http.Request) error {
	vID, ok := mux.Vars(r)["id"]
	if !ok {
		return fmt.Errorf("no id")
	}

	req.ID = vID

	return nil
}

// Validate validates request to find purchase by id.
func (req *idPurchaseRequest) Validate() error {
	switch {
	case req.ID == "":
		return fmt.Errorf("not correct id")
	default:
		return nil
	}
}

func (p *purchaseRouter) findByIDPurchase(w http.ResponseWriter, r *http.Request) {
	var req idPurchaseRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	purchase, err := p.services.Purchase.FindByID(r.Context(), req.IDPurchaseRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if purchase.ID == "" {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, purchase)
}

type lastUserIDPurchaseRequest struct {
	model.UserIDPurchaseRequest
}

// Build builds request to find last purchase by user id.
func (req *lastUserIDPurchaseRequest) Build(r *http.Request) error {
	vID, ok := mux.Vars(r)["id"]
	if !ok {
		return fmt.Errorf("no id")
	}

	req.ID = vID

	return nil
}

// Validate validates request to find last purchase by user id.
func (req *lastUserIDPurchaseRequest) Validate() error {
	switch {
	case req.ID == "":
		return fmt.Errorf("not correct id")
	default:
		return nil
	}
}

func (p *purchaseRouter) findLastByUserIDPurchase(w http.ResponseWriter, r *http.Request) {
	var req lastUserIDPurchaseRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	purchase, err := p.services.Purchase.FindLastByUserID(r.Context(), req.UserIDPurchaseRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if purchase.ID == "" {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, purchase)
}

type userIDPurchaseRequest struct {
	model.UserIDPurchaseRequest
}

// Build builds request to find all purchases by user id.
func (req *userIDPurchaseRequest) Build(r *http.Request) error {
	vID, ok := mux.Vars(r)["id"]
	if !ok {
		return fmt.Errorf("no id")
	}

	req.ID = vID

	return nil
}

// Validate validates request to find all purchases by user id.
func (req *userIDPurchaseRequest) Validate() error {
	switch {
	case req.ID == "":
		return fmt.Errorf("not correct id")
	default:
		return nil
	}
}

func (p *purchaseRouter) findAllByUserIDPurchase(w http.ResponseWriter, r *http.Request) {
	var req userIDPurchaseRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	purchases, err := p.services.Purchase.FindAllByUserID(r.Context(), req.UserIDPurchaseRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(purchases) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, purchases)
}

type userIDPeriodPurchaseRequest struct {
	model.UserIDPeriodPurchaseRequest
}

// Build builds request to find all purchases by user id and date period.
func (req *userIDPeriodPurchaseRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&req.UserIDPeriodPurchaseRequest)
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

// Validate validates request to find all purchases by user id and date period.
func (req *userIDPeriodPurchaseRequest) Validate() error {
	switch {
	case req.ID == "":
		return fmt.Errorf("not correct id")
	case req.Start == time.Time{}:
		return fmt.Errorf("start date is required")
	case req.End == time.Time{}:
		return fmt.Errorf("end date is required")
	default:
		return nil
	}
}

func (p *purchaseRouter) findByUserIDAndPeriodPurchase(w http.ResponseWriter, r *http.Request) {
	var req userIDPeriodPurchaseRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	purchases, err := p.services.Purchase.FindByUserIDAndPeriod(r.Context(), req.UserIDPeriodPurchaseRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(purchases) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, purchases)
}

type userIDAfterDatePurchaseRequest struct {
	model.UserIDAfterDatePurchaseRequest
}

// Build builds request to find all purchases by user id after date.
func (req *userIDAfterDatePurchaseRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&req.UserIDAfterDatePurchaseRequest)
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

// Validate validates request to find all purchases by user id after date.
func (req *userIDAfterDatePurchaseRequest) Validate() error {
	switch {
	case req.ID == "":
		return fmt.Errorf("not correct id")
	case req.Start == time.Time{}:
		return fmt.Errorf("start date is required")
	default:
		return nil
	}
}

func (p *purchaseRouter) findByUserIDAfterDatePurchase(w http.ResponseWriter, r *http.Request) {
	var req userIDAfterDatePurchaseRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	purchases, err := p.services.Purchase.FindByUserIDAfterDate(r.Context(), req.UserIDAfterDatePurchaseRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(purchases) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, purchases)
}

type userIDBeforeDatePurchaseRequest struct {
	model.UserIDBeforeDatePurchaseRequest
}

// Build builds request to find all purchases by user id before date.
func (req *userIDBeforeDatePurchaseRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&req.UserIDBeforeDatePurchaseRequest)
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

// Validate validates request to find all purchases by user id before date.
func (req *userIDBeforeDatePurchaseRequest) Validate() error {
	switch {
	case req.ID == "":
		return fmt.Errorf("not correct id")
	case req.End == time.Time{}:
		return fmt.Errorf("end date is required")
	default:
		return nil
	}
}

func (p *purchaseRouter) findByUserIDBeforeDatePurchase(w http.ResponseWriter, r *http.Request) {
	var req userIDBeforeDatePurchaseRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	purchases, err := p.services.Purchase.FindByUserIDBeforeDate(r.Context(), req.UserIDBeforeDatePurchaseRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(purchases) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, purchases)
}

type userIDFileIDPurchaseRequest struct {
	model.UserIDFileIDPurchaseRequest
}

// Build builds request to find all purchases by user id and file name.
func (req *userIDFileIDPurchaseRequest) Build(r *http.Request) error {
	vars := mux.Vars(r)
	vUserID, ok := vars["userID"]
	if !ok {
		return fmt.Errorf("no user id")
	}

	vFileID, ok := vars["fileID"]
	if !ok {
		return fmt.Errorf("no file id")
	}

	req.UserID = vUserID
	req.FileID = vFileID

	return nil
}

// Validate validates request to find all purchases by user id and file name.
func (req *userIDFileIDPurchaseRequest) Validate() error {
	switch {
	case req.UserID == "":
		return fmt.Errorf("not correct user id")
	case req.FileID == "":
		return fmt.Errorf("not correct file id")
	default:
		return nil
	}
}

func (p *purchaseRouter) findByUserIDAndFileIDPurchase(w http.ResponseWriter, r *http.Request) {
	var req userIDFileIDPurchaseRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	purchases, err := p.services.Purchase.FindByUserIDAndFileID(r.Context(), req.UserIDFileIDPurchaseRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(purchases) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, purchases)
}

func (p *purchaseRouter) findLast(w http.ResponseWriter, r *http.Request) {
	purchase, err := p.services.Purchase.FindLast(r.Context())
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if purchase.ID == "" {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, purchase)
}

func (p *purchaseRouter) findAll(w http.ResponseWriter, r *http.Request) {
	purchases, err := p.services.Purchase.FindAll(r.Context())
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(purchases) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, purchases)
}

type periodPurchaseRequest struct {
	model.PeriodPurchaseRequest
}

// Build builds request to find all purchases by date period.
func (req *periodPurchaseRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&req.PeriodPurchaseRequest)
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

// Validate validates request to find all purchases by date period.
func (req *periodPurchaseRequest) Validate() error {
	switch {
	case req.Start == time.Time{}:
		return fmt.Errorf("start date is required")
	case req.End == time.Time{}:
		return fmt.Errorf("end date is required")
	default:
		return nil
	}
}

func (p *purchaseRouter) findByPeriodPurchase(w http.ResponseWriter, r *http.Request) {
	var req periodPurchaseRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	purchases, err := p.services.Purchase.FindByPeriod(r.Context(), req.PeriodPurchaseRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(purchases) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, purchases)
}

type afterDatePurchaseRequest struct {
	model.AfterDatePurchaseRequest
}

// Build builds request to find all purchases after date.
func (req *afterDatePurchaseRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&req.AfterDatePurchaseRequest)
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

// Validate validates request to find all purchases after date.
func (req *afterDatePurchaseRequest) Validate() error {
	switch {
	case req.Start == time.Time{}:
		return fmt.Errorf("start date is required")
	default:
		return nil
	}
}

func (p *purchaseRouter) findAfterDatePurchase(w http.ResponseWriter, r *http.Request) {
	var req afterDatePurchaseRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	purchases, err := p.services.Purchase.FindAfterDate(r.Context(), req.AfterDatePurchaseRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(purchases) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, purchases)
}

type beforeDatePurchaseRequest struct {
	model.BeforeDatePurchaseRequest
}

// Build builds request to find all purchases before date.
func (req *beforeDatePurchaseRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&req.BeforeDatePurchaseRequest)
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

// Validate validates request to find all purchases before date.
func (req *beforeDatePurchaseRequest) Validate() error {
	switch {
	case req.End == time.Time{}:
		return fmt.Errorf("end date is required")
	default:
		return nil
	}
}

func (p *purchaseRouter) findBeforeDatePurchase(w http.ResponseWriter, r *http.Request) {
	var req beforeDatePurchaseRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	purchases, err := p.services.Purchase.FindBeforeDate(r.Context(), req.BeforeDatePurchaseRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(purchases) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, purchases)
}

type fileIDPurchaseRequest struct {
	model.FileIDPurchaseRequest
}

// Build builds request to find all purchases by file name.
func (req *fileIDPurchaseRequest) Build(r *http.Request) error {
	vID, ok := mux.Vars(r)["fileID"]
	if !ok {
		return fmt.Errorf("no file id")
	}

	req.FileID = vID

	return nil
}

// Validate validates request to find all purchases by file name.
func (req *fileIDPurchaseRequest) Validate() error {
	switch {
	case req.FileID == "":
		return fmt.Errorf("not correct file id")
	default:
		return nil
	}
}

func (p *purchaseRouter) findByFileIDPurchase(w http.ResponseWriter, r *http.Request) {
	var req fileIDPurchaseRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	purchases, err := p.services.Purchase.FindByFileID(r.Context(), req.FileIDPurchaseRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(purchases) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, purchases)
}
