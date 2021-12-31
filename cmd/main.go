package main

import (
	"context"

	"gbu-scanner/internal/app"

	"gbu-scanner/pkg/graceful"
	"gbu-scanner/pkg/logger"

	"github.com/pkg/errors"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	log := logger.NewLogrus()

	graceful.OnShutdown(cancel)

	err := app.Run(ctx, log)
	if err != nil {
		err = errors.Wrap(err, "error running app")
		log.Fatal(err)
	}
}
