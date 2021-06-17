package handler

import (
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/service"
	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/auth"
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

	return handler
}
