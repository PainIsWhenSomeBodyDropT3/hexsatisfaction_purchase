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

type fileRouter struct {
	*mux.Router
	services     *service.Services
	tokenManager auth.TokenManager
}

func newFile(services *service.Services, tokenManager auth.TokenManager) fileRouter {
	router := mux.NewRouter().PathPrefix(filePath).Subrouter()
	handler := fileRouter{
		router,
		services,
		tokenManager,
	}

	router.Path("/{name}").
		Methods(http.MethodGet).
		HandlerFunc(handler.findByNameFile)

	router.Path("/actual/").
		Methods(http.MethodGet).
		HandlerFunc(handler.findActualFile)

	router.Path("/expired/").
		Methods(http.MethodGet).
		HandlerFunc(handler.findNotActualFile)

	router.Path("/").
		Methods(http.MethodGet).
		HandlerFunc(handler.findAllFile)

	router.Path("/added").
		Methods(http.MethodPost).
		HandlerFunc(handler.findAddedByPeriodFile)

	router.Path("/updated").
		Methods(http.MethodPost).
		HandlerFunc(handler.findUpdatedByPeriodFile)

	secure := router.PathPrefix("/api").Subrouter()
	secure.Use(handler.tokenManager.UserIdentity)

	secure.Path("/").
		Methods(http.MethodPost).
		HandlerFunc(handler.createFile)

	secure.Path("/{id}").
		Methods(http.MethodPut).
		HandlerFunc(handler.updateFile)

	secure.Path("/{id}").
		Methods(http.MethodDelete).
		HandlerFunc(handler.deleteFile)

	secure.Path("/{id}").
		Methods(http.MethodGet).
		HandlerFunc(handler.findByIDFile)

	secure.Path("/author/{id}").
		Methods(http.MethodGet).
		HandlerFunc(handler.findByAuthorIDFile)

	return handler
}

type createFileRequest struct {
	model.CreateFileRequest
}

// Build builds request for create file.
func (req *createFileRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&req.CreateFileRequest)
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

// Validate validates request for create file.
func (req *createFileRequest) Validate() error {
	switch {
	case req.Name == "":
		return fmt.Errorf("name is required")
	case req.Description == "":
		return fmt.Errorf("description is required")
	case req.Size < 1:
		return fmt.Errorf("not correct size")
	case req.Path == "":
		return fmt.Errorf("path is required")
	case req.AddDate == time.Time{}:
		return fmt.Errorf("add date is required")
	case req.UpdateDate == time.Time{}:
		return fmt.Errorf("update date is required")
	case req.AuthorID == "":
		return fmt.Errorf("not correct author id")
	default:
		return nil
	}
}

func (f *fileRouter) createFile(w http.ResponseWriter, r *http.Request) {
	var req createFileRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	id, err := f.services.File.Create(r.Context(), req.CreateFileRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, id)
}

type updateFileRequest struct {
	model.UpdateFileRequest
}

// Build builds request for update file.
func (req *updateFileRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&req.UpdateFileRequest)
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

// Validate validates request for update file.
func (req *updateFileRequest) Validate() error {
	switch {
	case req.ID == "":
		return fmt.Errorf("not correct id")
	case req.Name == "":
		return fmt.Errorf("name is required")
	case req.Description == "":
		return fmt.Errorf("description is required")
	case req.Size < 1:
		return fmt.Errorf("not correct size")
	case req.Path == "":
		return fmt.Errorf("path is required")
	case req.AddDate == time.Time{}:
		return fmt.Errorf("add date is required")
	case req.UpdateDate == time.Time{}:
		return fmt.Errorf("update date is required")
	case req.AuthorID == "":
		return fmt.Errorf("not correct author id")
	default:
		return nil
	}
}

