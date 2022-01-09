package sleep

import (
	"context"
	"time"
)

// WithContext blocks for specified time.Duration or until context is closed.
// If context closes sooner than time passes, true returned, false otherwise.
func WithContext(ctx context.Context, duration time.Duration) bool {
	select {
	case <-ctx.Done():
		return true
	case <-time.After(duration):
		return false
	}
}
