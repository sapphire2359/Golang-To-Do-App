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
		return errors.New("invalid status")
	}
	//id := len(todos) + 1
	id := storage.GetNextId(todoItems)
	todoItem := models.TodoItem{Id: id, Description: description, Status: status}
	todoItems = append(todoItems, todoItem)

	slog.Info("Added new todo item...", "traceID", trace.GetTraceID(ctx), "id", id)
	return storage.SaveTodoItems(todoItems)
}

func ListTodoItems(ctx context.Context) ([]models.TodoItem, error) {
	slog.Info("List all todo items...", "traceID", trace.GetTraceID(ctx), "...", "...")
	return storage.LoadTodoItems()
}

func UpdateTodoItem(ctx context.Context, id int, description string, newStatus models.Status) error {
	todoItems, err := storage.LoadTodoItems()
	if err != nil {
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
					return errors.New("invalid status")
				}
				todoItems[i].Status = newStatus
			}
			updated = true
			break
		}
	}
	if !updated {
		return errors.New("todo item not found")
	}
	slog.Info("Updated todo item", "traceID", trace.GetTraceID(ctx), "id", id)
	return storage.SaveTodoItems(todoItems)
}

func DeleteTodoItem(ctx context.Context, id int) error {
	todoItems, err := storage.LoadTodoItems()
	if err != nil {
		return err
	}
	newTodoItems := make([]models.TodoItem, 0, len(todoItems))
	found := false
	for _, todoItem := range todoItems {
		if todoItem.Id == id {
			found = true
			continue
		}
		newTodoItems = append(newTodoItems, todoItem)
	}
	if !found {
		return errors.New("todo item not found")
	}
	slog.Info("Deleted todo", "traceID", trace.GetTraceID(ctx), "id", id)
	return storage.SaveTodoItems(newTodoItems)
}
