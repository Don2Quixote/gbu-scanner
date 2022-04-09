package app

import (
	"context"

	"gbu-scanner/pkg/wrappers/mongo"

	"github.com/pkg/errors"
)

// makeConnections makes required connections/clients.
func makeConnections(ctx context.Context, cfg appConfig) (*mongo.Client, error) {
	mongo, err := mongo.Connect(ctx, cfg.MongoHost, cfg.MongoUser, cfg.MongoPass, cfg.MongoSRV)
	if err != nil {
		return nil, errors.Wrap(err, "connect to mongo")
	}

	err = mongo.Ping(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "ping mongo")
	}

	return mongo, nil
}
