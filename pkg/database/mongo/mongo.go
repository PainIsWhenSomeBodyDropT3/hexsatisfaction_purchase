package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/JesusG2000/hexsatisfaction_purchase/internal/config"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const timeout = 5

// NewMongo creates new connection to mongo database.
func NewMongo(ctx context.Context, cfg config.MongoConfig) (*mongo.Database, error) {
	cfg.URI = fmt.Sprintf("%s://%s:%d", cfg.DatabaseDialect, cfg.Host, cfg.Port)
	c, err := mongo.NewClient(options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create mongo client")
	}

	err = c.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = checkConnection(ctx, c)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't connect to mongo client")
	}

	return c.Database(cfg.DatabaseName), nil
}

func checkConnection(ctx context.Context, c *mongo.Client) error {
	ctx, cancel := context.WithTimeout(ctx, timeout*time.Second)
	defer cancel()
	if err := c.Ping(ctx, readpref.Primary()); err != nil {
		return errors.Wrap(err, "couldn't ping to database")
	}
	return nil
}
