package service

import (
	"context"
	"net"
	"os"
	"strings"

	"github.com/JesusG2000/hexsatisfaction/pkg/grpc/api"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/config"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/grpc"
	"github.com/JesusG2000/hexsatisfaction_purchase/internal/repository"
	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/auth"
	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/database/mongo"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// TestAPI represents a struct for tests api.
type TestAPI struct {
	*Services
	auth.TokenManager
	GRPCClient api.ExistanceClient
}

const configPath = "config/main"

// InitTest4Mock initialize an a TestAPI for mock testing.
func InitTest4Mock() (*TestAPI, error) {
	envPath, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get path to current directory")
	}
	envPath = strings.SplitAfter(envPath, "hexsatisfaction_purchase")[0]
	if err := godotenv.Load(envPath + "/" + ".env"); err != nil {
		return nil, errors.Wrap(err, "couldn't load env file")
	}

	configPath := strings.Split(configPath, "/")

	viper.AddConfigPath(envPath + "/" + configPath[0])
	viper.SetConfigName(configPath[1])

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "couldn't load config file")
	}

	return initServices4Test()
}

func initServices4Test() (*TestAPI, error) {
	cfg, err := config.Init(configPath)
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
