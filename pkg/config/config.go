package config

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/pkg/errors"
)

// Parse is a wrapper around package confita to parse env variables into struct
// Example:
// var cfg struct {
//     Value string `config:"VALUE,required"`
// }
// err = config.Parse(&cfg)
//
// "required" mark in struct tag marks field as required and Parse will return
// error if value not found for field.
func Parse(cfg interface{}) error {
	err := confita.NewLoader(env.NewBackend()).Load(context.Background(), cfg)
	if err != nil {
		return errors.Wrap(err, "load env vars")
	}
	return nil
}
