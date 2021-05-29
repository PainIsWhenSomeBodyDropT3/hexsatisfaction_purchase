package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/JesusG2000/hexsatisfaction_purchase/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const timeout = 5

func NewMongo(ctx context.Context, cfg config.MongoConfig) (*mongo.Database, error) {
	cfg.URI = fmt.Sprintf("%s://%s:%d", cfg.Dialect, cfg.Host, cfg.Port)
	c, err := mongo.NewClient(options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}

	err = c.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = checkConnection(ctx, c)
	if err != nil {
		return nil, err
	}

	return c.Database(cfg.Name), nil
}

func checkConnection(ctx context.Context, c *mongo.Client) error {
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()
	return c.Ping(ctx, readpref.Primary())
}
