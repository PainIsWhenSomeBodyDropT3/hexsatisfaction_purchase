package app

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/JesusG2000/hexsatisfaction_purchase/internal/config"
	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/database/mongo"
)

// Run runs hexsatisfaction_purchase service
func Run(configPath string) {
	ctx := context.Background()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	cfg, err := config.Init(configPath)
	if err != nil {
		log.Fatal("Init config error: ", err)
	}

	_, err = mongo.NewMongo(ctx, cfg.Mongo)
	if err != nil {
		log.Fatal("Init db error: ", err)
	}
}
