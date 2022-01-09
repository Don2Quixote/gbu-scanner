package posts

import (
	"context"

	"gbu-scanner/internal/entity"

	"gbu-scanner/pkg/logger"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Posts is implementation for scanner.Posts interface.
type Posts struct {
	mongoDB *mongo.Database
	log     logger.Logger
}

// New returns scanner.Posts implementation.
func New(mongo *mongo.Client, database string, log logger.Logger) *Posts {
	return &Posts{
		mongoDB: mongo.Database(database),
		log:     log,
	}
}

func (p *Posts) Add(ctx context.Context, post entity.Post) error {
	_, err := p.mongoDB.Collection(publishedPostsCollection).InsertOne(ctx, post)
	if err != nil {
		return errors.Wrap(err, "can't insert post")
	}
	return nil
}

func (p *Posts) GetAll(ctx context.Context) ([]entity.Post, error) {
	cursor, err := p.mongoDB.Collection(publishedPostsCollection).Find(ctx, bson.D{}) // bson.D{} means find all
	if err != nil {
		return nil, errors.Wrap(err, "can't find records")
	}
	defer p.closeCursor(ctx, cursor)

	var posts []entity.Post
	err = cursor.All(ctx, &posts)
	if err != nil {
		return nil, errors.Wrap(err, "can't read all records from cursor")
	}

	return posts, nil
}
