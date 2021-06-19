package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/JesusG2000/hexsatisfaction_purchase/internal/config"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/grpc"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/handler"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/repository"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/server"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/service"
	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/auth"
	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/database/mongo"
	"github.com/go-openapi/runtime/middleware"
)

// Run runs hexsatisfaction_purchase service
func Run() {
	ctx := context.Background()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	cfg, err := config.Init()
	if err != nil {
		log.Fatal("Init config error: ", err)
	}

	db, err := mongo.NewMongo(ctx, cfg.Mongo)
	if err != nil {
		log.Fatal("Init db error: ", err)
	}

	tokenManager, err := auth.NewManager(cfg.Auth.SigningKey)
	if err != nil {
		log.Fatal("Init jwt-token error: ", err)
	}
	addr := net.JoinHostPort(cfg.GRPC.Host, cfg.GRPC.Port)
	grpcClient, err := grpc.NewGRPCClient(addr)
	if err != nil {
		log.Fatal("Init grpc client error: ", err)
	}
	repos := repository.NewRepositories(db)

	services := service.NewServices(service.Deps{
		Repos:        repos,
		TokenManager: tokenManager,
		GRPCClient:   grpcClient,
	})

	router := handler.NewHandler(services, tokenManager)
	routeSwagger(router)

	srv := server.NewServer(cfg, router)
	go startService(ctx, srv)
	log.Printf("server started")

	<-stop

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Printf("failed to stop server: %v", err)
	}

	log.Printf("shutting down server...")
}

func startService(ctx context.Context, coreService *server.Server) {
	if err := coreService.Run(); err != nil {
		log.Fatal(ctx, "service shutdown: ", err.Error())
	}
}
func routeSwagger(router *handler.API) {
	ops := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(ops, nil)

	router.Handle("/docs", sh)
	router.Handle("/swagger.yaml", http.FileServer(http.Dir("./docs/")))
}
