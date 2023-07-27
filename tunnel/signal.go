package tunnel

import (
	"os"
	"os/signal"
	"syscall"
)

func NewQuitSignal() <-chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	return quit
}
