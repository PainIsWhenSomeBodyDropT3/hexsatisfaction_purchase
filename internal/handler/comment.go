package handler

import (
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/service"
	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/auth"
	"github.com/gorilla/mux"
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

	return handler
}
