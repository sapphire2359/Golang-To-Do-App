package api

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"todoapp-json/logic"
	"todoapp-json/models"
)

// handler to retrieve/get all todo items
func ListTodoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	todos, err := logic.ListTodoItems(ctx)
	if err != nil {
		http.Error(w, "could not fetch todos", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// handler to create a new todo item to the list
func AddTodoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// var req addTodoRequest
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var item models.TodoItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := logic.AddTodoItem(ctx, item.Description, item.Status)
	if err != nil {
		if err.Error() == "missing description" || err.Error() == "invalid status" {
			http.Error(w, err.Error(), http.StatusBadRequest) // returns HTTP 400
			return
		}
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// handler to update a todo item from the list using id
func UpdateTodoHandler(write http.ResponseWriter, request *http.Request) {
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
		http.Error(write, "Update failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	write.WriteHeader(http.StatusOK)
}

// handler to delete a todo item
func DeleteTodoHandler(write http.ResponseWriter, request *http.Request) {
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
	// calling logic layer
	if err := logic.DeleteTodoItem(ctx, id); err != nil {
		if err.Error() == "todo item not found" {
			http.Error(write, err.Error(), http.StatusNotFound) // returns HTTP 400
			return
		}
		http.Error(write, "Delete failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	write.WriteHeader(http.StatusOK)
}

// handler to list data dynamically in a web page
func ListHandler(write http.ResponseWriter, request *http.Request) {
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
