package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect is wrapper for original mongo.Connect function but with
// credential-parameters instead of connection string
func Connect(ctx context.Context, host, user, pass string, srv bool) (*mongo.Client, error) {
	claims := ""
	if user != "" {
		claims = fmt.Sprintf("%s:%s@", user, pass)
	}

	connString := fmt.Sprintf("mongodb://%s%s", claims, host)
	if srv {
		// https://docs.mongodb.com/manual/reference/connection-string/#dns-seed-list-connection-format
		connString = fmt.Sprintf("mongodb+srv://%s%s", claims, host)
	}

	return mongo.Connect(ctx, options.Client().ApplyURI(connString))
}

// Forward type
type Client = mongo.Client