func (f *fileRouter) updateFile(w http.ResponseWriter, r *http.Request) {
	var req updateFileRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	id, err := f.services.File.Update(r.Context(), req.UpdateFileRequest)
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

type deleteFileRequest struct {
	model.DeleteFileRequest
}

// Build builds request for delete file.
func (req *deleteFileRequest) Build(r *http.Request) error {
	vID, ok := mux.Vars(r)["id"]
	if !ok {
		return fmt.Errorf("no id")
	}

	req.ID = vID

	return nil
}

// Validate validates request for delete file.
func (req *deleteFileRequest) Validate() error {
	switch {
	case req.ID == "":
		return fmt.Errorf("not correct id")
	default:
		return nil
	}
}

func (f *fileRouter) deleteFile(w http.ResponseWriter, r *http.Request) {
	var req deleteFileRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	id, err := f.services.File.Delete(r.Context(), req.DeleteFileRequest)
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

type idFileRequest struct {
	model.IDFileRequest
}

// Build builds request to find file by id.
func (req *idFileRequest) Build(r *http.Request) error {
	vID, ok := mux.Vars(r)["id"]
	if !ok {
		return fmt.Errorf("no id")
	}

	req.ID = vID

	return nil
}

// Validate validates request to find file by id.
func (req *idFileRequest) Validate() error {
	switch {
	case req.ID == "":
		return fmt.Errorf("not correct id")
	default:
		return nil
	}
}

func (f *fileRouter) findByIDFile(w http.ResponseWriter, r *http.Request) {
	var req idFileRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	file, err := f.services.File.FindByID(r.Context(), req.IDFileRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if file.ID == "" {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, file)
}

type nameFileRequest struct {
	model.NameFileRequest
}

// Build builds request to find file by name.
func (req *nameFileRequest) Build(r *http.Request) error {
	name, ok := mux.Vars(r)["name"]
	if !ok {
		return fmt.Errorf("no id")
	}

	req.Name = name

	return nil
}

// Validate validates request to find file by name.
func (req *nameFileRequest) Validate() error {
	switch {
	case req.Name == "":
		return fmt.Errorf("name is required")
	default:
		return nil
	}
}

func (f *fileRouter) findByNameFile(w http.ResponseWriter, r *http.Request) {
	var req nameFileRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	files, err := f.services.File.FindByName(r.Context(), req.NameFileRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(files) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, files)
}

func (f *fileRouter) findAllFile(w http.ResponseWriter, r *http.Request) {
	files, err := f.services.File.FindAll(r.Context())
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(files) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, files)
}

type authorIDFileRequest struct {
	model.AuthorIDFileRequest
}

// Build builds request to find file by author id.
func (req *authorIDFileRequest) Build(r *http.Request) error {
	vID, ok := mux.Vars(r)["id"]
	if !ok {
		return fmt.Errorf("no id")
	}

	req.ID = vID

	return nil
}

// Validate validates request fto find file by author id.
func (req *authorIDFileRequest) Validate() error {
	switch {
	case req.ID == "":
		return fmt.Errorf("not correct id")
	default:
		return nil
	}
}

func (f *fileRouter) findByAuthorIDFile(w http.ResponseWriter, r *http.Request) {
	var req authorIDFileRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	files, err := f.services.File.FindByAuthorID(r.Context(), req.AuthorIDFileRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(files) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, files)
}

func (f *fileRouter) findNotActualFile(w http.ResponseWriter, r *http.Request) {
	files, err := f.services.File.FindNotActual(r.Context())
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(files) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, files)
}

func (f *fileRouter) findActualFile(w http.ResponseWriter, r *http.Request) {
	files, err := f.services.File.FindActual(r.Context())
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(files) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, files)
}

type addedPeriodFileRequest struct {
	model.AddedPeriodFileRequest
}

// Build builds request to find added file by date period.
func (req *addedPeriodFileRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&req.AddedPeriodFileRequest)
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

// Validate validates request to find added file by date period.
func (req *addedPeriodFileRequest) Validate() error {
	switch {
	case req.Start == time.Time{}:
		return fmt.Errorf("start is required")
	case req.End == time.Time{}:
		return fmt.Errorf("end is required")
	default:
		return nil
	}
}

func (f *fileRouter) findAddedByPeriodFile(w http.ResponseWriter, r *http.Request) {
	var req addedPeriodFileRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	files, err := f.services.File.FindAddedByPeriod(r.Context(), req.AddedPeriodFileRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(files) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, files)
}

type updatedPeriodFileRequest struct {
	model.UpdatedPeriodFileRequest
}

// Build builds request to find updated file by date period.
func (req *updatedPeriodFileRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&req.UpdatedPeriodFileRequest)
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

// Validate validates request to find updated file by date period.
func (req *updatedPeriodFileRequest) Validate() error {
	switch {
	case req.Start == time.Time{}:
		return fmt.Errorf("start is required")
	case req.End == time.Time{}:
		return fmt.Errorf("end is required")
	default:
		return nil
	}
}

func (f *fileRouter) findUpdatedByPeriodFile(w http.ResponseWriter, r *http.Request) {
	var req updatedPeriodFileRequest
	err := middleware.ParseRequest(r, &req)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	files, err := f.services.File.FindUpdatedByPeriod(r.Context(), req.UpdatedPeriodFileRequest)
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	if len(files) == 0 {
		middleware.Empty(w, http.StatusNotFound)
		return
	}

	middleware.JSONReturn(w, http.StatusOK, files)
}
