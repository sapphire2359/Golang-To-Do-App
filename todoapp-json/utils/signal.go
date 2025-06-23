package utils

import (
	"os"
	"os/signal"
	"syscall"
)

// Handle graceful shutdown
func WaitForInterrupt() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	// slog.Info("Shutdown signal received")
}
