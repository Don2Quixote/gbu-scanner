package graceful

import (
	"os"
	"os/signal"
	"syscall"
)

// OnShutdown registers a handler for SIGTERM/SIGINT signals
// This function is not-blocking, launches goroutine.
func OnShutdown(action func()) {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
		<-ch
		action()
	}()
}
