package handler

import (
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/service"
	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/auth"
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

	return handler
}
