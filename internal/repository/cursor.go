package repository

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *Repository) closeCursor(ctx context.Context, cur *mongo.Cursor) {
	err := cur.Close(ctx)
	if err != nil {
		err = errors.Wrap(err, "can't close cursor")
		r.log.Error(err)
	}
}
