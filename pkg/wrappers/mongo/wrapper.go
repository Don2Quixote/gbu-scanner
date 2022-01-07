package mongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewClient is wrapper for original mongo.NewClient function but with
// credential-parameters instead of connection string
func NewClient(host, user, pass string, srv bool) (*mongo.Client, error) {
	claims := ""
	if user != "" {
		claims = fmt.Sprintf("%s:%s@", user, pass)
	}
	connString := fmt.Sprintf("mongodb://%s%s", claims, host)
	if srv {
		connString = fmt.Sprintf("mongodb+srv://%s%s", claims, host)
	}
	return mongo.NewClient(options.Client().ApplyURI(connString))
}

// Forward type
type Client = mongo.Client
