package scanner

import (
	"context"

	"gbu-scanner/internal/entity"
)

// Blog is interface for getting posts from blog
// (Expected that implementation gets posts from go.dev/blog/all)
type Blog interface {
	// GetPosts returns all available posts from blog
	// Returned posts are ordered from newest to oldest
	GetPosts(context.Context) ([]entity.Post, error)
}

// Publisher is interface for interacting with message broker
// to publish events about new posts
type Publisher interface {
	// Publish publishes post to message broker and
	// other services can process it anyhow
	Publish(context.Context, entity.Post) error
}

// Repository is interface for interacting with storage where
// information about last published post stored
type Repository interface {
	// GetPublishedPosts reutrns all posts published to a message broker
	// Notice: "published posts" is not same thing as "posted in blog"
	// "Published" means "published to message broker"
	GetPublishedPosts(ctx context.Context) ([]entity.Post, error)
	// AddPublishedPost saves post to list of published posts
	// Notice: "published posts" is not same thing as "posted in blog"
	// "Published" means "published to message broker"
	AddPublishedPost(ctx context.Context, post entity.Post) error
}
