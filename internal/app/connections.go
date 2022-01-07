package app

import (
	"context"

	"gbu-scanner/pkg/wrappers/mongo"

	"github.com/pkg/errors"
)

// makeConnections returns required connections/clients
func makeConnections(ctx context.Context, cfg appConfig) (*mongo.Client, error) {
	mongo, err := mongo.NewClient(cfg.MongoHost, cfg.MongoUser, cfg.MongoPass, cfg.MongoSRV)
	if err != nil {
		return nil, errors.Wrap(err, "can't create mongo client")
	}

	err = mongo.Connect(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to mongo")
	}

	err = mongo.Ping(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't ping mongo")
	}

	return mongo, nil
}
