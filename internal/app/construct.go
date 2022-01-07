package app

import (
	"context"
	"net/http"
	"time"

	"gbu-scanner/internal/posts"
	"gbu-scanner/internal/publisher"
	"gbu-scanner/internal/repository"
	"gbu-scanner/internal/scanner"

	"gbu-scanner/pkg/logger"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// construct constructs all scanner's dependencies
func construct(ctx context.Context,
	cfg appConfig,
	mongo *mongo.Client,
	log logger.Logger,
) (
	scanner.Posts,
	scanner.Publisher,
	scanner.Repository,
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

	repo := repository.New(mongo, cfg.MongoDatabase, log)

	posts := posts.New(cfg.BlogHost, cfg.BlogPath, cfg.BlogHTTPS, &http.Client{
		Timeout: time.Duration(cfg.BlogScanNetworkTimeout) * time.Second,
	}, log)

	return posts, publisher, repo, nil
}
