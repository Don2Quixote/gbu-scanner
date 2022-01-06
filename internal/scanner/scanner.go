package scanner

import (
	"context"
	"time"

	"gbu-scanner/internal/entity"
	"gbu-scanner/pkg/logger"

	"github.com/pkg/errors"
)

// Scanner is struct that incapsulates business-logic's dependencies (interfaces) and configuration
type Scanner struct {
	posts     Posts
	publisher Publisher
	repo      Repository
	interval  time.Duration
	log       logger.Logger
}

// New returns new scanner with main business-logic of this service - method Scan
func New(posts Posts, publisher Publisher, repo Repository, interval time.Duration, log logger.Logger) *Scanner {
	return &Scanner{
		posts:     posts,
		publisher: publisher,
		repo:      repo,
		interval:  interval,
		log:       log,
	}
}

// Scan is a blocking method until context cancelled, it does blog's posts scanning in a loop
// Once new post posted in blog, information about it published to message broker and
// consumers (other services) can do whatever they please with this information
// Scanning current implementation always returns nil-error when context is closed
func (s *Scanner) Scan(ctx context.Context) error {
	s.log.Info("starting scanning")

	// Loop executes scanning interations with specified inteval (s.interval)
	for isCtxClosed := false; !isCtxClosed; isCtxClosed = sleepWithContext(ctx, s.interval) {
		for _, err := range s.scanIteration(ctx) {
			s.log.Error(errors.Wrap(err, "error during scanning"))
		}
	}

	s.log.Info("scanning finished")

	return nil
}

// scanIteration called in Scan method to reduce it's loop's complexity
// More than one error allowed in iteration so it returns []error
func (s *Scanner) scanIteration(ctx context.Context) []error {
	var errs []error

	posts, err := s.posts.GetAll(ctx)
	if err != nil {
		return append(errs, errors.Wrap(err, "can't get posts"))
	}
	if len(posts) == 0 {
		s.log.Warn("0 posts")
		return nil
	}

	publihsedPosts, err := s.repo.GetPublishedPosts(ctx)
	if err != nil {
		return append(errs, errors.Wrap(err, "can't get published posts"))
	}

	var notPublishedPosts []entity.Post
	for _, p := range posts {
		found := false
		for _, pp := range publihsedPosts {
			if p.URL == pp.URL {
				found = true
				break
			}
		}
		if !found {
			notPublishedPosts = append(notPublishedPosts, p)
		}
	}

	if len(notPublishedPosts) == 0 {
		s.log.Info("no new posts")
		return nil
	}

	// Publish not published posts from oldest to newest
	// (in most cases expected only one not published post per scan iteration)
	for i := len(notPublishedPosts) - 1; i >= 0; i-- {
		s.log.Infof("publishing post post %q", notPublishedPosts[i].Title)
		err = s.publisher.Publish(ctx, notPublishedPosts[i])
		if err != nil {
			errs = append(errs, errors.Wrap(err, "can't publish post"))
			continue
		}

		// The saddest story - post published, but can't submit this information, so post will be published again
		// It is a problem "at least once / at most once", where I have chosen "at least once"
		err = s.repo.AddPublishedPost(ctx, notPublishedPosts[i])
		if err != nil {
			errs = append(errs, errors.Wrap(err, "can't add published post"))
		}
	}

	return errs
}

// sleepWithContext block for specified time.Duration or until context is closed
// If context closes sooner than time passes, true returned, false otherwise
func sleepWithContext(ctx context.Context, duration time.Duration) bool {
	select {
	case <-ctx.Done():
		return true
	case <-time.After(duration):
		return false
	}
}
