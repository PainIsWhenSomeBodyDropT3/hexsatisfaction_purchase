package service

import (
	"context"
	"net"

	"github.com/JesusG2000/hexsatisfaction/pkg/grpc/api"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/config"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/grpc"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/repository"
	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/auth"
	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/database/mongo"
	"github.com/pkg/errors"
)

// TestAPI represents a struct for tests api.
type TestAPI struct {
	*Services
	auth.TokenManager
	GRPCClient api.ExistanceClient
}

// InitTest4Mock initialize an a TestAPI for mock testing.
func InitTest4Mock() (*TestAPI, error) {
	test, err := initServices4Test()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't init tests for mock")
	}

	return test, nil
}

func initServices4Test() (*TestAPI, error) {
	cfg, err := config.Init()
	if err != nil {
		return nil, errors.Wrap(err, "config file error")
	}
	db, err := mongo.NewMongo(context.Background(), cfg.Mongo)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't init mongo")
	}

	tokenManager, err := auth.NewManager(cfg.Auth.SigningKey)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't init token manager")
	}

	addr := net.JoinHostPort(cfg.GRPC.Host, cfg.GRPC.Port)
	grpcClient, err := grpc.NewGRPCClient(addr)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't init grpc client")
	}

	return &TestAPI{
		Services: NewServices(Deps{
			Repos:        repository.NewRepositories(db),
			TokenManager: tokenManager,
			GRPCClient:   grpcClient,
		}),
		TokenManager: tokenManager,
		GRPCClient:   grpcClient,
	}, nil
}
