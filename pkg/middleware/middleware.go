package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type request interface {
	Build(*http.Request) error
	Validate() error
}

// SwagError represents a struct for swagger errors.
type SwagError struct {
	Message string `json:"message"`
}

// SwagEmptyError represents a struct for swagger errors without message.
type SwagEmptyError struct {
}

// ParseRequest parses request from http Request, stores it in the value pointed to by s and validates it.
// You must close r.Body in the Build method if you used it.
func ParseRequest(r *http.Request, s request) error {
	err := s.Build(r)
	if err != nil {
		return err
	}
	return s.Validate()
}

// JSONReturn returns server response in JSON format.
func JSONReturn(w http.ResponseWriter, statusCode int, jsonObject interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(jsonObject)
	if err != nil {
		fmt.Printf("could not encode json :%v", err.Error())
	}
}

// Empty returns server response just with status code.
func Empty(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
}

// JSONError returns error from server in JSON format.
func JSONError(w http.ResponseWriter, err error, httpStatus int) {
	JSONReturn(w, httpStatus, err.Error())
}
