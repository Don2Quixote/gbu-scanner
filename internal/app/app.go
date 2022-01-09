package app

import (
	"context"
	"time"

	"gbu-scanner/internal/scanner"

	"gbu-scanner/pkg/config"
	"gbu-scanner/pkg/logger"

	"github.com/pkg/errors"
)

// Run runs app. If returned error is not nil, program exited
// unexpectedly and non-zero code should be returned (os.Exit(1) or log.Fatal(...)).
func Run(ctx context.Context, log logger.Logger) error {
	log.Info("starting app")

	// Getting configuration
	var cfg appConfig
	err := config.Parse(&cfg)
	if err != nil {
		return errors.Wrap(err, "can't parse config")
	}
	cfg.setDefaults(log)

	// Getting required connections/clients
	mongo, err := makeConnections(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "can't make connections")
	}
	defer func() {
		err := mongo.Disconnect(ctx) // Disconnects without error if context closed
		if err != nil {
			log.Error(errors.Wrap(err, "can't disconnect mongo client"))
		}
	}()

	// Making dependencies for scanner
	blog, publisher, posts, err := makeDependencies(ctx, cfg, mongo, log)
	if err != nil {
		return errors.Wrap(err, "can't construct dependencies")
	}

	// Constructing and launching scanner
	blogScanInterval := time.Duration(cfg.BlogScanInterval) * time.Second
	err = scanner.New(blog, publisher, posts, blogScanInterval, log).Scan(ctx)
	if err != nil {
		return errors.Wrap(err, "error during scanning")
	}

	log.Info("app finished")

	return nil
}
