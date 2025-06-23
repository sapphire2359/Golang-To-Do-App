package logic

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"todoapp-json/models"
	"todoapp-json/storage"
	"todoapp-json/trace"
)

// add todo item
func AddTodoItem(ctx context.Context, description string, status models.Status) error {
	//check if description is empty
	if strings.TrimSpace(description) == "" {
		slog.Info("Missing description - description is blank", "traceID", trace.GetTraceID(ctx), "description", description)
		return errors.New("missing description")
	}
	//status validation
	if !models.IsValidStatus(status) {
		slog.Info("Invalid status - input only [not started|started|completed]", "traceID", trace.GetTraceID(ctx), "status", status)
		return errors.New("invalid status")
	}
	respCh := make(chan storage.StoreResponse)
	storage.StoreChan <- storage.StoreCommand{
		Type:       storage.CmdAdd,
		Item:       models.TodoItem{Description: description, Status: status},
		ResponseCh: respCh,
	}
	resp := <-respCh
	slog.Info("sending todo item data to add to json file via channel...", "traceID", trace.GetTraceID(ctx), "description", description)
	return resp.Err
}

// list todo items
func ListTodoItems(ctx context.Context) ([]models.TodoItem, error) {
	respCh := make(chan storage.StoreResponse)
	storage.StoreChan <- storage.StoreCommand{
		Type:       storage.CmdList,
		ResponseCh: respCh,
	}
	resp := <-respCh
	slog.Info("sending data to retrieve all todo items from json file via channel...", "traceID", trace.GetTraceID(ctx))
	return resp.TodoItems, resp.Err
}

// update todo item
func UpdateTodoItem(ctx context.Context, id int, description string, newStatus models.Status) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	if strings.TrimSpace(description) == "" {
		return errors.New("description cannot be empty")
	}
	if !models.IsValidStatus(newStatus) {
		return errors.New("invalid status")
	}

	respCh := make(chan storage.StoreResponse)
	if id <= 0 {
		return errors.New("invalid id")
	}
	storage.StoreChan <- storage.StoreCommand{
		Type: storage.CmdUpdate,
		Item: models.TodoItem{
			Id:          id,
			Description: description,
			Status:      newStatus,
		},
		ResponseCh: respCh,
	}

	resp := <-respCh
	slog.Info("sending data to update an existing todo item in json file via channel...", "traceID", trace.GetTraceID(ctx), "id", id)
	return resp.Err
}

// delete todo item
func DeleteTodoItem(ctx context.Context, id int) error {
	respCh := make(chan storage.StoreResponse)
	storage.StoreChan <- storage.StoreCommand{
		Type:       storage.CmdDelete,
		Id:         id,
		ResponseCh: respCh,
	}
	resp := <-respCh
	slog.Info("sending data to delete an existing todo item from json file via channel...", "traceID", trace.GetTraceID(ctx), "id", id)
	return resp.Err
}
