package posts

import (
	"context"

	"gbu-scanner/internal/entity"
	"gbu-scanner/internal/scanner"

	"gbu-scanner/pkg/logger"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// Posts is implementation for scanner.Posts interface.
type Posts struct {
	mongo   *mongo.Client
	mongoDB *mongo.Database
	log     logger.Logger
}

var _ scanner.Posts = &Posts{}

// New returns scanner.Posts implementation.
func New(mongo *mongo.Client, database string, log logger.Logger) *Posts {
	return &Posts{
		mongo:   mongo,
		mongoDB: mongo.Database(database),
		log:     log,
	}
}

// Init creates document with empty posts array if it doesn't exist.
func (p *Posts) Init(ctx context.Context) error {
	res := p.mongoDB.Collection(publishedPostsCollection).FindOne(ctx, bson.D{})
	if res.Err() != nil && !errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return errors.Wrap(res.Err(), "find document")
	}

	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		_, err := p.mongoDB.Collection(publishedPostsCollection).InsertOne(ctx, bson.D{
			{Key: "posts", Value: bson.A{}},
		})
		if err != nil {
			return errors.Wrap(err, "insert document")
		}
	}

	return nil
}

// Because of transationc's usage it's impossible to use app with standalone mongo instance:
// https://docs.mongodb.com/manual/core/transactions/#feature-compatibility-version--fcv-
// Why it's impossible to do transactions in a standalone? Who knows.
func (p *Posts) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	sess, err := p.mongo.StartSession(options.Session().SetCausalConsistency(true))
	if err != nil {
		return errors.Wrap(err, "start session")
	}
	defer sess.EndSession(ctx)

	txOpts := options.Transaction().
		SetReadConcern(readconcern.Snapshot()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	_, err = sess.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx)
	}, txOpts)
	if err != nil {
		return errors.Wrap(err, "make transaction")
	}

	return nil
}

func (p *Posts) Add(ctx context.Context, post entity.Post) error {
	// As only one document with posts array in collection - empty filter used
	_, err := p.mongoDB.Collection(publishedPostsCollection).UpdateOne(ctx, bson.D{}, bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "posts", Value: post},
		}},
	})
	if err != nil {
		return errors.Wrap(err, "insert post")
	}

	return nil
}

func (p *Posts) GetAll(ctx context.Context) ([]entity.Post, error) {
	// As only one document with posts array in collection - empty filter used
	res := p.mongoDB.Collection(publishedPostsCollection).FindOne(ctx, bson.D{})
	if res.Err() != nil {
		return nil, errors.Wrap(res.Err(), "find document")
	}

	var doc struct {
		Posts []entity.Post `bson:"posts"`
	}

	err := res.Decode(&doc)
	if err != nil {
		return nil, errors.Wrap(err, "decode document")
	}

	return doc.Posts, nil
}
