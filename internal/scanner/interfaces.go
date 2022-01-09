package scanner

import (
	"context"

	"gbu-scanner/internal/entity"
)

// Blog is interface for getting posts from blog
// (Expected that implementation gets posts from go.dev/blog/all).
type Blog interface {
	// GetPosts returns all available posts from blog.
	// Returned posts are ordered from newest to oldest.
	GetPosts(context.Context) ([]entity.Post, error)
}

// Publisher is interface for interacting with message broker
// to publish events about new posts.
type Publisher interface {
	// Publish publishes post to message broker and
	// other services can process it anyhow
	Publish(context.Context, entity.Post) error
}

// Posts is interface for interacting with storage where
// information about published posts stored.
// Notice: "published posts" is not same thing as "posted in blog":
// "Published" means "published to message broker".
type Posts interface {
	Transaction(ctx context.Context, fn func(txCtx context.Context) error) error
	// Add saves post to list of published posts
	Add(ctx context.Context, post entity.Post) error
	// GetAll reutrns all posts published to a message broker
	GetAll(ctx context.Context) ([]entity.Post, error)
}
