package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"todoapp-json/api"
	"todoapp-json/cmd"
	"todoapp-json/storage"
	"todoapp-json/trace"
	"todoapp-json/utils"
)

func main() {
	// context
	ctx := trace.WithTraceID(context.Background(), trace.NewTraceID())
	// setup structured logging
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	// Start the command loop
	storage.StoreChan = make(chan storage.StoreCommand, 100)
	// go routine
	go storage.StartStoreLoop()

	// CLI mode check
	if len(os.Args) > 1 {
		// if strings.Split(os.Args[1], "=")[1] == "REPL" {
		// 	cmd.RunREPL(context.Background())
		// 	return
		// }
		cmd.RunCLI(context.Background())
		return
	}

	// HTTP server setup
	mux := http.NewServeMux()
	api.RegisterHandlers(mux)
	server := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		slog.Info("Server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", "err", err)
		}
	}()

	// Handle graceful shutdown
	utils.WaitForInterrupt()
	slog.Info("Shutdown signal received")
	server.Shutdown(ctx)
}
