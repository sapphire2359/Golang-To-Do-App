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

// handler to retrieve/get all todo items
func getHandler(write http.ResponseWriter, request *http.Request) {
	ctx := request.Context() // adding TraceID to context from middleware
	todoItems, err := logic.ListTodoItems(ctx)
	if err != nil {
		http.Error(write, "Failed to read todos", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(write).Encode(todoItems)
}

// handler to create a new todo item to the list
func createHandler(write http.ResponseWriter, request *http.Request) {
	ctx := request.Context() // adding TraceID to context from middleware
	if request.Method != http.MethodPost {
		http.Error(write, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var todo models.TodoItem
	if err := json.NewDecoder(request.Body).Decode(&todo); err != nil {
		http.Error(write, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := logic.AddTodoItem(ctx, todo.Description, todo.Status); err != nil {
		if err.Error() == "invalid status" {
			http.Error(write, err.Error(), http.StatusBadRequest) // returns HTTP 400
			return
		}
		http.Error(write, "Failed to create todo", http.StatusInternalServerError)
		return
	}
	slog.Info("Created todo", "traceID", trace.GetTraceID(ctx), "desc", todo.Description)
	write.WriteHeader(http.StatusCreated)
}

// handler to update a todo item from the list using id
func updateHandler(write http.ResponseWriter, request *http.Request) {
	ctx := request.Context() // adding TraceID to context from middleware
	if request.Method != http.MethodPut {
		http.Error(write, "Only PUT allowed", http.StatusMethodNotAllowed)
		return
	}

	var todo models.TodoItem
	if err := json.NewDecoder(request.Body).Decode(&todo); err != nil {
		http.Error(write, "Invalid request body", http.StatusBadRequest)
		return
	}
	if todo.Id == 0 {
		http.Error(write, "Missing ID", http.StatusBadRequest)
		return
	}
	if err := logic.UpdateTodoItem(ctx, todo.Id, todo.Description, todo.Status); err != nil {
		if err.Error() == "invalid status" {
			http.Error(write, err.Error(), http.StatusBadRequest) // returns HTTP 400
			return
		}

		http.Error(write, "Update failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	write.WriteHeader(http.StatusOK)
}

// handler to delete a todo item
func deleteHandler(write http.ResponseWriter, request *http.Request) {
	ctx := request.Context() // adding TraceID to context from middleware
	if request.Method != http.MethodDelete {
		http.Error(write, "Only DELETE allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := request.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(write, "Invalid or missing ID", http.StatusBadRequest)
		return
	}
	if err := logic.DeleteTodoItem(ctx, id); err != nil {
		http.Error(write, "Delete failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	write.WriteHeader(http.StatusOK)
}

// handler to list data dynamically in a web page
func listHandler(write http.ResponseWriter, request *http.Request) {
	ctx := request.Context() // adding TraceID to context from middleware
	todos, err := logic.ListTodoItems(ctx)
	if err != nil {
		http.Error(write, "Failed to retrieve todos", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/dynamic/list.html")
	if err != nil {
		http.Error(write, "Template parsing error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(write, todos); err != nil {
		http.Error(write, "Template execution error", http.StatusInternalServerError)
	}
}
