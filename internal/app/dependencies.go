package app

import (
	"context"
	"net/http"
	"time"

	"gbu-scanner/internal/blog"
	"gbu-scanner/internal/posts"
	"gbu-scanner/internal/publisher"
	"gbu-scanner/internal/scanner"

	"gbu-scanner/pkg/logger"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// makeDependencies maeks all scanner's dependencies.
func makeDependencies(
	ctx context.Context,
	cfg appConfig,
	mongo *mongo.Client,
	log logger.Logger,
) (
	scanner.Blog,
	scanner.Publisher,
	scanner.Posts,
	error,
) {
	publisher := publisher.New(publisher.RabbitConfig{
		Host:           cfg.RabbitHost,
		User:           cfg.RabbitUser,
		Pass:           cfg.RabbitPass,
		Vhost:          cfg.RabbitVhost,
		Amqps:          cfg.RabbitAmqps,
		ReconnectDelay: time.Duration(cfg.RabbitReconnectDelay) * time.Second,
	}, log)
	err := publisher.Init(ctx)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "can't init publisher")
	}

	posts := posts.New(mongo, cfg.MongoDatabase, log)
	err = posts.Init(ctx)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "can't init database")
	}

	blog := blog.New(cfg.BlogHost, cfg.BlogPath, cfg.BlogHTTPS, &http.Client{
		Timeout: time.Duration(cfg.BlogScanNetworkTimeout) * time.Second,
	}, log)

	return blog, publisher, posts, nil
}
