package api

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"todoapp-json/logic"
	"todoapp-json/models"
	"todoapp-json/trace"
)

// registering all handlers wrapped with middleware trace function
func RegisterHandlers(mux *http.ServeMux) {
	mux.Handle("/get", TraceMiddleware(http.HandlerFunc(getHandler)))
	mux.Handle("/create", TraceMiddleware(http.HandlerFunc(createHandler)))
	mux.Handle("/update", TraceMiddleware(http.HandlerFunc(updateHandler)))
	mux.Handle("/delete", TraceMiddleware(http.HandlerFunc(deleteHandler)))

	// Serve static /about page
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve dynamic list page
	mux.Handle("/list", TraceMiddleware(http.HandlerFunc(listHandler)))
}

// handler to retrieve/get all todo items
func getHandler(w http.ResponseWriter, r *http.Request) {
	//ctx := trace.WithTraceID(r.Context(), trace.NewTraceID())
	ctx := r.Context()
	todoItems, err := logic.ListTodoItems(ctx)
	if err != nil {
		http.Error(w, "Failed to read todos", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(todoItems)
}

// handler to create a new todo item to the list
func createHandler(w http.ResponseWriter, r *http.Request) {
	//ctx := trace.WithTraceID(r.Context(), trace.NewTraceID())
	ctx := r.Context()
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var todo models.TodoItem
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := logic.AddTodoItem(ctx, todo.Description, todo.Status); err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}
	slog.Info("Created todo", "traceID", trace.GetTraceID(ctx), "desc", todo.Description)
	w.WriteHeader(http.StatusCreated)
}

// handler to update a todo item from the list using id
func updateHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := trace.WithTraceID(r.Context(), trace.NewTraceID())
	ctx := r.Context()
	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT allowed", http.StatusMethodNotAllowed)
		return
	}

	var todo models.TodoItem
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if todo.Id == 0 {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}
	if err := logic.UpdateTodoItem(ctx, todo.Id, todo.Description, todo.Status); err != nil {
		http.Error(w, "Update failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// handler to delete a todo item
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := trace.WithTraceID(r.Context(), trace.NewTraceID())
	ctx := r.Context()
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid or missing ID", http.StatusBadRequest)
		return
	}
	if err := logic.DeleteTodoItem(ctx, id); err != nil {
		http.Error(w, "Delete failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// handler to list data in dynamically in a web page
func listHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := trace.WithTraceID(r.Context(), trace.NewTraceID())
	ctx := r.Context()
	todos, err := logic.ListTodoItems(ctx)
	if err != nil {
		http.Error(w, "Failed to retrieve todos", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/dynamic/list.html")
	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, todos); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}

// middleware func that adds a TraceID to each request's context.
func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxWithTrace := trace.WithTraceID(r.Context(), trace.NewTraceID())
		r = r.WithContext(ctxWithTrace)
		next.ServeHTTP(w, r)
	})
}
