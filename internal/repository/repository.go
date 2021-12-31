package repository

import (
	"context"

	"gbu-scanner/internal/entity"
	"gbu-scanner/pkg/logger"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	mongoDB *mongo.Database
	log     logger.Logger
}

func New(mongo *mongo.Client, database string, log logger.Logger) *repository {
	return &repository{
		mongoDB: mongo.Database(database),
		log:     log,
	}
}

func (r *repository) AddPublishedPost(ctx context.Context, post entity.Post) error {
	_, err := r.mongoDB.Collection(publishedPostsCollection).InsertOne(ctx, post)
	if err != nil {
		return errors.Wrap(err, "can't insert post")
	}
	return nil
}

func (r *repository) GetPublishedPosts(ctx context.Context) ([]entity.Post, error) {
	cursor, err := r.mongoDB.Collection(publishedPostsCollection).Find(ctx, bson.D{}) // bson.D{} means find all
	if err != nil {
		return nil, errors.Wrap(err, "can't find records")
	}
	defer r.closeCursor(ctx, cursor)

	var posts []entity.Post
	err = cursor.All(ctx, &posts)
	if err != nil {
		return nil, errors.Wrap(err, "can't read all records from cursor")
	}

	return posts, nil
}
