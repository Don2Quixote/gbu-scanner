package app

import (
	"context"
	"net/http"
	"time"

	"gbu-scanner/internal/posts"
	"gbu-scanner/internal/publisher"
	"gbu-scanner/internal/repository"
	"gbu-scanner/internal/scanner"

	"gbu-scanner/pkg/config"
	"gbu-scanner/pkg/logger"
	"gbu-scanner/pkg/wrappers/mongo"
	"gbu-scanner/pkg/wrappers/rabbit"

	"github.com/pkg/errors"
)

// Run runs app. If returned error is not nil, program exited
// unexpectedly and non-zero code should be returned (os.Exit(1) or log.Fatal(...))
func Run(ctx context.Context, log logger.Logger) error {
	log.Info("starting app")

	var cfg appConfig
	err := config.Parse(&cfg)
	if err != nil {
		return errors.Wrap(err, "can't parse config")
	}
	cfg.setDefaults(log)

	mongo, err := mongo.NewClient(cfg.MongoHost, cfg.MongoUser, cfg.MongoPass, cfg.MongoSRV)
	if err != nil {
		return errors.Wrap(err, "can't create mongo client")
	}

	err = mongo.Connect(ctx)
	if err != nil {
		return errors.Wrap(err, "can't connect to mongo")
	}
	defer func() {
		err := mongo.Disconnect(ctx)
		if err != nil {
			log.Error(errors.Wrap(err, "can't disconnect mongo client"))
		}
	}()

	err = mongo.Ping(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "can't ping mongo")
	}

	rabbit, err := rabbit.Dial(cfg.RabbitHost, cfg.RabbitUser, cfg.RabbitPass, cfg.RabbitVhost, cfg.RabbitAmqps)
	if err != nil {
		return errors.Wrap(err, "can't connect to rabbit")
	}
	defer rabbit.Close()

	rabbitChan, err := rabbit.Channel()
	if err != nil {
		return errors.Wrap(err, "can't get rabbit channel")
	}

	publisher := publisher.New(rabbitChan)
	err = publisher.Init(ctx)
	if err != nil {
		return errors.Wrap(err, "can't init publisher")
	}
	repo := repository.New(mongo, cfg.MongoDatabase, log)
	posts := posts.New(cfg.BlogHost, cfg.BlogPath, cfg.BlogHTTPS, &http.Client{
		Timeout: time.Duration(cfg.BlogScanNetworkTimeout) * time.Second,
	}, log)
	scanner := scanner.New(posts, publisher, repo, time.Duration(cfg.BlogScanInterval)*time.Second, log)

	err = scanner.Scan(ctx)
	if err != nil {
		return errors.Wrap(err, "error during scanning")
	}

	log.Info("app finished")

	return nil
}
