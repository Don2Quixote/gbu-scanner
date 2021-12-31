package scanner

import (
	"context"
	"time"

	"gbu-scanner/internal/entity"
	"gbu-scanner/pkg/logger"

	"github.com/pkg/errors"
)

type scanner struct {
	posts     Posts
	publisher Publisher
	repo      Repository
	interval  time.Duration
	log       logger.Logger
}

// New returns new scanner with main business-logic of this service - method Scan
func New(posts Posts, publisher Publisher, repo Repository, interval time.Duration, log logger.Logger) *scanner {
	return &scanner{
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
func (s *scanner) Scan(ctx context.Context) error {
	s.log.Info("starting scanning")

	firstIteration := true
	for {
		// First iteration should not wait
		if !firstIteration && !sleepWithContext(ctx, s.interval) {
			break // loop will be stopped if context is closed
		}
		firstIteration = false

		posts, err := s.posts.GetAll(ctx)
		if err != nil {
			err = errors.Wrap(err, "can't get posts")
			s.log.Error(err)
			continue
		}
		if len(posts) == 0 {
			s.log.Warn("0 posts")
			continue
		}

		publihsedPosts, err := s.repo.GetPublishedPosts(ctx)
		if err != nil {
			err = errors.Wrap(err, "can't get published posts")
			s.log.Error(err)
			continue
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
			continue
		}

		// Publish not published posts from oldest to newest
		// (in most cases expected only one not published post per scan iteration)
		for i := len(notPublishedPosts) - 1; i >= 0; i-- {
			s.log.Infof("publishing post post %q", notPublishedPosts[i].Title)
			err = s.publisher.Publish(ctx, notPublishedPosts[i])
			if err != nil {
				err = errors.Wrap(err, "can't publish post")
				s.log.Error(err)
			}

			// The saddest story - post published, but can't submit this information, so post will be published again
			// It is a problem "at least once / at most once", where I have chosen "at least once"
			err = s.repo.AddPublishedPost(ctx, notPublishedPosts[i])
			if err != nil {
				err = errors.Wrap(err, "can't add published post")
				s.log.Error(err)
			}
		}
	}

	s.log.Info("scanning finished")

	return nil
}

// sleepWithContext block for specified time.Duration
// If context closes sooner than time passes, false returned, true otherwise
func sleepWithContext(ctx context.Context, duration time.Duration) bool {
	select {
	case <-ctx.Done():
		return false
	case <-time.After(duration):
		return true
	}
}
