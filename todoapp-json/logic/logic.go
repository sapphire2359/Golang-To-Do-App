package logic

import (
	"context"
	"errors"
	"log/slog"
	"todoapp-json/models"
	"todoapp-json/storage"
	"todoapp-json/trace"
)

// add todo item
func AddTodoItem(ctx context.Context, description string, status models.Status) error {
	todoItems, err := storage.LoadTodoItems()
	if err != nil {
		return err
	}
	//status validation
	if !models.IsValidStatus(status) {
		slog.Info("Invalid status...", "traceID", trace.GetTraceID(ctx), "status", status)
		return errors.New("invalid status")
	}
	id := storage.GetNextId(todoItems)
	todoItem := models.TodoItem{Id: id, Description: description, Status: status}
	todoItems = append(todoItems, todoItem)

	slog.Info("Added new todo item...", "traceID", trace.GetTraceID(ctx), "id", id)
	return storage.SaveTodoItems(todoItems)
}

// list todo items
func ListTodoItems(ctx context.Context) ([]models.TodoItem, error) {
	slog.Info("List all todo items...", "traceID", trace.GetTraceID(ctx))
	return storage.LoadTodoItems()
}

// update todo item
func UpdateTodoItem(ctx context.Context, id int, description string, newStatus models.Status) error {
	todoItems, err := storage.LoadTodoItems()
	if err != nil {
		slog.Info("Error loading todo items", "traceID", trace.GetTraceID(ctx), "Todo items", todoItems)
		return err
	}
	updated := false
	for i, todoItem := range todoItems {
		if todoItem.Id == id {
			if description != "" {
				todoItems[i].Description = description
			}
			if newStatus != "" {
				//status validation
				if !models.IsValidStatus(newStatus) {
					slog.Info("Invalid status...", "traceID", trace.GetTraceID(ctx), "status", newStatus)
					return errors.New("invalid status")
				}
				todoItems[i].Status = newStatus
			}
			updated = true
			break
		}
	}
	if !updated {
		slog.Info("Todo item not found.", "traceID", trace.GetTraceID(ctx), "updated", updated)
		return errors.New("todo item not found")
	}
	slog.Info("Updated todo item", "traceID", trace.GetTraceID(ctx), "id", id)
	return storage.SaveTodoItems(todoItems)
}

// delete todo item
func DeleteTodoItem(ctx context.Context, id int) error {
	todoItems, err := storage.LoadTodoItems()
	if err != nil {
		slog.Info("Error loading todo items", "traceID", trace.GetTraceID(ctx), "Todo items", todoItems)
		return err
	}
	// newTodoItems := make([]models.TodoItem, 0, len(todoItems))
	found := false
	for index, todoItem := range todoItems {
		if todoItem.Id == id {
			todoItems = append(todoItems[:index], todoItems[index+1:]...)
			found = true
			continue
		}
		// newTodoItems = append(newTodoItems, todoItem)
	}
	if !found {
		return errors.New("todo item not found")
	}
	slog.Info("Deleted todo", "traceID", trace.GetTraceID(ctx), "id", id)
	return storage.SaveTodoItems(todoItems)
}
