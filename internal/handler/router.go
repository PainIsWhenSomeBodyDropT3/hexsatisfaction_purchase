package handler

import (
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/service"
	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/auth"
	"github.com/gorilla/mux"
)

const (
	purchasePath = "/purchase"
	commentPath  = "/comment"
	filePath     = "/file"
)

// API represents a structure with APIs.
type API struct {
	*mux.Router
}

// NewHandler creates and serves endpoints of API.
func NewHandler(services *service.Services, tokenManager auth.TokenManager) *API {
	api := API{
		mux.NewRouter(),
	}
	api.PathPrefix(purchasePath).Handler(newPurchase(services, tokenManager))
	api.PathPrefix(commentPath).Handler(newComment(services, tokenManager))
	api.PathPrefix(filePath).Handler(newFile(services, tokenManager))

	return &api
}
