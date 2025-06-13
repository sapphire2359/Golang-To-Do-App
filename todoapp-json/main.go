package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"todoapp-json/api"
	"todoapp-json/logic"
	"todoapp-json/models"
	"todoapp-json/trace"
)

// CLI flags
var (
	description = flag.String("desc", "", "Todo description")
	status      = flag.String("status", "not started", "Status [not started|started|completed]")
	action      = flag.String("action", "list", "Action to perform [add|list|update|delete|serve]")
	id          = flag.Int("id", 0, "Todo ID for update/delete")
)

func main() {
	flag.Parse()
	ctx := trace.WithTraceID(context.Background(), trace.NewTraceID())

	// setup structured logging
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	switch *action {
	case "add":
		// error handler when user does not input "description" or "status"
		if *description == "" || *status == "" {
			fmt.Println("-description or -status cannot be blank")
			return
		}
		logic.AddTodoItem(ctx, *description, models.Status(*status))
	case "list":
		todos, _ := logic.ListTodoItems(ctx)
		for _, todo := range todos {
			fmt.Printf("[%d] %s (%s)\n", todo.Id, todo.Description, todo.Status)
		}
	case "update":
		// error handler when user does not input "description" and "status"
		if *description == "" && *status == "" {
			fmt.Println(" Usage:   -action=update -Id=\"...\" -desc=\"...\" -status=\"...\"")
			fmt.Println(" Description and status cannot be blank")
			return
		}
		logic.UpdateTodoItem(ctx, *id, *description, models.Status(*status))
	case "delete":
		if *id > 0 {
			fmt.Println(" Usage:   -action=delete -Id=\"...\"")
			fmt.Println(" Id must be greater than 0")
		}
		logic.DeleteTodoItem(ctx, *id)
	case "serve":
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
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		slog.Info("Shutdown signal received")
		server.Shutdown(ctx)
	default:
		fmt.Println("How to use Todo App CLI :")
		fmt.Println("  -action=add -desc=\"...\" -status=\"not started|started|completed\"")
		fmt.Println("  -action=list")
		fmt.Println("  -action=update -Id=\"...\" -desc=\"...\" -status=\"...\"")
		fmt.Println("  -action=delete= -Id=\"...\"")
		fmt.Println("  -action=serve")
	}
}
