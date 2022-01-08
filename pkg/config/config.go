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
func Parse(cfg interface{}) error {
	err := confita.NewLoader(env.NewBackend()).Load(context.Background(), cfg)
	if err != nil {
		return errors.Wrap(err, "can't load env vars")
	}
	return nil
}
